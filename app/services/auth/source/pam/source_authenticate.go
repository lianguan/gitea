// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package pam

import (
	"context"
	"fmt"
	"strings"

	"gitmin.com/gitmin/app/models/auth"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/auth/pam"
	"gitmin.com/gitmin/app/modules/optional"
	"gitmin.com/gitmin/app/modules/setting"

	"github.com/google/uuid"
)

// Authenticate queries if login/password is valid against the PAM,
// and create a local user if success when enabled.
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, userName, password string) (*user_model.User, error) {
	pamLogin, err := pam.Auth(source.ServiceName, userName, password)
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failure") {
			return nil, user_model.ErrUserNotExist{Name: userName}
		}
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	// Allow PAM sources with `@` in their name, like from Active Directory
	username := pamLogin
	email := pamLogin
	idx := strings.Index(pamLogin, "@")
	if idx > -1 {
		username = pamLogin[:idx]
	}
	if user_model.ValidateEmail(email) != nil {
		if source.EmailDomain != "" {
			email = fmt.Sprintf("%s@%s", username, source.EmailDomain)
		} else {
			email = fmt.Sprintf("%s@%s", username, setting.Service.NoReplyAddress)
		}
		if user_model.ValidateEmail(email) != nil {
			email = uuid.New().String() + "@localhost"
		}
	}

	user = &user_model.User{
		LowerName:   strings.ToLower(username),
		Name:        username,
		Email:       email,
		Passwd:      password,
		LoginType:   auth.PAM,
		LoginSource: source.authSource.ID,
		LoginName:   userName, // This is what the user typed in
	}
	overwriteDefault := &user_model.CreateUserOverwriteOptions{
		IsActive: optional.Some(true),
	}

	if err := user_model.CreateUser(ctx, user, &user_model.Meta{}, overwriteDefault); err != nil {
		return user, err
	}

	return user, nil
}

// IsSkipLocalTwoFA returns if this source should skip local 2fa for password authentication
func (source *Source) IsSkipLocalTwoFA() bool {
	return source.SkipLocalTwoFA
}

// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webauthn

import (
	"context"
	"encoding/binary"
	"encoding/gob"

	"gitmin.com/gitmin/app/models/auth"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/app/modules/util"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthn represents the global WebAuthn instance
var WebAuthn *webauthn.WebAuthn

// Init initializes the WebAuthn instance from the config.
func Init() {
	gob.Register(&webauthn.SessionData{})

	appURL, _ := protocol.FullyQualifiedOrigin(setting.AppURL)

	WebAuthn = &webauthn.WebAuthn{
		Config: &webauthn.Config{
			RPDisplayName: setting.AppName,
			RPID:          setting.Domain,
			RPOrigins:     []string{appURL},
			AuthenticatorSelection: protocol.AuthenticatorSelection{
				UserVerification: protocol.VerificationDiscouraged,
			},
			AttestationPreference: protocol.PreferDirectAttestation,
		},
	}
}

// user represents an implementation of webauthn.User based on User model
type user struct {
	ctx  context.Context
	User *user_model.User

	defaultAuthFlags protocol.AuthenticatorFlags
}

var _ webauthn.User = (*user)(nil)

func NewWebAuthnUser(ctx context.Context, u *user_model.User, defaultAuthFlags ...protocol.AuthenticatorFlags) webauthn.User {
	return &user{ctx: ctx, User: u, defaultAuthFlags: util.OptionalArg(defaultAuthFlags)}
}

// WebAuthnID implements the webauthn.User interface
func (u *user) WebAuthnID() []byte {
	id := make([]byte, 8)
	binary.PutVarint(id, u.User.ID)
	return id
}

// WebAuthnName implements the webauthn.User interface
func (u *user) WebAuthnName() string {
	return util.IfZero(u.User.LoginName, u.User.Name)
}

// WebAuthnDisplayName implements the webauthn.User interface
func (u *user) WebAuthnDisplayName() string {
	return u.User.DisplayName()
}

// WebAuthnCredentials implements the webauthn.User interface
func (u *user) WebAuthnCredentials() []webauthn.Credential {
	dbCreds, err := auth.GetWebAuthnCredentialsByUID(u.ctx, u.User.ID)
	if err != nil {
		return nil
	}
	return dbCreds.ToCredentials(u.defaultAuthFlags)
}
// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	auth_model "code.gitea.io/gitea/app/models/auth"
	"code.gitea.io/gitea/app/services/auth"
	"code.gitea.io/gitea/app/services/auth/source/db"
)

// This test file exists to assert that our Source exposes the interfaces that we expect
// It tightly binds the interfaces and implementation without breaking go import cycles

type sourceInterface interface {
	auth.PasswordAuthenticator
	auth_model.Config
}

var _ (sourceInterface) = &db.Source{}
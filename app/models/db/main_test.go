// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	"testing"

	"code.gitea.io/gitea/app/models/unittest"

	_ "code.gitea.io/gitea/app/models"
	_ "code.gitea.io/gitea/app/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
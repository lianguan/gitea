// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth_test

import (
	"testing"

	"code.gitea.io/gitea/app/models/unittest"

	_ "code.gitea.io/gitea/app/models"
	_ "code.gitea.io/gitea/app/models/actions"
	_ "code.gitea.io/gitea/app/models/activities"
	_ "code.gitea.io/gitea/app/models/auth"
	_ "code.gitea.io/gitea/app/models/perm/access"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

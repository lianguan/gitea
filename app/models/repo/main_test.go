// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"code.gitea.io/gitea/app/models/unittest"

	_ "code.gitea.io/gitea/app/models" // register table model
	_ "code.gitea.io/gitea/app/models/actions"
	_ "code.gitea.io/gitea/app/models/activities"
	_ "code.gitea.io/gitea/app/models/perm/access" // register table model
	_ "code.gitea.io/gitea/app/models/repo"        // register table model
	_ "code.gitea.io/gitea/app/models/user"        // register table model
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

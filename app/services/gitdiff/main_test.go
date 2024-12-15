// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gitdiff

import (
	"testing"

	"code.gitea.io/gitea/app/models/unittest"

	_ "code.gitea.io/gitea/app/models"
	_ "code.gitea.io/gitea/app/models/actions"
	_ "code.gitea.io/gitea/app/models/activities"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

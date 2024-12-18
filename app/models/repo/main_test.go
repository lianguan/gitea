// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models" // register table model
	_ "gitmin.com/gitmin/app/models/actions"
	_ "gitmin.com/gitmin/app/models/activities"
	_ "gitmin.com/gitmin/app/models/perm/access" // register table model
	_ "gitmin.com/gitmin/app/models/repo"        // register table model
	_ "gitmin.com/gitmin/app/models/user"        // register table model
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models"
	_ "gitmin.com/gitmin/app/models/actions"
	_ "gitmin.com/gitmin/app/models/activities"
	_ "gitmin.com/gitmin/app/models/auth"
	_ "gitmin.com/gitmin/app/models/perm/access"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package organization_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models"
	_ "gitmin.com/gitmin/app/models/actions"
	_ "gitmin.com/gitmin/app/models/activities"
	_ "gitmin.com/gitmin/app/models/organization"
	_ "gitmin.com/gitmin/app/models/repo"
	_ "gitmin.com/gitmin/app/models/user"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

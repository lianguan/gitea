// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models"
	_ "gitmin.com/gitmin/app/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

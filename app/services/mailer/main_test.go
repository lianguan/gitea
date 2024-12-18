// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package mailer

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models"
	_ "gitmin.com/gitmin/app/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

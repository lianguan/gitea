// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package convert

import (
	"testing"

	"code.gitea.io/gitea/app/models/unittest"

	_ "code.gitea.io/gitea/app/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

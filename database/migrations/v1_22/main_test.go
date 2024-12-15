// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_22 //nolint

import (
	"testing"

	"code.gitea.io/gitea/database/migrations/base"
)

func TestMain(m *testing.M) {
	base.MainTest(m)
}
// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package project

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		FixtureFiles: []string{
			"project.yml",
			"project_board.yml",
			"project_issue.yml",
			"repository.yml",
		},
	})
}

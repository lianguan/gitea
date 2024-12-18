// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues_test

import (
	"testing"

	issues_model "gitmin.com/gitmin/app/models/issues"
	"gitmin.com/gitmin/app/models/unittest"

	_ "gitmin.com/gitmin/app/models"
	_ "gitmin.com/gitmin/app/models/actions"
	_ "gitmin.com/gitmin/app/models/activities"
	_ "gitmin.com/gitmin/app/models/repo"
	_ "gitmin.com/gitmin/app/models/user"

	"github.com/stretchr/testify/assert"
)

func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&issues_model.Issue{},
		&issues_model.MergeRequest{},
		&issues_model.Milestone{},
		&issues_model.Label{},
	)
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

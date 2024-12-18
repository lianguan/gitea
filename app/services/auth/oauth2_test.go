// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth

import (
	"context"
	"testing"

	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/web/middleware"
	"gitmin.com/gitmin/app/services/actions"

	"github.com/stretchr/testify/assert"
)

func TestUserIDFromToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	t.Run("Actions JWT", func(t *testing.T) {
		const RunningTaskID = 47
		token, err := actions.CreateAuthorizationToken(RunningTaskID, 1, 2)
		assert.NoError(t, err)

		ds := make(middleware.ContextData)

		o := OAuth2{}
		uid := o.userIDFromToken(context.Background(), token, ds)
		assert.Equal(t, int64(user_model.ActionsUserID), uid)
		assert.Equal(t, true, ds["IsActionsToken"])
		assert.Equal(t, ds["ActionsTaskID"], int64(RunningTaskID))
	})
}

func TestCheckTaskIsRunning(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	cases := map[string]struct {
		TaskID   int64
		Expected bool
	}{
		"Running":   {TaskID: 47, Expected: true},
		"Missing":   {TaskID: 1, Expected: false},
		"Cancelled": {TaskID: 46, Expected: false},
	}

	for name := range cases {
		c := cases[name]
		t.Run(name, func(t *testing.T) {
			actual := CheckTaskIsRunning(context.Background(), c.TaskID)
			assert.Equal(t, c.Expected, actual)
		})
	}
}

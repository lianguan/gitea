// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/db"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"

	"github.com/stretchr/testify/assert"
)

func TestIsFollowing(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.True(t, user_model.IsFollowing(db.DefaultContext, 4, 2))
	assert.False(t, user_model.IsFollowing(db.DefaultContext, 2, 4))
	assert.False(t, user_model.IsFollowing(db.DefaultContext, 5, unittest.NonexistentID))
	assert.False(t, user_model.IsFollowing(db.DefaultContext, unittest.NonexistentID, 5))
	assert.False(t, user_model.IsFollowing(db.DefaultContext, unittest.NonexistentID, unittest.NonexistentID))
}

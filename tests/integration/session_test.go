// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"testing"

	"gitmin.com/gitmin/app/models/auth"
	"gitmin.com/gitmin/app/models/db"
	"gitmin.com/gitmin/app/models/unittest"
	"gitmin.com/gitmin/tests"

	"github.com/stretchr/testify/assert"
)

func Test_RegenerateSession(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	assert.NoError(t, unittest.PrepareTestDatabase())

	key := "new_key890123456"  // it must be 16 characters long
	key2 := "new_key890123457" // it must be 16 characters
	exist, err := auth.ExistSession(db.DefaultContext, key)
	assert.NoError(t, err)
	assert.False(t, exist)

	sess, err := auth.RegenerateSession(db.DefaultContext, "", key)
	assert.NoError(t, err)
	assert.EqualValues(t, key, sess.Key)
	assert.Empty(t, sess.Data)

	sess, err = auth.ReadSession(db.DefaultContext, key2)
	assert.NoError(t, err)
	assert.EqualValues(t, key2, sess.Key)
	assert.Empty(t, sess.Data)
}

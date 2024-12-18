// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	"gitmin.com/gitmin/app/models/db"
	"gitmin.com/gitmin/app/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestRunnerToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	token := unittest.AssertExistsAndLoadBean(t, &ActionRunnerToken{ID: 3})
	expectedToken, err := GetLatestRunnerToken(db.DefaultContext, 1, 0)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedToken, token)
}

func TestNewRunnerToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	token, err := NewRunnerToken(db.DefaultContext, 1, 0)
	assert.NoError(t, err)
	expectedToken, err := GetLatestRunnerToken(db.DefaultContext, 1, 0)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedToken, token)
}

func TestUpdateRunnerToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	token := unittest.AssertExistsAndLoadBean(t, &ActionRunnerToken{ID: 3})
	token.IsActive = true
	assert.NoError(t, UpdateRunnerToken(db.DefaultContext, token))
	expectedToken, err := GetLatestRunnerToken(db.DefaultContext, 1, 0)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedToken, token)
}

// Copyright 2025 The Gitmin Authors. All rights reserved.
// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package migrations

import (
	"testing"

	"gitmin.com/gitmin/app/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	defer test.MockVariableValue(&preparedMigrations)()
	preparedMigrations = []*migration{
		{idNumber: 0},
		{idNumber: 1},
	}
	assert.EqualValues(t, 2, calcDBVersion(preparedMigrations))
	assert.EqualValues(t, 2, ExpectedDBVersion())

	assert.EqualValues(t, 1, migrationIDNumberToDBVersion(0))

	assert.EqualValues(t, []*migration{{idNumber: 0}, {idNumber: 1}}, getPendingMigrations(0, preparedMigrations))
	assert.EqualValues(t, []*migration{{idNumber: 1}}, getPendingMigrations(1, preparedMigrations))
	assert.EqualValues(t, []*migration{}, getPendingMigrations(2, preparedMigrations))
}

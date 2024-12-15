// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package activitypub

import (
	"testing"

	"code.gitea.io/gitea/app/models/db"
	"code.gitea.io/gitea/app/models/unittest"
	user_model "code.gitea.io/gitea/app/models/user"

	_ "code.gitea.io/gitea/app/models" // https://forum.gitea.com/t/testfixtures-could-not-clean-table-access-no-such-table-access/4137/4

	"github.com/stretchr/testify/assert"
)

func TestUserSettings(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	pub, priv, err := GetKeyPair(db.DefaultContext, user1)
	assert.NoError(t, err)
	pub1, err := GetPublicKey(db.DefaultContext, user1)
	assert.NoError(t, err)
	assert.Equal(t, pub, pub1)
	priv1, err := GetPrivateKey(db.DefaultContext, user1)
	assert.NoError(t, err)
	assert.Equal(t, priv, priv1)
}

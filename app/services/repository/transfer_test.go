// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repository

import (
	"sync"
	"testing"

	"gitmin.com/gitmin/app/models"
	activities_model "gitmin.com/gitmin/app/models/activities"
	"gitmin.com/gitmin/app/models/db"
	"gitmin.com/gitmin/app/models/organization"
	access_model "gitmin.com/gitmin/app/models/perm/access"
	repo_model "gitmin.com/gitmin/app/models/repo"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/util"
	"gitmin.com/gitmin/app/services/feed"
	notify_service "gitmin.com/gitmin/app/services/notify"

	"github.com/stretchr/testify/assert"
)

var notifySync sync.Once

func registerNotifier() {
	notifySync.Do(func() {
		notify_service.RegisterNotifier(feed.NewNotifier())
	})
}

func TestTransferOwnership(t *testing.T) {
	registerNotifier()

	assert.NoError(t, unittest.PrepareTestDatabase())

	doer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	repo.Owner = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
	assert.NoError(t, TransferOwnership(db.DefaultContext, doer, doer, repo, nil))

	transferredRepo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	assert.EqualValues(t, 2, transferredRepo.OwnerID)

	exist, err := util.IsExist(repo_model.RepoPath("org3", "repo3"))
	assert.NoError(t, err)
	assert.False(t, exist)
	exist, err = util.IsExist(repo_model.RepoPath("user2", "repo3"))
	assert.NoError(t, err)
	assert.True(t, exist)
	unittest.AssertExistsAndLoadBean(t, &activities_model.Action{
		OpType:    activities_model.ActionTransferRepo,
		ActUserID: 2,
		RepoID:    3,
		Content:   "org3/repo3",
	})

	unittest.CheckConsistencyFor(t, &repo_model.Repository{}, &user_model.User{}, &organization.Team{})
}

func TestStartRepositoryTransferSetPermission(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	doer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})
	recipient := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 5})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})
	repo.Owner = unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	hasAccess, err := access_model.HasAnyUnitAccess(db.DefaultContext, recipient.ID, repo)
	assert.NoError(t, err)
	assert.False(t, hasAccess)

	assert.NoError(t, StartRepositoryTransfer(db.DefaultContext, doer, recipient, repo, nil))

	hasAccess, err = access_model.HasAnyUnitAccess(db.DefaultContext, recipient.ID, repo)
	assert.NoError(t, err)
	assert.True(t, hasAccess)

	unittest.CheckConsistencyFor(t, &repo_model.Repository{}, &user_model.User{}, &organization.Team{})
}

func TestRepositoryTransfer(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	doer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 3})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 3})

	transfer, err := models.GetPendingRepositoryTransfer(db.DefaultContext, repo)
	assert.NoError(t, err)
	assert.NotNil(t, transfer)

	// Cancel transfer
	assert.NoError(t, CancelRepositoryTransfer(db.DefaultContext, repo))

	transfer, err = models.GetPendingRepositoryTransfer(db.DefaultContext, repo)
	assert.Error(t, err)
	assert.Nil(t, transfer)
	assert.True(t, models.IsErrNoPendingTransfer(err))

	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	assert.NoError(t, models.CreatePendingRepositoryTransfer(db.DefaultContext, doer, user2, repo.ID, nil))

	transfer, err = models.GetPendingRepositoryTransfer(db.DefaultContext, repo)
	assert.NoError(t, err)
	assert.NoError(t, transfer.LoadAttributes(db.DefaultContext))
	assert.Equal(t, "user2", transfer.Recipient.Name)

	org6 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// Only transfer can be started at any given time
	err = models.CreatePendingRepositoryTransfer(db.DefaultContext, doer, org6, repo.ID, nil)
	assert.Error(t, err)
	assert.True(t, models.IsErrRepoTransferInProgress(err))

	// Unknown user
	err = models.CreatePendingRepositoryTransfer(db.DefaultContext, doer, &user_model.User{ID: 1000, LowerName: "user1000"}, repo.ID, nil)
	assert.Error(t, err)

	// Cancel transfer
	assert.NoError(t, CancelRepositoryTransfer(db.DefaultContext, repo))
}

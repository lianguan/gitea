// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"testing"

	"gitmin.com/gitmin/app/models/db"
	issues_model "gitmin.com/gitmin/app/models/issues"
	pull_model "gitmin.com/gitmin/app/models/pull"
	repo_model "gitmin.com/gitmin/app/models/repo"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/private"
	repo_module "gitmin.com/gitmin/app/modules/repository"
	"gitmin.com/gitmin/app/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestHandlePullRequestMerging(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	pr, err := issues_model.GetUnmergedPullRequest(db.DefaultContext, 1, 1, "branch2", "master", issues_model.MergeRequestFlowGithub)
	assert.NoError(t, err)
	assert.NoError(t, pr.LoadBaseRepo(db.DefaultContext))

	user1 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

	err = pull_model.ScheduleAutoMerge(db.DefaultContext, user1, pr.ID, repo_model.MergeStyleSquash, "squash merge a pr")
	assert.NoError(t, err)

	autoMerge := unittest.AssertExistsAndLoadBean(t, &pull_model.AutoMerge{MergeRequestID: pr.ID})

	ctx, resp := contexttest.MockPrivateContext(t, "/")
	handlePullRequestMerging(ctx, &private.HookOptions{
		PullRequestID: pr.ID,
		UserID:        2,
	}, pr.BaseRepo.OwnerName, pr.BaseRepo.Name, []*repo_module.PushUpdateOptions{
		{NewCommitID: "01234567"},
	})
	assert.Empty(t, resp.Body.String())
	pr, err = issues_model.GetPullRequestByID(db.DefaultContext, pr.ID)
	assert.NoError(t, err)
	assert.True(t, pr.HasMerged)
	assert.EqualValues(t, "01234567", pr.MergedCommitID)

	unittest.AssertNotExistsBean(t, &pull_model.AutoMerge{ID: autoMerge.ID})
}

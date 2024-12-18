// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/db"
	issues_model "gitmin.com/gitmin/app/models/issues"
	"gitmin.com/gitmin/app/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestMergeRequestList_LoadAttributes(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	prs := []*issues_model.MergeRequest{
		unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{ID: 2}),
	}
	assert.NoError(t, issues_model.MergeRequestList(prs).LoadAttributes(db.DefaultContext))
	for _, pr := range prs {
		assert.NotNil(t, pr.Issue)
		assert.Equal(t, pr.IssueID, pr.Issue.ID)
	}

	assert.NoError(t, issues_model.MergeRequestList([]*issues_model.MergeRequest{}).LoadAttributes(db.DefaultContext))
}

func TestMergeRequestList_LoadReviewCommentsCounts(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	prs := []*issues_model.MergeRequest{
		unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{ID: 2}),
	}
	reviewComments, err := issues_model.MergeRequestList(prs).LoadReviewCommentsCounts(db.DefaultContext)
	assert.NoError(t, err)
	assert.Len(t, reviewComments, 2)
	for _, pr := range prs {
		assert.EqualValues(t, 1, reviewComments[pr.IssueID])
	}
}

func TestMergeRequestList_LoadReviews(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	prs := []*issues_model.MergeRequest{
		unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{ID: 1}),
		unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{ID: 2}),
	}
	reviewList, err := issues_model.MergeRequestList(prs).LoadReviews(db.DefaultContext)
	assert.NoError(t, err)
	// 1, 7, 8, 9, 10, 22
	assert.Len(t, reviewList, 6)
	assert.EqualValues(t, 1, reviewList[0].ID)
	assert.EqualValues(t, 7, reviewList[1].ID)
	assert.EqualValues(t, 8, reviewList[2].ID)
	assert.EqualValues(t, 9, reviewList[3].ID)
	assert.EqualValues(t, 10, reviewList[4].ID)
	assert.EqualValues(t, 22, reviewList[5].ID)
}

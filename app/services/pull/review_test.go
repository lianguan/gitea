// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package pull_test

import (
	"testing"

	"gitmin.com/gitmin/app/models/db"
	issues_model "gitmin.com/gitmin/app/models/issues"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	pull_service "gitmin.com/gitmin/app/services/pull"

	"github.com/stretchr/testify/assert"
)

func TestDismissReview(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	pull := unittest.AssertExistsAndLoadBean(t, &issues_model.MergeRequest{})
	assert.NoError(t, pull.LoadIssue(db.DefaultContext))
	issue := pull.Issue
	assert.NoError(t, issue.LoadRepo(db.DefaultContext))
	reviewer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	review, err := issues_model.CreateReview(db.DefaultContext, issues_model.CreateReviewOptions{
		Issue:    issue,
		Reviewer: reviewer,
		Type:     issues_model.ReviewTypeReject,
	})

	assert.NoError(t, err)
	issue.IsClosed = true
	pull.HasMerged = false
	assert.NoError(t, issues_model.UpdateIssueCols(db.DefaultContext, issue, "is_closed"))
	assert.NoError(t, pull.UpdateCols(db.DefaultContext, "has_merged"))
	_, err = pull_service.DismissReview(db.DefaultContext, review.ID, issue.RepoID, "", &user_model.User{}, false, false)
	assert.Error(t, err)
	assert.True(t, pull_service.IsErrDismissRequestOnClosedPR(err))

	pull.HasMerged = true
	pull.Issue.IsClosed = false
	assert.NoError(t, issues_model.UpdateIssueCols(db.DefaultContext, issue, "is_closed"))
	assert.NoError(t, pull.UpdateCols(db.DefaultContext, "has_merged"))
	_, err = pull_service.DismissReview(db.DefaultContext, review.ID, issue.RepoID, "", &user_model.User{}, false, false)
	assert.Error(t, err)
	assert.True(t, pull_service.IsErrDismissRequestOnClosedPR(err))
}

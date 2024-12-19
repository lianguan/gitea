// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	"gitmin.com/gitmin/app/modules/git"
	api "gitmin.com/gitmin/app/modules/structs"
	webhook_module "gitmin.com/gitmin/app/modules/webhook"

	"github.com/stretchr/testify/assert"
)

func TestDetectMatched(t *testing.T) {
	testCases := []struct {
		desc         string
		commit       *git.Commit
		triggedEvent webhook_module.HookEventType
		payload      api.Payloader
		yamlOn       string
		expected     bool
	}{
		{
			desc:         "HookEventCreate(create) matches GithubEventCreate(create)",
			triggedEvent: webhook_module.HookEventCreate,
			payload:      nil,
			yamlOn:       "on: create",
			expected:     true,
		},
		{
			desc:         "HookEventIssues(issues) `opened` action matches GithubEventIssues(issues)",
			triggedEvent: webhook_module.HookEventIssues,
			payload:      &api.IssuePayload{Action: api.HookIssueOpened},
			yamlOn:       "on: issues",
			expected:     true,
		},
		{
			desc:         "HookEventIssues(issues) `milestoned` action matches GithubEventIssues(issues)",
			triggedEvent: webhook_module.HookEventIssues,
			payload:      &api.IssuePayload{Action: api.HookIssueMilestoned},
			yamlOn:       "on: issues",
			expected:     true,
		},
		{
			desc:         "HookEventPullRequestSync(merge_request_sync) matches GithubEventPullRequest(merge_request)",
			triggedEvent: webhook_module.HookEventPullRequestSync,
			payload:      &api.MergeRequestPayload{Action: api.HookIssueSynchronized},
			yamlOn:       "on: merge_request",
			expected:     true,
		},
		{
			desc:         "HookEventPullRequest(merge_request) `label_updated` action doesn't match GithubEventPullRequest(merge_request) with no activity type",
			triggedEvent: webhook_module.HookEventPullRequest,
			payload:      &api.MergeRequestPayload{Action: api.HookIssueLabelUpdated},
			yamlOn:       "on: merge_request",
			expected:     false,
		},
		{
			desc:         "HookEventPullRequest(merge_request) `closed` action doesn't match GithubEventPullRequest(merge_request) with no activity type",
			triggedEvent: webhook_module.HookEventPullRequest,
			payload:      &api.MergeRequestPayload{Action: api.HookIssueClosed},
			yamlOn:       "on: merge_request",
			expected:     false,
		},
		{
			desc:         "HookEventPullRequest(merge_request) `closed` action doesn't match GithubEventPullRequest(merge_request) with branches",
			triggedEvent: webhook_module.HookEventPullRequest,
			payload: &api.MergeRequestPayload{
				Action: api.HookIssueClosed,
				MergeRequest: &api.MergeRequest{
					Base: &api.PRBranchInfo{},
				},
			},
			yamlOn:   "on:\n  merge_request:\n    branches: [main]",
			expected: false,
		},
		{
			desc:         "HookEventPullRequest(merge_request) `label_updated` action matches GithubEventPullRequest(merge_request) with `label` activity type",
			triggedEvent: webhook_module.HookEventPullRequest,
			payload:      &api.MergeRequestPayload{Action: api.HookIssueLabelUpdated},
			yamlOn:       "on:\n  merge_request:\n    types: [labeled]",
			expected:     true,
		},
		{
			desc:         "HookEventPullRequestReviewComment(merge_request_review_comment) matches GithubEventPullRequestReviewComment(merge_request_review_comment)",
			triggedEvent: webhook_module.HookEventPullRequestReviewComment,
			payload:      &api.MergeRequestPayload{Action: api.HookIssueReviewed},
			yamlOn:       "on:\n  merge_request_review_comment:\n    types: [created]",
			expected:     true,
		},
		{
			desc:         "HookEventPullRequestReviewRejected(merge_request_review_rejected) doesn't match GithubEventPullRequestReview(merge_request_review) with `dismissed` activity type (we don't support `dismissed` at present)",
			triggedEvent: webhook_module.HookEventPullRequestReviewRejected,
			payload:      &api.MergeRequestPayload{Action: api.HookIssueReviewed},
			yamlOn:       "on:\n  merge_request_review:\n    types: [dismissed]",
			expected:     false,
		},
		{
			desc:         "HookEventRelease(release) `published` action matches GithubEventRelease(release) with `published` activity type",
			triggedEvent: webhook_module.HookEventRelease,
			payload:      &api.ReleasePayload{Action: api.HookReleasePublished},
			yamlOn:       "on:\n  release:\n    types: [published]",
			expected:     true,
		},
		{
			desc:         "HookEventPackage(package) `created` action doesn't match GithubEventRegistryPackage(registry_package) with `updated` activity type",
			triggedEvent: webhook_module.HookEventPackage,
			payload:      &api.PackagePayload{Action: api.HookPackageCreated},
			yamlOn:       "on:\n  registry_package:\n    types: [updated]",
			expected:     false,
		},
		{
			desc:         "HookEventWiki(wiki) matches GithubEventGollum(gollum)",
			triggedEvent: webhook_module.HookEventWiki,
			payload:      nil,
			yamlOn:       "on: gollum",
			expected:     true,
		},
		{
			desc:         "HookEventSchedue(schedule) matches GithubEventSchedule(schedule)",
			triggedEvent: webhook_module.HookEventSchedule,
			payload:      nil,
			yamlOn:       "on: schedule",
			expected:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			evts, err := GetEventsFromContent([]byte(tc.yamlOn))
			assert.NoError(t, err)
			assert.Len(t, evts, 1)
			assert.Equal(t, tc.expected, detectMatched(nil, tc.commit, tc.triggedEvent, tc.payload, evts[0]))
		})
	}
}

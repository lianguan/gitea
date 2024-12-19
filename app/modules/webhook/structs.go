// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

// HookEvents is a set of web hook events
type HookEvents struct {
	Create                   bool `json:"create"`
	Delete                   bool `json:"delete"`
	Fork                     bool `json:"fork"`
	Issues                   bool `json:"issues"`
	IssueAssign              bool `json:"issue_assign"`
	IssueLabel               bool `json:"issue_label"`
	IssueMilestone           bool `json:"issue_milestone"`
	IssueComment             bool `json:"issue_comment"`
	Push                     bool `json:"push"`
	PullRequest              bool `json:"merge_request"`
	PullRequestAssign        bool `json:"merge_request_assign"`
	PullRequestLabel         bool `json:"merge_request_label"`
	PullRequestMilestone     bool `json:"merge_request_milestone"`
	PullRequestComment       bool `json:"merge_request_comment"`
	PullRequestReview        bool `json:"merge_request_review"`
	PullRequestSync          bool `json:"merge_request_sync"`
	PullRequestReviewRequest bool `json:"merge_request_review_request"`
	Wiki                     bool `json:"wiki"`
	Repository               bool `json:"repository"`
	Release                  bool `json:"release"`
	Package                  bool `json:"package"`
}

// HookEvent represents events that will delivery hook.
type HookEvent struct {
	PushOnly       bool   `json:"push_only"`
	SendEverything bool   `json:"send_everything"`
	ChooseEvents   bool   `json:"choose_events"`
	BranchFilter   string `json:"branch_filter"`

	HookEvents `json:"events"`
}

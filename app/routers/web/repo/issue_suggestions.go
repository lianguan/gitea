// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"net/http"

	"gitmin.com/gitmin/app/models/db"
	issues_model "gitmin.com/gitmin/app/models/issues"
	"gitmin.com/gitmin/app/models/unit"
	issue_indexer "gitmin.com/gitmin/app/modules/indexer/issues"
	"gitmin.com/gitmin/app/modules/optional"
	"gitmin.com/gitmin/app/modules/structs"
	"gitmin.com/gitmin/app/services/context"
)

// IssueSuggestions returns a list of issue suggestions
func IssueSuggestions(ctx *context.Context) {
	keyword := ctx.Req.FormValue("q")

	canReadIssues := ctx.Repo.CanRead(unit.TypeIssues)
	canReadPulls := ctx.Repo.CanRead(unit.TypeMergeRequests)

	var isMergeRequest optional.Option[bool]
	if canReadPulls && !canReadIssues {
		isMergeRequest = optional.Some(true)
	} else if canReadIssues && !canReadPulls {
		isMergeRequest = optional.Some(false)
	}

	searchOpt := &issue_indexer.SearchOptions{
		Paginator: &db.ListOptions{
			Page:     0,
			PageSize: 5,
		},
		Keyword:        keyword,
		RepoIDs:        []int64{ctx.Repo.Repository.ID},
		IsMergeRequest: isMergeRequest,
		IsClosed:       nil,
		SortBy:         issue_indexer.SortByUpdatedDesc,
	}

	ids, _, err := issue_indexer.SearchIssues(ctx, searchOpt)
	if err != nil {
		ctx.ServerError("SearchIssues", err)
		return
	}
	issues, err := issues_model.GetIssuesByIDs(ctx, ids, true)
	if err != nil {
		ctx.ServerError("FindIssuesByIDs", err)
		return
	}

	suggestions := make([]*structs.Issue, 0, len(issues))

	for _, issue := range issues {
		suggestion := &structs.Issue{
			ID:    issue.ID,
			Index: issue.Index,
			Title: issue.Title,
			State: issue.State(),
		}

		if issue.IsMergeRequest {
			if err := issue.LoadPullRequest(ctx); err != nil {
				ctx.ServerError("LoadPullRequest", err)
				return
			}
			if issue.MergeRequest != nil {
				suggestion.MergeRequest = &structs.MergeRequestMeta{
					HasMerged:        issue.MergeRequest.HasMerged,
					IsWorkInProgress: issue.MergeRequest.IsWorkInProgress(ctx),
				}
			}
		}

		suggestions = append(suggestions, suggestion)
	}

	ctx.JSON(http.StatusOK, suggestions)
}

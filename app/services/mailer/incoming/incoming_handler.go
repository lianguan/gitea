// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package incoming

import (
	"bytes"
	"context"
	"fmt"

	issues_model "gitmin.com/gitmin/app/models/issues"
	access_model "gitmin.com/gitmin/app/models/perm/access"
	repo_model "gitmin.com/gitmin/app/models/repo"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/log"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/app/modules/util"
	attachment_service "gitmin.com/gitmin/app/services/attachment"
	"gitmin.com/gitmin/app/services/context/upload"
	issue_service "gitmin.com/gitmin/app/services/issue"
	incoming_payload "gitmin.com/gitmin/app/services/mailer/incoming/payload"
	"gitmin.com/gitmin/app/services/mailer/token"
	pull_service "gitmin.com/gitmin/app/services/pull"
)

type MailHandler interface {
	Handle(ctx context.Context, content *MailContent, doer *user_model.User, payload []byte) error
}

var handlers = map[token.HandlerType]MailHandler{
	token.ReplyHandlerType:       &ReplyHandler{},
	token.UnsubscribeHandlerType: &UnsubscribeHandler{},
}

// ReplyHandler handles incoming emails to create a reply from them
type ReplyHandler struct{}

func (h *ReplyHandler) Handle(ctx context.Context, content *MailContent, doer *user_model.User, payload []byte) error {
	if doer == nil {
		return util.NewInvalidArgumentErrorf("doer can't be nil")
	}

	ref, err := incoming_payload.GetReferenceFromPayload(ctx, payload)
	if err != nil {
		return err
	}

	var issue *issues_model.Issue

	switch r := ref.(type) {
	case *issues_model.Issue:
		issue = r
	case *issues_model.Comment:
		comment := r

		if err := comment.LoadIssue(ctx); err != nil {
			return err
		}

		issue = comment.Issue
	default:
		return util.NewInvalidArgumentErrorf("unsupported reply reference: %v", ref)
	}

	if err := issue.LoadRepo(ctx); err != nil {
		return err
	}

	perm, err := access_model.GetUserRepoPermission(ctx, issue.Repo, doer)
	if err != nil {
		return err
	}

	// Locked issues require write permissions
	if issue.IsLocked && !perm.CanWriteIssuesOrPulls(issue.IsMergeRequest) && !doer.IsAdmin {
		log.Debug("can't write issue or pull")
		return nil
	}

	if !perm.CanReadIssuesOrPulls(issue.IsMergeRequest) {
		log.Debug("can't read issue or pull")
		return nil
	}

	attachmentIDs := make([]string, 0, len(content.Attachments))
	if setting.Attachment.Enabled {
		for _, attachment := range content.Attachments {
			a, err := attachment_service.UploadAttachment(ctx, bytes.NewReader(attachment.Content), setting.Attachment.AllowedTypes, int64(len(attachment.Content)), &repo_model.Attachment{
				Name:       attachment.Name,
				UploaderID: doer.ID,
				RepoID:     issue.Repo.ID,
			})
			if err != nil {
				if upload.IsErrFileTypeForbidden(err) {
					log.Info("Skipping disallowed attachment type: %s", attachment.Name)
					continue
				}
				return err
			}
			attachmentIDs = append(attachmentIDs, a.UUID)
		}
	}

	if content.Content == "" && len(attachmentIDs) == 0 {
		return nil
	}

	switch r := ref.(type) {
	case *issues_model.Issue:
		_, err := issue_service.CreateIssueComment(ctx, doer, issue.Repo, issue, content.Content, attachmentIDs)
		if err != nil {
			return fmt.Errorf("CreateIssueComment failed: %w", err)
		}
	case *issues_model.Comment:
		comment := r

		switch comment.Type {
		case issues_model.CommentTypeCode:
			_, err := pull_service.CreateCodeComment(
				ctx,
				doer,
				nil,
				issue,
				comment.Line,
				content.Content,
				comment.TreePath,
				false, // not pending review but a single review
				comment.ReviewID,
				"",
				attachmentIDs,
			)
			if err != nil {
				return fmt.Errorf("CreateCodeComment failed: %w", err)
			}
		default:
			_, err := issue_service.CreateIssueComment(ctx, doer, issue.Repo, issue, content.Content, attachmentIDs)
			if err != nil {
				return fmt.Errorf("CreateIssueComment failed: %w", err)
			}
		}
	}
	return nil
}

// UnsubscribeHandler handles unwatching issues/pulls
type UnsubscribeHandler struct{}

func (h *UnsubscribeHandler) Handle(ctx context.Context, _ *MailContent, doer *user_model.User, payload []byte) error {
	if doer == nil {
		return util.NewInvalidArgumentErrorf("doer can't be nil")
	}

	ref, err := incoming_payload.GetReferenceFromPayload(ctx, payload)
	if err != nil {
		return err
	}

	switch r := ref.(type) {
	case *issues_model.Issue:
		issue := r

		if err := issue.LoadRepo(ctx); err != nil {
			return err
		}

		perm, err := access_model.GetUserRepoPermission(ctx, issue.Repo, doer)
		if err != nil {
			return err
		}

		if !perm.CanReadIssuesOrPulls(issue.IsMergeRequest) {
			log.Debug("can't read issue or pull")
			return nil
		}

		return issues_model.CreateOrUpdateIssueWatch(ctx, doer.ID, issue.ID, false)
	}

	return fmt.Errorf("unsupported unsubscribe reference: %v", ref)
}

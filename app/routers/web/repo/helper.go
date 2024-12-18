// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"net/url"

	"gitmin.com/gitmin/app/modules/git"
	"gitmin.com/gitmin/app/services/context"
)

func HandleGitError(ctx *context.Context, msg string, err error) {
	if git.IsErrNotExist(err) {
		refType := ""
		switch {
		case ctx.Repo.IsViewBranch:
			refType = "branch"
		case ctx.Repo.IsViewTag:
			refType = "tag"
		case ctx.Repo.IsViewCommit:
			refType = "commit"
		}
		ctx.Data["NotFoundPrompt"] = ctx.Locale.Tr("repo.tree_path_not_found_"+refType, ctx.Repo.TreePath, url.PathEscape(ctx.Repo.RefName))
		ctx.Data["NotFoundGoBackURL"] = ctx.Repo.RepoLink + "/src/" + refType + "/" + url.PathEscape(ctx.Repo.RefName)
		ctx.NotFound(msg, err)
	} else {
		ctx.ServerError(msg, err)
	}
}

// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"path/filepath"

	repo_model "gitmin.com/gitmin/app/models/repo"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/app/modules/util"
	"gitmin.com/gitmin/app/services/context"
	repo_service "gitmin.com/gitmin/app/services/repository"
)

// AdoptOrDeleteRepository adopts or deletes a repository
func AdoptOrDeleteRepository(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings.adopt")
	ctx.Data["PageIsSettingsRepos"] = true
	allowAdopt := ctx.IsUserSiteAdmin() || setting.Repository.AllowAdoptionOfUnadoptedRepositories
	ctx.Data["allowAdopt"] = allowAdopt
	allowDelete := ctx.IsUserSiteAdmin() || setting.Repository.AllowDeleteOfUnadoptedRepositories
	ctx.Data["allowDelete"] = allowDelete

	dir := ctx.FormString("id")
	action := ctx.FormString("action")

	ctxUser := ctx.Doer
	root := user_model.UserPath(ctxUser.LowerName)

	// check not a repo
	has, err := repo_model.IsRepositoryModelExist(ctx, ctxUser, dir)
	if err != nil {
		ctx.ServerError("IsRepositoryExist", err)
		return
	}

	isDir, err := util.IsDir(filepath.Join(root, dir+".git"))
	if err != nil {
		ctx.ServerError("IsDir", err)
		return
	}
	if has || !isDir {
		// Fallthrough to failure mode
	} else if action == "adopt" && allowAdopt {
		if _, err := repo_service.AdoptRepository(ctx, ctxUser, ctxUser, repo_service.CreateRepoOptions{
			Name:      dir,
			IsPrivate: true,
		}); err != nil {
			ctx.ServerError("repository.AdoptRepository", err)
			return
		}
		ctx.Flash.Success(ctx.Tr("repo.adopt_preexisting_success", dir))
	} else if action == "delete" && allowDelete {
		if err := repo_service.DeleteUnadoptedRepository(ctx, ctxUser, ctxUser, dir); err != nil {
			ctx.ServerError("repository.AdoptRepository", err)
			return
		}
		ctx.Flash.Success(ctx.Tr("repo.delete_preexisting_success", dir))
	}

	ctx.Redirect(setting.AppSubURL + "/user/settings/repos")
}

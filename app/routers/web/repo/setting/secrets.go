// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"errors"
	"net/http"

	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/base"
	"gitmin.com/gitmin/app/modules/setting"
	shared "gitmin.com/gitmin/app/routers/web/shared/secrets"
	shared_user "gitmin.com/gitmin/app/routers/web/shared/user"
	"gitmin.com/gitmin/app/services/context"
)

const (
	// TODO: Separate secrets from runners when layout is ready
	tplRepoSecrets base.TplName = "repo/settings/actions"
	tplOrgSecrets  base.TplName = "org/settings/actions"
	tplUserSecrets base.TplName = "user/settings/actions"
)

type secretsCtx struct {
	OwnerID         int64
	RepoID          int64
	IsRepo          bool
	IsOrg           bool
	IsUser          bool
	SecretsTemplate base.TplName
	RedirectLink    string
}

func getSecretsCtx(ctx *context.Context) (*secretsCtx, error) {
	if ctx.Data["PageIsRepoSettings"] == true {
		return &secretsCtx{
			OwnerID:         0,
			RepoID:          ctx.Repo.Repository.ID,
			IsRepo:          true,
			SecretsTemplate: tplRepoSecrets,
			RedirectLink:    ctx.Repo.RepoLink + "/settings/actions/secrets",
		}, nil
	}

	if ctx.Data["PageIsOrgSettings"] == true {
		err := shared_user.LoadHeaderCount(ctx)
		if err != nil {
			ctx.ServerError("LoadHeaderCount", err)
			return nil, nil
		}
		return &secretsCtx{
			OwnerID:         ctx.ContextUser.ID,
			RepoID:          0,
			IsOrg:           true,
			SecretsTemplate: tplOrgSecrets,
			RedirectLink:    ctx.Org.OrgLink + "/settings/actions/secrets",
		}, nil
	}

	if ctx.Data["PageIsUserSettings"] == true {
		return &secretsCtx{
			OwnerID:         ctx.Doer.ID,
			RepoID:          0,
			IsUser:          true,
			SecretsTemplate: tplUserSecrets,
			RedirectLink:    setting.AppSubURL + "/user/settings/actions/secrets",
		}, nil
	}

	return nil, errors.New("unable to set Secrets context")
}

func Secrets(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("actions.actions")
	ctx.Data["PageType"] = "secrets"
	ctx.Data["PageIsSharedSettingsSecrets"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	sCtx, err := getSecretsCtx(ctx)
	if err != nil {
		ctx.ServerError("getSecretsCtx", err)
		return
	}

	if sCtx.IsRepo {
		ctx.Data["DisableSSH"] = setting.SSH.Disabled
	}

	shared.SetSecretsContext(ctx, sCtx.OwnerID, sCtx.RepoID)
	if ctx.Written() {
		return
	}
	ctx.HTML(http.StatusOK, sCtx.SecretsTemplate)
}

func SecretsPost(ctx *context.Context) {
	sCtx, err := getSecretsCtx(ctx)
	if err != nil {
		ctx.ServerError("getSecretsCtx", err)
		return
	}

	if ctx.HasError() {
		ctx.JSONError(ctx.GetErrMsg())
		return
	}

	shared.PerformSecretsPost(
		ctx,
		sCtx.OwnerID,
		sCtx.RepoID,
		sCtx.RedirectLink,
	)
}

func SecretsDelete(ctx *context.Context) {
	sCtx, err := getSecretsCtx(ctx)
	if err != nil {
		ctx.ServerError("getSecretsCtx", err)
		return
	}
	shared.PerformSecretsDelete(
		ctx,
		sCtx.OwnerID,
		sCtx.RepoID,
		sCtx.RedirectLink,
	)
}

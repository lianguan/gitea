// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package org

import (
	"net/http"

	"gitmin.com/gitmin/app/modules/base"
	shared_user "gitmin.com/gitmin/app/routers/web/shared/user"
	"gitmin.com/gitmin/app/services/context"
)

const (
	tplSettingsBlockedUsers base.TplName = "org/settings/blocked_users"
)

func BlockedUsers(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("user.block.list")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsBlockedUsers"] = true

	shared_user.BlockedUsers(ctx, ctx.ContextUser)
	if ctx.Written() {
		return
	}

	ctx.HTML(http.StatusOK, tplSettingsBlockedUsers)
}

func BlockedUsersPost(ctx *context.Context) {
	shared_user.BlockedUsersPost(ctx, ctx.ContextUser)
	if ctx.Written() {
		return
	}

	ctx.Redirect(ctx.ContextUser.OrganisationLink() + "/settings/blocked_users")
}

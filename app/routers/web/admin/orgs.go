// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors.
// SPDX-License-Identifier: MIT

package admin

import (
	"code.gitea.io/gitea/app/models/db"
	user_model "code.gitea.io/gitea/app/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/app/routers/web/explore"
	"code.gitea.io/gitea/app/services/context"
)

const (
	tplOrgs base.TplName = "admin/org/list"
)

// Organizations show all the organizations
func Organizations(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.organizations")
	ctx.Data["PageIsAdminOrganizations"] = true

	if ctx.FormString("sort") == "" {
		ctx.SetFormString("sort", UserSearchDefaultAdminSort)
	}

	explore.RenderUserSearch(ctx, &user_model.SearchUserOptions{
		Actor:           ctx.Doer,
		Type:            user_model.UserTypeOrganization,
		IncludeReserved: true, // administrator needs to list all accounts include reserved
		ListOptions: db.ListOptions{
			PageSize: setting.UI.Admin.OrgPagingNum,
		},
		Visible: []structs.VisibleType{structs.VisibleTypePublic, structs.VisibleTypeLimited, structs.VisibleTypePrivate},
	}, tplOrgs)
}

// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	api "gitmin.com/gitmin/app/modules/structs"
	"gitmin.com/gitmin/app/modules/util"
	"gitmin.com/gitmin/app/modules/web"
	"gitmin.com/gitmin/app/routers/common"
	"gitmin.com/gitmin/app/services/context"
)

// Markup render markup document to HTML
func Markup(ctx *context.Context) {
	form := web.GetForm(ctx).(*api.MarkupOption)
	mode := util.Iif(form.Wiki, "wiki", form.Mode) //nolint:staticcheck
	common.RenderMarkup(ctx.Base, ctx.Repo, mode, form.Text, form.Context, form.FilePath)
}

// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package utils

import (
	"code.gitea.io/gitea/app/models/db"
	"code.gitea.io/gitea/app/services/context"
	"code.gitea.io/gitea/app/services/convert"
)

// GetListOptions returns list options using the page and limit parameters
func GetListOptions(ctx *context.APIContext) db.ListOptions {
	return db.ListOptions{
		Page:     ctx.FormInt("page"),
		PageSize: convert.ToCorrectPageSize(ctx.FormInt("limit")),
	}
}

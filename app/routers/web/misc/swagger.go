// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	"net/http"

	"gitmin.com/gitmin/app/modules/base"
	"gitmin.com/gitmin/app/services/context"
)

// tplSwagger swagger page template
const tplSwagger base.TplName = "swagger/ui"

// Swagger render swagger-ui page with v1 json
func Swagger(ctx *context.Context) {
	ctx.Data["APIJSONVersion"] = "v1"
	ctx.HTML(http.StatusOK, tplSwagger)
}

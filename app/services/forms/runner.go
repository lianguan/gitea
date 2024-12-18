// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package forms

import (
	"net/http"

	"gitmin.com/gitmin/app/modules/web/middleware"
	"gitmin.com/gitmin/app/services/context"

	"gitea.com/go-chi/binding"
)

// EditRunnerForm form for admin to create runner
type EditRunnerForm struct {
	Description string
}

// Validate validates form fields
func (f *EditRunnerForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, f, ctx.Locale)
}

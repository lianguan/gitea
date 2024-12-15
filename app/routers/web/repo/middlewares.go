// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"fmt"
	"strconv"

	system_model "code.gitea.io/gitea/app/models/system"
	user_model "code.gitea.io/gitea/app/models/user"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/app/services/context"
	user_service "code.gitea.io/gitea/app/services/user"
)

// SetEditorconfigIfExists set editor config as render variable
func SetEditorconfigIfExists(ctx *context.Context) {
	if ctx.Repo.Repository.IsEmpty {
		return
	}

	ec, _, err := ctx.Repo.GetEditorconfig()

	if err != nil && !git.IsErrNotExist(err) {
		description := fmt.Sprintf("Error while getting .editorconfig file: %v", err)
		if err := system_model.CreateRepositoryNotice(description); err != nil {
			ctx.ServerError("ErrCreatingReporitoryNotice", err)
		}
		return
	}

	ctx.Data["Editorconfig"] = ec
}

// SetDiffViewStyle set diff style as render variable
func SetDiffViewStyle(ctx *context.Context) {
	queryStyle := ctx.FormString("style")

	if !ctx.IsSigned {
		ctx.Data["IsSplitStyle"] = queryStyle == "split"
		return
	}

	var (
		userStyle = ctx.Doer.DiffViewStyle
		style     string
	)

	if queryStyle == "unified" || queryStyle == "split" {
		style = queryStyle
	} else if userStyle == "unified" || userStyle == "split" {
		style = userStyle
	} else {
		style = "unified"
	}

	ctx.Data["IsSplitStyle"] = style == "split"

	opts := &user_service.UpdateOptions{
		DiffViewStyle: optional.Some(style),
	}
	if err := user_service.UpdateUser(ctx, ctx.Doer, opts); err != nil {
		ctx.ServerError("UpdateUser", err)
	}
}

// SetWhitespaceBehavior set whitespace behavior as render variable
func SetWhitespaceBehavior(ctx *context.Context) {
	const defaultWhitespaceBehavior = "show-all"
	whitespaceBehavior := ctx.FormString("whitespace")
	switch whitespaceBehavior {
	case "", "ignore-all", "ignore-eol", "ignore-change":
		break
	default:
		whitespaceBehavior = defaultWhitespaceBehavior
	}
	if ctx.IsSigned {
		userWhitespaceBehavior, err := user_model.GetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyDiffWhitespaceBehavior, defaultWhitespaceBehavior)
		if err == nil {
			if whitespaceBehavior == "" {
				whitespaceBehavior = userWhitespaceBehavior
			} else if whitespaceBehavior != userWhitespaceBehavior {
				_ = user_model.SetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyDiffWhitespaceBehavior, whitespaceBehavior)
			}
		} // else: we can ignore the error safely
	}

	// these behaviors are for gitdiff.GetWhitespaceFlag
	if whitespaceBehavior == "" {
		ctx.Data["WhitespaceBehavior"] = defaultWhitespaceBehavior
	} else {
		ctx.Data["WhitespaceBehavior"] = whitespaceBehavior
	}
}

// SetShowOutdatedComments set the show outdated comments option as context variable
func SetShowOutdatedComments(ctx *context.Context) {
	showOutdatedCommentsValue := ctx.FormString("show-outdated")
	if showOutdatedCommentsValue != "true" && showOutdatedCommentsValue != "false" {
		// invalid or no value for this form string -> use default or stored user setting
		showOutdatedCommentsValue = "true"
		if ctx.IsSigned {
			showOutdatedCommentsValue, _ = user_model.GetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyShowOutdatedComments, showOutdatedCommentsValue)
		}
	} else if ctx.IsSigned {
		// valid value -> update user setting if user is logged in
		_ = user_model.SetUserSetting(ctx, ctx.Doer.ID, user_model.SettingsKeyShowOutdatedComments, showOutdatedCommentsValue)
	}
	ctx.Data["ShowOutdatedComments"], _ = strconv.ParseBool(showOutdatedCommentsValue)
}
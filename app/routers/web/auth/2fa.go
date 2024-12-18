// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth

import (
	"errors"
	"net/http"

	"gitmin.com/gitmin/app/models/auth"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/base"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/app/modules/web"
	"gitmin.com/gitmin/app/services/context"
	"gitmin.com/gitmin/app/services/externalaccount"
	"gitmin.com/gitmin/app/services/forms"
)

var (
	tplTwofa        base.TplName = "user/auth/twofa"
	tplTwofaScratch base.TplName = "user/auth/twofa_scratch"
)

// TwoFactor shows the user a two-factor authentication page.
func TwoFactor(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("twofa")

	if CheckAutoLogin(ctx) {
		return
	}

	// Ensure user is in a 2FA session.
	if ctx.Session.Get("twofaUid") == nil {
		ctx.ServerError("UserSignIn", errors.New("not in 2FA session"))
		return
	}

	ctx.HTML(http.StatusOK, tplTwofa)
}

// TwoFactorPost validates a user's two-factor authentication token.
func TwoFactorPost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.TwoFactorAuthForm)
	ctx.Data["Title"] = ctx.Tr("twofa")

	// Ensure user is in a 2FA session.
	idSess := ctx.Session.Get("twofaUid")
	if idSess == nil {
		ctx.ServerError("UserSignIn", errors.New("not in 2FA session"))
		return
	}

	id := idSess.(int64)
	twofa, err := auth.GetTwoFactorByUID(ctx, id)
	if err != nil {
		ctx.ServerError("UserSignIn", err)
		return
	}

	// Validate the passcode with the stored TOTP secret.
	ok, err := twofa.ValidateTOTP(form.Passcode)
	if err != nil {
		ctx.ServerError("UserSignIn", err)
		return
	}

	if ok && twofa.LastUsedPasscode != form.Passcode {
		remember := ctx.Session.Get("twofaRemember").(bool)
		u, err := user_model.GetUserByID(ctx, id)
		if err != nil {
			ctx.ServerError("UserSignIn", err)
			return
		}

		if ctx.Session.Get("linkAccount") != nil {
			err = externalaccount.LinkAccountFromStore(ctx, ctx.Session, u)
			if err != nil {
				ctx.ServerError("UserSignIn", err)
				return
			}
		}

		twofa.LastUsedPasscode = form.Passcode
		if err = auth.UpdateTwoFactor(ctx, twofa); err != nil {
			ctx.ServerError("UserSignIn", err)
			return
		}

		handleSignIn(ctx, u, remember)
		return
	}

	ctx.RenderWithErr(ctx.Tr("auth.twofa_passcode_incorrect"), tplTwofa, forms.TwoFactorAuthForm{})
}

// TwoFactorScratch shows the scratch code form for two-factor authentication.
func TwoFactorScratch(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("twofa_scratch")

	if CheckAutoLogin(ctx) {
		return
	}

	// Ensure user is in a 2FA session.
	if ctx.Session.Get("twofaUid") == nil {
		ctx.ServerError("UserSignIn", errors.New("not in 2FA session"))
		return
	}

	ctx.HTML(http.StatusOK, tplTwofaScratch)
}

// TwoFactorScratchPost validates and invalidates a user's two-factor scratch token.
func TwoFactorScratchPost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.TwoFactorScratchAuthForm)
	ctx.Data["Title"] = ctx.Tr("twofa_scratch")

	// Ensure user is in a 2FA session.
	idSess := ctx.Session.Get("twofaUid")
	if idSess == nil {
		ctx.ServerError("UserSignIn", errors.New("not in 2FA session"))
		return
	}

	id := idSess.(int64)
	twofa, err := auth.GetTwoFactorByUID(ctx, id)
	if err != nil {
		ctx.ServerError("UserSignIn", err)
		return
	}

	// Validate the passcode with the stored TOTP secret.
	if twofa.VerifyScratchToken(form.Token) {
		// Invalidate the scratch token.
		_, err = twofa.GenerateScratchToken()
		if err != nil {
			ctx.ServerError("UserSignIn", err)
			return
		}
		if err = auth.UpdateTwoFactor(ctx, twofa); err != nil {
			ctx.ServerError("UserSignIn", err)
			return
		}

		remember := ctx.Session.Get("twofaRemember").(bool)
		u, err := user_model.GetUserByID(ctx, id)
		if err != nil {
			ctx.ServerError("UserSignIn", err)
			return
		}

		handleSignInFull(ctx, u, remember, false)
		if ctx.Written() {
			return
		}
		ctx.Flash.Info(ctx.Tr("auth.twofa_scratch_used"))
		ctx.Redirect(setting.AppSubURL + "/user/settings/security")
		return
	}

	ctx.RenderWithErr(ctx.Tr("auth.twofa_scratch_token_incorrect"), tplTwofaScratch, forms.TwoFactorScratchAuthForm{})
}

// Copyright 2016 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors.
// SPDX-License-Identifier: MIT

package repo

import (
	"errors"
	"fmt"
	"net/http"

	"gitmin.com/gitmin/app/models/organization"
	"gitmin.com/gitmin/app/models/perm"
	access_model "gitmin.com/gitmin/app/models/perm/access"
	repo_model "gitmin.com/gitmin/app/models/repo"
	user_model "gitmin.com/gitmin/app/models/user"
	api "gitmin.com/gitmin/app/modules/structs"
	"gitmin.com/gitmin/app/modules/util"
	"gitmin.com/gitmin/app/modules/web"
	"gitmin.com/gitmin/app/routers/api/v1/utils"
	"gitmin.com/gitmin/app/services/context"
	"gitmin.com/gitmin/app/services/convert"
	repo_service "gitmin.com/gitmin/app/services/repository"
)

// ListForks list a repository's forks
func ListForks(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/forks repository listForks
	// ---
	// summary: List a repository's forks
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/RepositoryList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	forks, total, err := repo_service.FindForks(ctx, ctx.Repo.Repository, ctx.Doer, utils.GetListOptions(ctx))
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "FindForks", err)
		return
	}
	if err := repo_model.RepositoryList(forks).LoadOwners(ctx); err != nil {
		ctx.Error(http.StatusInternalServerError, "LoadOwners", err)
		return
	}
	if err := repo_model.RepositoryList(forks).LoadUnits(ctx); err != nil {
		ctx.Error(http.StatusInternalServerError, "LoadUnits", err)
		return
	}

	apiForks := make([]*api.Repository, len(forks))
	for i, fork := range forks {
		permission, err := access_model.GetUserRepoPermission(ctx, fork, ctx.Doer)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, "GetUserRepoPermission", err)
			return
		}
		apiForks[i] = convert.ToRepo(ctx, fork, permission)
	}

	ctx.SetTotalCountHeader(total)
	ctx.JSON(http.StatusOK, apiForks)
}

// CreateFork create a fork of a repo
func CreateFork(ctx *context.APIContext) {
	// swagger:operation POST /repos/{owner}/{repo}/forks repository createFork
	// ---
	// summary: Fork a repository
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo to fork
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo to fork
	//   type: string
	//   required: true
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/CreateForkOption"
	// responses:
	//   "202":
	//     "$ref": "#/responses/Repository"
	//   "403":
	//     "$ref": "#/responses/forbidden"
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "409":
	//     description: The repository with the same name already exists.
	//   "422":
	//     "$ref": "#/responses/validationError"

	form := web.GetForm(ctx).(*api.CreateForkOption)
	repo := ctx.Repo.Repository
	var forker *user_model.User // user/org that will own the fork
	if form.Organization == nil {
		forker = ctx.Doer
	} else {
		org, err := organization.GetOrgByName(ctx, *form.Organization)
		if err != nil {
			if organization.IsErrOrgNotExist(err) {
				ctx.Error(http.StatusUnprocessableEntity, "", err)
			} else {
				ctx.Error(http.StatusInternalServerError, "GetOrgByName", err)
			}
			return
		}
		isMember, err := org.IsOrgMember(ctx, ctx.Doer.ID)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, "IsOrgMember", err)
			return
		} else if !isMember {
			ctx.Error(http.StatusForbidden, "isMemberNot", fmt.Sprintf("User is no Member of Organisation '%s'", org.Name))
			return
		}
		forker = org.AsUser()
	}

	var name string
	if form.Name == nil {
		name = repo.Name
	} else {
		name = *form.Name
	}

	fork, err := repo_service.ForkRepository(ctx, ctx.Doer, forker, repo_service.ForkRepoOptions{
		BaseRepo:    repo,
		Name:        name,
		Description: repo.Description,
	})
	if err != nil {
		if errors.Is(err, util.ErrAlreadyExist) || repo_model.IsErrReachLimitOfRepo(err) {
			ctx.Error(http.StatusConflict, "ForkRepository", err)
		} else if errors.Is(err, user_model.ErrBlockedUser) {
			ctx.Error(http.StatusForbidden, "ForkRepository", err)
		} else {
			ctx.Error(http.StatusInternalServerError, "ForkRepository", err)
		}
		return
	}

	// TODO change back to 201
	ctx.JSON(http.StatusAccepted, convert.ToRepo(ctx, fork, access_model.Permission{AccessMode: perm.AccessModeOwner}))
}

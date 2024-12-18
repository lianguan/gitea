// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repository

import (
	"context"

	"gitmin.com/gitmin/app/models/organization"
	repo_model "gitmin.com/gitmin/app/models/repo"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/setting"
)

func CanUserForkBetweenOwners(id1, id2 int64) bool {
	if id1 != id2 {
		return true
	}
	return setting.Repository.AllowForkIntoSameOwner
}

// CanUserForkRepo returns true if specified user can fork repository.
func CanUserForkRepo(ctx context.Context, user *user_model.User, repo *repo_model.Repository) (bool, error) {
	if user == nil {
		return false, nil
	}
	if CanUserForkBetweenOwners(repo.OwnerID, user.ID) && !repo_model.HasForkedRepo(ctx, user.ID, repo.ID) {
		return true, nil
	}
	ownedOrgs, err := organization.GetOrgsCanCreateRepoByUserID(ctx, user.ID)
	if err != nil {
		return false, err
	}
	for _, org := range ownedOrgs {
		if repo.OwnerID != org.ID && !repo_model.HasForkedRepo(ctx, org.ID, repo.ID) {
			return true, nil
		}
	}
	return false, nil
}

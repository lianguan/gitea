// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repository

import (
	"context"
	"slices"

	actions_model "gitmin.com/gitmin/app/models/actions"
	"gitmin.com/gitmin/app/models/db"
	repo_model "gitmin.com/gitmin/app/models/repo"
	"gitmin.com/gitmin/app/models/unit"
	"gitmin.com/gitmin/app/modules/log"
	actions_service "gitmin.com/gitmin/app/services/actions"
)

// UpdateRepositoryUnits updates a repository's units
func UpdateRepositoryUnits(ctx context.Context, repo *repo_model.Repository, units []repo_model.RepoUnit, deleteUnitTypes []unit.Type) (err error) {
	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return err
	}
	defer committer.Close()

	// Delete existing settings of units before adding again
	for _, u := range units {
		deleteUnitTypes = append(deleteUnitTypes, u.Type)
	}

	if slices.Contains(deleteUnitTypes, unit.TypeActions) {
		if err := actions_model.CleanRepoScheduleTasks(ctx, repo); err != nil {
			log.Error("CleanRepoScheduleTasks: %v", err)
		}
	}

	for _, u := range units {
		if u.Type == unit.TypeActions {
			if err := actions_service.DetectAndHandleSchedules(ctx, repo); err != nil {
				log.Error("DetectAndHandleSchedules: %v", err)
			}
			break
		}
	}

	if _, err = db.GetEngine(ctx).Where("repo_id = ?", repo.ID).In("type", deleteUnitTypes).Delete(new(repo_model.RepoUnit)); err != nil {
		return err
	}

	if len(units) > 0 {
		if err = db.Insert(ctx, units); err != nil {
			return err
		}
	}

	return committer.Commit()
}

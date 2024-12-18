// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package stats

import (
	"context"
	"testing"
	"time"

	"code.gitea.io/gitea/app/models/db"
	repo_model "code.gitea.io/gitea/app/models/repo"
	"code.gitea.io/gitea/app/models/unittest"
	"code.gitea.io/gitea/app/modules/queue"
	"code.gitea.io/gitea/app/modules/setting"

	_ "code.gitea.io/gitea/app/models"
	_ "code.gitea.io/gitea/app/models/actions"
	_ "code.gitea.io/gitea/app/models/activities"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

func TestRepoStatsIndex(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	setting.CfgProvider, _ = setting.NewConfigProviderFromData("")

	setting.LoadQueueSettings()

	err := Init()
	assert.NoError(t, err)

	repo, err := repo_model.GetRepositoryByID(db.DefaultContext, 1)
	assert.NoError(t, err)

	err = UpdateRepoIndexer(repo)
	assert.NoError(t, err)

	assert.NoError(t, queue.GetManager().FlushAll(context.Background(), 10*time.Second))

	status, err := repo_model.GetIndexerStatus(db.DefaultContext, repo, repo_model.RepoIndexerTypeStats)
	assert.NoError(t, err)
	assert.Equal(t, "65f1bf27bc3bf70f64657658635e66094edbcb4d", status.CommitSha)
	langs, err := repo_model.GetTopLanguageStats(db.DefaultContext, repo, 5)
	assert.NoError(t, err)
	assert.Empty(t, langs)
}

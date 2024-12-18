// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"sync"
	"testing"

	auth_model "gitmin.com/gitmin/app/models/auth"
	"gitmin.com/gitmin/app/models/db"
	issues_model "gitmin.com/gitmin/app/models/issues"
	repo_model "gitmin.com/gitmin/app/models/repo"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/git"
	"gitmin.com/gitmin/app/modules/gitrepo"
	files_service "gitmin.com/gitmin/app/services/repository/files"

	"github.com/stretchr/testify/assert"
)

func TestDataAsyncDoubleRead_Issue29101(t *testing.T) {
	onGiteaRun(t, func(t *testing.T, u *url.URL) {
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

		testContent := bytes.Repeat([]byte{'a'}, 10000)
		resp, err := files_service.ChangeRepoFiles(db.DefaultContext, repo, user, &files_service.ChangeRepoFilesOptions{
			Files: []*files_service.ChangeRepoFile{
				{
					Operation:     "create",
					TreePath:      "test.txt",
					ContentReader: bytes.NewReader(testContent),
				},
			},
			OldBranch: repo.DefaultBranch,
			NewBranch: repo.DefaultBranch,
		})
		assert.NoError(t, err)

		sha := resp.Commit.SHA

		gitRepo, err := gitrepo.OpenRepository(db.DefaultContext, repo)
		assert.NoError(t, err)

		commit, err := gitRepo.GetCommit(sha)
		assert.NoError(t, err)

		entry, err := commit.GetTreeEntryByPath("test.txt")
		assert.NoError(t, err)

		b := entry.Blob()
		r1, err := b.DataAsync()
		assert.NoError(t, err)
		defer r1.Close()
		r2, err := b.DataAsync()
		assert.NoError(t, err)
		defer r2.Close()

		var data1, data2 []byte
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			data1, _ = io.ReadAll(r1)
			assert.NoError(t, err)
			wg.Done()
		}()
		go func() {
			data2, _ = io.ReadAll(r2)
			assert.NoError(t, err)
			wg.Done()
		}()
		wg.Wait()
		assert.Equal(t, testContent, data1)
		assert.Equal(t, testContent, data2)
	})
}

func TestAgitPullPush(t *testing.T) {
	onGiteaRun(t, func(t *testing.T, u *url.URL) {
		baseAPITestContext := NewAPITestContext(t, "user2", "repo1", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)

		u.Path = baseAPITestContext.GitPath()
		u.User = url.UserPassword("user2", userPassword)

		dstPath := t.TempDir()
		doGitClone(dstPath, u)(t)

		gitRepo, err := git.OpenRepository(context.Background(), dstPath)
		assert.NoError(t, err)
		defer gitRepo.Close()

		doGitCreateBranch(dstPath, "test-agit-push")

		// commit 1
		_, err = generateCommitWithNewData(testFileSizeSmall, dstPath, "user2@example.com", "User Two", "branch-data-file-")
		assert.NoError(t, err)

		// push to create an agit pull request
		err = git.NewCommand(git.DefaultContext, "push", "origin",
			"-o", "title=test-title", "-o", "description=test-description",
			"HEAD:refs/for/master/test-agit-push",
		).Run(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)

		// check pull request exist
		pr := unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{BaseRepoID: 1, Flow: issues_model.PullRequestFlowAGit, HeadBranch: "user2/test-agit-push"})
		assert.NoError(t, pr.LoadIssue(db.DefaultContext))
		assert.Equal(t, "test-title", pr.Issue.Title)
		assert.Equal(t, "test-description", pr.Issue.Content)

		// commit 2
		_, err = generateCommitWithNewData(testFileSizeSmall, dstPath, "user2@example.com", "User Two", "branch-data-file-2-")
		assert.NoError(t, err)

		// push 2
		err = git.NewCommand(git.DefaultContext, "push", "origin", "HEAD:refs/for/master/test-agit-push").Run(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)

		// reset to first commit
		err = git.NewCommand(git.DefaultContext, "reset", "--hard", "HEAD~1").Run(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)

		// test force push without confirm
		_, stderr, err := git.NewCommand(git.DefaultContext, "push", "origin", "HEAD:refs/for/master/test-agit-push").RunStdString(&git.RunOpts{Dir: dstPath})
		assert.Error(t, err)
		assert.Contains(t, stderr, "[remote rejected] HEAD -> refs/for/master/test-agit-push (request `force-push` push option)")

		// test force push with confirm
		err = git.NewCommand(git.DefaultContext, "push", "origin", "HEAD:refs/for/master/test-agit-push", "-o", "force-push").Run(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)
	})
}

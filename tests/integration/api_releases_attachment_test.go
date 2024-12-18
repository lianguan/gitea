// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"fmt"
	"net/http"
	"testing"

	auth_model "gitmin.com/gitmin/app/models/auth"
	repo_model "gitmin.com/gitmin/app/models/repo"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/app/modules/test"
	"gitmin.com/gitmin/tests"
)

func TestAPIEditReleaseAttachmentWithUnallowedFile(t *testing.T) {
	// Limit the allowed release types (since by default there is no restriction)
	defer test.MockVariableValue(&setting.Repository.Release.AllowedTypes, ".exe")()
	defer tests.PrepareTestEnv(t)()

	attachment := unittest.AssertExistsAndLoadBean(t, &repo_model.Attachment{ID: 9})
	release := unittest.AssertExistsAndLoadBean(t, &repo_model.Release{ID: attachment.ReleaseID})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: attachment.RepoID})
	repoOwner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, repoOwner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

	filename := "file.bad"
	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/releases/%d/assets/%d", repoOwner.Name, repo.Name, release.ID, attachment.ID)
	req := NewRequestWithValues(t, "PATCH", urlStr, map[string]string{
		"name": filename,
	}).AddTokenAuth(token)

	session.MakeRequest(t, req, http.StatusUnprocessableEntity)
}

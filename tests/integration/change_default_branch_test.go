// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"fmt"
	"net/http"
	"testing"

	repo_model "gitmin.com/gitmin/app/models/repo"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/tests"
)

func TestChangeDefaultBranch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	branchesURL := fmt.Sprintf("/%s/%s/settings/branches", owner.Name, repo.Name)

	csrf := GetUserCSRFToken(t, session)
	req := NewRequestWithValues(t, "POST", branchesURL, map[string]string{
		"_csrf":  csrf,
		"action": "default_branch",
		"branch": "DefaultBranch",
	})
	session.MakeRequest(t, req, http.StatusSeeOther)

	csrf = GetUserCSRFToken(t, session)
	req = NewRequestWithValues(t, "POST", branchesURL, map[string]string{
		"_csrf":  csrf,
		"action": "default_branch",
		"branch": "does_not_exist",
	})
	session.MakeRequest(t, req, http.StatusNotFound)
}

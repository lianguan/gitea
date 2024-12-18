// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"
	"time"

	auth_model "gitmin.com/gitmin/app/models/auth"
	"gitmin.com/gitmin/app/models/db"
	"gitmin.com/gitmin/app/models/unittest"
	user_model "gitmin.com/gitmin/app/models/user"
	api "gitmin.com/gitmin/app/modules/structs"
	"gitmin.com/gitmin/app/services/convert"
	"gitmin.com/gitmin/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPITeamUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadOrganization)
	req := NewRequest(t, "GET", "/api/v1/teams/1/members/user1").
		AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)

	req = NewRequest(t, "GET", "/api/v1/teams/1/members/user2").
		AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)
	var user2 *api.User
	DecodeJSON(t, resp, &user2)
	user2.Created = user2.Created.In(time.Local)
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: "user2"})

	expectedUser := convert.ToUser(db.DefaultContext, user, user)

	// test time via unix timestamp
	assert.EqualValues(t, expectedUser.LastLogin.Unix(), user2.LastLogin.Unix())
	assert.EqualValues(t, expectedUser.Created.Unix(), user2.Created.Unix())
	expectedUser.LastLogin = user2.LastLogin
	expectedUser.Created = user2.Created

	assert.Equal(t, expectedUser, user2)
}

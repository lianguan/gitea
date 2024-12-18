// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package install

import (
	"net/http/httptest"
	"testing"

	"gitmin.com/gitmin/app/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	r := Routes()
	assert.NotNil(t, r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.EqualValues(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `class="page-content install"`)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/no-such", nil)
	r.ServeHTTP(w, req)
	assert.EqualValues(t, 404, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/assets/img/gitmin.svg", nil)
	r.ServeHTTP(w, req)
	assert.EqualValues(t, 200, w.Code)
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

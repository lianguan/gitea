// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	"gitmin.com/gitmin/app/models/unittest"
	"gitmin.com/gitmin/app/modules/hostmatcher"
	"gitmin.com/gitmin/app/modules/setting"

	_ "gitmin.com/gitmin/app/models"
	_ "gitmin.com/gitmin/app/models/actions"
)

func TestMain(m *testing.M) {
	// for tests, allow only loopback IPs
	setting.Webhook.AllowedHostList = hostmatcher.MatchBuiltinLoopback
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return Init()
		},
	})
}

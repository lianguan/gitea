// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"gitmin.com/gitmin/app/modules/graceful"
	"gitmin.com/gitmin/app/modules/log"
	"gitmin.com/gitmin/app/modules/queue"
	"gitmin.com/gitmin/app/modules/setting"
	notify_service "gitmin.com/gitmin/app/services/notify"
)

func Init() {
	if !setting.Actions.Enabled {
		return
	}

	jobEmitterQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "actions_ready_job", jobEmitterQueueHandler)
	if jobEmitterQueue == nil {
		log.Fatal("Unable to create actions_ready_job queue")
	}
	go graceful.GetManager().RunWithCancel(jobEmitterQueue)

	notify_service.RegisterNotifier(NewNotifier())
}

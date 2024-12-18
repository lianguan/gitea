// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_13 //nolint

import (
	"gitmin.com/gitmin/app/modules/timeutil"
	"gitmin.com/gitmin/database/migrations/base"

	"xorm.io/xorm"
)

func AddPrimaryKeyToRepoTopic(x *xorm.Engine) error {
	// Topic represents a topic of repositories
	type Topic struct {
		ID          int64  `xorm:"pk autoincr"`
		Name        string `xorm:"UNIQUE VARCHAR(25)"`
		RepoCount   int
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	// RepoTopic represents associated repositories and topics
	type RepoTopic struct {
		RepoID  int64 `xorm:"pk"`
		TopicID int64 `xorm:"pk"`
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := base.RecreateTable(sess, &Topic{}); err != nil {
		return err
	}
	if err := base.RecreateTable(sess, &RepoTopic{}); err != nil {
		return err
	}

	return sess.Commit()
}

// Copyright 2022 Gitea. All rights reserved.
// SPDX-License-Identifier: MIT

package pull

import (
	"context"
	"fmt"

	"gitmin.com/gitmin/app/models/db"
	repo_model "gitmin.com/gitmin/app/models/repo"
	user_model "gitmin.com/gitmin/app/models/user"
	"gitmin.com/gitmin/app/modules/timeutil"
)

// AutoMerge represents a pull request scheduled for merging when checks succeed
type AutoMerge struct {
	ID          int64                 `xorm:"pk autoincr"`
	MergeRequestID      int64                 `xorm:"UNIQUE"`
	DoerID      int64                 `xorm:"INDEX NOT NULL"`
	Doer        *user_model.User      `xorm:"-"`
	MergeStyle  repo_model.MergeStyle `xorm:"varchar(30)"`
	Message     string                `xorm:"LONGTEXT"`
	CreatedUnix timeutil.TimeStamp    `xorm:"created"`
}

// TableName return database table name for xorm
func (AutoMerge) TableName() string {
	return "pull_auto_merge"
}

func init() {
	db.RegisterModel(new(AutoMerge))
}

// ErrAlreadyScheduledToAutoMerge represents a "PullRequestHasMerged"-error
type ErrAlreadyScheduledToAutoMerge struct {
	MergeRequestID int64
}

func (err ErrAlreadyScheduledToAutoMerge) Error() string {
	return fmt.Sprintf("pull request is already scheduled to auto merge when checks succeed [merge_request_id: %d]", err.MergeRequestID)
}

// IsErrAlreadyScheduledToAutoMerge checks if an error is a ErrAlreadyScheduledToAutoMerge.
func IsErrAlreadyScheduledToAutoMerge(err error) bool {
	_, ok := err.(ErrAlreadyScheduledToAutoMerge)
	return ok
}

// ScheduleAutoMerge schedules a pull request to be merged when all checks succeed
func ScheduleAutoMerge(ctx context.Context, doer *user_model.User, mergeRequestID int64, style repo_model.MergeStyle, message string) error {
	// Check if we already have a merge scheduled for that pull request
	if exists, _, err := GetScheduledMergeByMergeRequestID(ctx, mergeRequestID); err != nil {
		return err
	} else if exists {
		return ErrAlreadyScheduledToAutoMerge{MergeRequestID: mergeRequestID}
	}

	_, err := db.GetEngine(ctx).Insert(&AutoMerge{
		DoerID:     doer.ID,
		MergeRequestID:     mergeRequestID,
		MergeStyle: style,
		Message:    message,
	})
	return err
}

// GetScheduledMergeByMergeRequestID gets a scheduled pull request merge by pull request id
func GetScheduledMergeByMergeRequestID(ctx context.Context, mergeRequestID int64) (bool, *AutoMerge, error) {
	scheduledPRM := &AutoMerge{}
	exists, err := db.GetEngine(ctx).Where("merge_request_id = ?", mergeRequestID).Get(scheduledPRM)
	if err != nil || !exists {
		return false, nil, err
	}

	doer, err := user_model.GetUserByID(ctx, scheduledPRM.DoerID)
	if err != nil {
		return false, nil, err
	}

	scheduledPRM.Doer = doer
	return true, scheduledPRM, nil
}

// DeleteScheduledAutoMerge delete a scheduled pull request
func DeleteScheduledAutoMerge(ctx context.Context, mergeRequestID int64) error {
	exist, scheduledPRM, err := GetScheduledMergeByMergeRequestID(ctx, mergeRequestID)
	if err != nil {
		return err
	} else if !exist {
		return db.ErrNotExist{Resource: "auto_merge", ID: mergeRequestID}
	}

	_, err = db.GetEngine(ctx).ID(scheduledPRM.ID).Delete(&AutoMerge{})
	return err
}

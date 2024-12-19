// Copyright 2015 The Gogs Authors. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"
	"fmt"

	"gitmin.com/gitmin/app/modules/git"
	"gitmin.com/gitmin/app/modules/log"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/database/migrations/v1_0"

	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const minDBVersion = 0 // Gitmin 1.0.0

type migration struct {
	idNumber    int64 // DB version is "the last migration's idNumber" + 1
	description string
	migrate     func(*xorm.Engine) error
}

// newMigration creates a new migration
func newMigration(idNumber int64, desc string, fn func(*xorm.Engine) error) *migration {
	return &migration{idNumber, desc, fn}
}

// Migrate executes the migration
func (m *migration) Migrate(x *xorm.Engine) error {
	return m.migrate(x)
}

// Version describes the version table. Should have only one row with id==1
type Version struct {
	ID      int64 `xorm:"pk autoincr"`
	Version int64 // DB version is "the last migration's idNumber" + 1
}

// Use noopMigration when there is a migration that has been no-oped
var noopMigration = func(_ *xorm.Engine) error { return nil }

var preparedMigrations []*migration

// This is a sequence of migrations. Add new migrations to the bottom of the list.
// If you want to "retire" a migration, remove it from the top of the list and
// update minDBVersion accordingly
func prepareMigrationTasks() []*migration {
	if preparedMigrations != nil {
		return preparedMigrations
	}
	preparedMigrations = []*migration{
		// Gitmin 1.0.0 starts at database version 0

		newMigration(0, "Initial but do nothing", v1_0.InitialDoNothing),
		newMigration(1, "Add TimeEstimate to Issue table", v1_0.AddTimeEstimateColumnToIssueTable),
	}
	return preparedMigrations
}

// GetCurrentDBVersion returns the current db version
func GetCurrentDBVersion(x *xorm.Engine) (int64, error) {
	if err := x.Sync(new(Version)); err != nil {
		return -1, fmt.Errorf("sync: %w", err)
	}

	currentVersion := &Version{ID: 1}
	has, err := x.Get(currentVersion)
	if err != nil {
		return -1, fmt.Errorf("get: %w", err)
	}
	if !has {
		return -1, nil
	}
	return currentVersion.Version, nil
}

func calcDBVersion(migrations []*migration) int64 {
	dbVer := int64(minDBVersion + len(migrations))
	if migrations[0].idNumber != minDBVersion {
		panic("migrations should start at minDBVersion")
	}
	if dbVer != migrations[len(migrations)-1].idNumber+1 {
		panic("migrations are not in order")
	}
	return dbVer
}

// ExpectedDBVersion returns the expected db version
func ExpectedDBVersion() int64 {
	return calcDBVersion(prepareMigrationTasks())
}

// EnsureUpToDate will check if the db is at the correct version
func EnsureUpToDate(x *xorm.Engine) error {
	currentDB, err := GetCurrentDBVersion(x)
	if err != nil {
		return err
	}

	if currentDB < 0 {
		return fmt.Errorf("database has not been initialized")
	}

	if minDBVersion > currentDB {
		return fmt.Errorf("DB version %d (<= %d) is too old for auto-migration. Upgrade to Gitea 1.6.4 first then upgrade to this version", currentDB, minDBVersion)
	}

	expectedDB := ExpectedDBVersion()

	if currentDB != expectedDB {
		return fmt.Errorf(`current database version %d is not equal to the expected version %d. Please run "gitea [--config /path/to/app.ini] migrate" to update the database version`, currentDB, expectedDB)
	}

	return nil
}

func getPendingMigrations(curDBVer int64, migrations []*migration) []*migration {
	return migrations[curDBVer-minDBVersion:]
}

func migrationIDNumberToDBVersion(idNumber int64) int64 {
	return idNumber + 1
}

// Migrate database to current version
func Migrate(x *xorm.Engine) error {
	migrations := prepareMigrationTasks()
	maxDBVer := calcDBVersion(migrations)

	// Set a new clean the default mapper to GonicMapper as that is the default for Gitea.
	x.SetMapper(names.GonicMapper{})
	if err := x.Sync(new(Version)); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	currentVersion := &Version{ID: 1}
	has, err := x.Get(currentVersion)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	} else if !has {
		// If the version record does not exist, it is a fresh installation, and we can skip all migrations.
		// XORM model framework will create all tables when initializing.
		currentVersion.ID = 0
		currentVersion.Version = maxDBVer
		if _, err = x.InsertOne(currentVersion); err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	curDBVer := currentVersion.Version
	// Outdated Gitea database version is not supported
	if curDBVer < minDBVersion {
		log.Fatal(`Gitea no longer supports auto-migration from your previously installed version.
Please try upgrading to a lower version first (suggested v1.6.4), then upgrade to this version.`)
		return nil
	}

	// Downgrading Gitea's database version not supported
	if maxDBVer < curDBVer {
		msg := fmt.Sprintf("Your database (migration version: %d) is for a newer Gitea, you can not use the newer database for this old Gitea release (%d).", curDBVer, maxDBVer)
		msg += "\nGitea will exit to keep your database safe and unchanged. Please use the correct Gitea release, do not change the migration version manually (incorrect manual operation may lose data)."
		if !setting.IsProd {
			msg += fmt.Sprintf("\nIf you are in development and really know what you're doing, you can force changing the migration version by executing: UPDATE version SET version=%d WHERE id=1;", maxDBVer)
		}
		log.Fatal("Migration Error: %s", msg)
		return nil
	}

	// Some migration tasks depend on the git command
	if git.DefaultContext == nil {
		if err = git.InitSimple(context.Background()); err != nil {
			return err
		}
	}

	// Migrate
	for _, m := range getPendingMigrations(curDBVer, migrations) {
		log.Info("Migration[%d]: %s", m.idNumber, m.description)
		// Reset the mapper between each migration - migrations are not supposed to depend on each other
		x.SetMapper(names.GonicMapper{})
		if err = m.Migrate(x); err != nil {
			return fmt.Errorf("migration[%d]: %s failed: %w", m.idNumber, m.description, err)
		}
		currentVersion.Version = migrationIDNumberToDBVersion(m.idNumber)
		if _, err = x.ID(1).Update(currentVersion); err != nil {
			return err
		}
	}
	return nil
}

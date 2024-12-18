// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package routers

import (
	"context"
	"net/http"
	"reflect"
	"runtime"

	"gitmin.com/gitmin/app/models"
	authmodel "gitmin.com/gitmin/app/models/auth"
	"gitmin.com/gitmin/app/modules/cache"
	"gitmin.com/gitmin/app/modules/eventsource"
	"gitmin.com/gitmin/app/modules/git"
	"gitmin.com/gitmin/app/modules/highlight"
	"gitmin.com/gitmin/app/modules/log"
	"gitmin.com/gitmin/app/modules/markup"
	"gitmin.com/gitmin/app/modules/markup/external"
	"gitmin.com/gitmin/app/modules/setting"
	"gitmin.com/gitmin/app/modules/ssh"
	"gitmin.com/gitmin/app/modules/storage"
	"gitmin.com/gitmin/app/modules/svg"
	"gitmin.com/gitmin/app/modules/system"
	"gitmin.com/gitmin/app/modules/templates"
	"gitmin.com/gitmin/app/modules/translation"
	"gitmin.com/gitmin/app/modules/util"
	"gitmin.com/gitmin/app/modules/web"
	"gitmin.com/gitmin/app/modules/web/routing"
	actions_router "gitmin.com/gitmin/app/routers/api/actions"
	packages_router "gitmin.com/gitmin/app/routers/api/packages"
	apiv1 "gitmin.com/gitmin/app/routers/api/v1"
	"gitmin.com/gitmin/app/routers/common"
	"gitmin.com/gitmin/app/routers/private"
	web_routers "gitmin.com/gitmin/app/routers/web"
	actions_service "gitmin.com/gitmin/app/services/actions"
	asymkey_service "gitmin.com/gitmin/app/services/asymkey"
	"gitmin.com/gitmin/app/services/auth"
	"gitmin.com/gitmin/app/services/auth/source/oauth2"
	"gitmin.com/gitmin/app/services/automerge"
	"gitmin.com/gitmin/app/services/cron"
	feed_service "gitmin.com/gitmin/app/services/feed"
	indexer_service "gitmin.com/gitmin/app/services/indexer"
	"gitmin.com/gitmin/app/services/mailer"
	mailer_incoming "gitmin.com/gitmin/app/services/mailer/incoming"
	markup_service "gitmin.com/gitmin/app/services/markup"
	repo_migrations "gitmin.com/gitmin/app/services/migrations"
	mirror_service "gitmin.com/gitmin/app/services/mirror"
	"gitmin.com/gitmin/app/services/oauth2_provider"
	pull_service "gitmin.com/gitmin/app/services/pull"
	release_service "gitmin.com/gitmin/app/services/release"
	repo_service "gitmin.com/gitmin/app/services/repository"
	"gitmin.com/gitmin/app/services/repository/archiver"
	"gitmin.com/gitmin/app/services/task"
	"gitmin.com/gitmin/app/services/uinotification"
	"gitmin.com/gitmin/app/services/webhook"
)

func mustInit(fn func() error) {
	err := fn()
	if err != nil {
		ptr := reflect.ValueOf(fn).Pointer()
		fi := runtime.FuncForPC(ptr)
		log.Fatal("%s failed: %v", fi.Name(), err)
	}
}

func mustInitCtx(ctx context.Context, fn func(ctx context.Context) error) {
	err := fn(ctx)
	if err != nil {
		ptr := reflect.ValueOf(fn).Pointer()
		fi := runtime.FuncForPC(ptr)
		log.Fatal("%s(ctx) failed: %v", fi.Name(), err)
	}
}

func syncAppConfForGit(ctx context.Context) error {
	runtimeState := new(system.RuntimeState)
	if err := system.AppState.Get(ctx, runtimeState); err != nil {
		return err
	}

	updated := false
	if runtimeState.LastAppPath != setting.AppPath {
		log.Info("AppPath changed from '%s' to '%s'", runtimeState.LastAppPath, setting.AppPath)
		runtimeState.LastAppPath = setting.AppPath
		updated = true
	}
	if runtimeState.LastCustomConf != setting.CustomConf {
		log.Info("CustomConf changed from '%s' to '%s'", runtimeState.LastCustomConf, setting.CustomConf)
		runtimeState.LastCustomConf = setting.CustomConf
		updated = true
	}

	if updated {
		log.Info("re-sync repository hooks ...")
		mustInitCtx(ctx, repo_service.SyncRepositoryHooks)

		log.Info("re-write ssh public keys ...")
		mustInitCtx(ctx, asymkey_service.RewriteAllPublicKeys)

		return system.AppState.Set(ctx, runtimeState)
	}
	return nil
}

func InitWebInstallPage(ctx context.Context) {
	translation.InitLocales(ctx)
	setting.LoadSettingsForInstall()
	mustInit(svg.Init)
}

// InitWebInstalled is for global installed configuration.
func InitWebInstalled(ctx context.Context) {
	mustInitCtx(ctx, git.InitFull)
	log.Info("Git version: %s (home: %s)", git.DefaultFeatures().VersionInfo(), git.HomeDir())
	if !git.DefaultFeatures().SupportHashSha256 {
		log.Warn("sha256 hash support is disabled - requires Git >= 2.42." + util.Iif(git.DefaultFeatures().UsingGogit, " Gogit is currently unsupported.", ""))
	}

	// Setup i18n
	translation.InitLocales(ctx)

	setting.LoadSettings()
	mustInit(storage.Init)

	mailer.NewContext(ctx)
	mustInit(cache.Init)
	mustInit(feed_service.Init)
	mustInit(uinotification.Init)
	mustInitCtx(ctx, archiver.Init)

	highlight.NewContext()
	external.RegisterRenderers()
	markup.Init(markup_service.ProcessorHelper())

	if setting.EnableSQLite3 {
		log.Info("SQLite3 support is enabled")
	} else if setting.Database.Type.IsSQLite3() {
		log.Fatal("SQLite3 support is disabled, but it is used for database setting. Please get or build a Gitea release with SQLite3 support.")
	}

	mustInitCtx(ctx, common.InitDBEngine)
	log.Info("ORM engine initialization successful!")
	mustInit(system.Init)
	mustInitCtx(ctx, oauth2.Init)
	mustInitCtx(ctx, oauth2_provider.Init)
	mustInit(release_service.Init)

	mustInitCtx(ctx, models.Init)
	mustInitCtx(ctx, authmodel.Init)
	mustInitCtx(ctx, repo_service.Init)

	// Booting long running goroutines.
	mustInit(indexer_service.Init)

	mirror_service.InitSyncMirrors()
	mustInit(webhook.Init)
	mustInit(pull_service.Init)
	mustInit(automerge.Init)
	mustInit(task.Init)
	mustInit(repo_migrations.Init)
	eventsource.GetManager().Init()
	mustInitCtx(ctx, mailer_incoming.Init)

	mustInitCtx(ctx, syncAppConfForGit)

	mustInit(ssh.Init)

	auth.Init()
	mustInit(svg.Init)

	actions_service.Init()

	mustInit(repo_service.InitLicenseClassifier)

	// Finally start up the cron
	cron.NewContext(ctx)
}

// NormalRoutes represents non install routes
func NormalRoutes() *web.Router {
	_ = templates.HTMLRenderer()
	r := web.NewRouter()
	r.Use(common.ProtocolMiddlewares()...)

	r.Mount("/", web_routers.Routes())
	r.Mount("/api/v1", apiv1.Routes())
	r.Mount("/api/internal", private.Routes())

	r.Post("/-/fetch-redirect", common.FetchRedirectDelegate)

	if setting.Packages.Enabled {
		// This implements package support for most package managers
		r.Mount("/api/packages", packages_router.CommonRoutes())
		// This implements the OCI API, this container registry "/v2" endpoint must be in the root of the site.
		// If site admin deploys Gitea in a sub-path, they must configure their reverse proxy to map the "https://host/v2" endpoint to Gitea.
		r.Mount("/v2", packages_router.ContainerRoutes())
	}

	if setting.Actions.Enabled {
		prefix := "/api/actions"
		r.Mount(prefix, actions_router.Routes(prefix))

		// TODO: Pipeline api used for runner internal communication with gitea server. but only artifact is used for now.
		// In Github, it uses ACTIONS_RUNTIME_URL=https://pipelines.actions.githubusercontent.com/fLgcSHkPGySXeIFrg8W8OBSfeg3b5Fls1A1CwX566g8PayEGlg/
		// TODO: this prefix should be generated with a token string with runner ?
		prefix = "/api/actions_pipeline"
		r.Mount(prefix, actions_router.ArtifactsRoutes(prefix))
		prefix = actions_router.ArtifactV4RouteBase
		r.Mount(prefix, actions_router.ArtifactsV4Routes(prefix))
	}

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		routing.UpdateFuncInfo(req.Context(), routing.GetFuncInfo(http.NotFound, "GlobalNotFound"))
		http.NotFound(w, req)
	})
	return r
}

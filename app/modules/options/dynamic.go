// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !bindata

package options

import (
	"code.gitea.io/gitea/app/modules/assetfs"
	"code.gitea.io/gitea/app/modules/setting"
)

func BuiltinAssets() *assetfs.Layer {
	return assetfs.Local("builtin(static)", setting.StaticRootPath, "bundles", "options")
}

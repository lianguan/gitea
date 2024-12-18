// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package templates

import (
	"slices"
	"strings"

	"gitmin.com/gitmin/app/modules/assetfs"
	"gitmin.com/gitmin/app/modules/setting"
)

func AssetFS() *assetfs.LayeredFS {
	return assetfs.Layered(CustomAssets(), BuiltinAssets())
}

func CustomAssets() *assetfs.Layer {
	return assetfs.Local("custom", setting.CustomPath, "bundles/templates")
}

func ListWebTemplateAssetNames(assets *assetfs.LayeredFS) ([]string, error) {
	files, err := assets.ListAllFiles(".", true)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(files, func(file string) bool {
		return strings.HasPrefix(file, "mail/") || !strings.HasSuffix(file, ".tmpl")
	}), nil
}

func ListMailTemplateAssetNames(assets *assetfs.LayeredFS) ([]string, error) {
	files, err := assets.ListAllFiles(".", true)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(files, func(file string) bool {
		return !strings.HasPrefix(file, "mail/") || !strings.HasSuffix(file, ".tmpl")
	}), nil
}

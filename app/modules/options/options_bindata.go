// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build bindata

package options

//go:generate go run ../../../generate/generate-bindata.go ../../../bundles/options options bindata.go

// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_0 //nolint

import (
	"xorm.io/xorm"
)

func InitialDoNothing(x *xorm.Engine) error {
	return nil
}

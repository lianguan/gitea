// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package convert

import (
	secret_model "gitmin.com/gitmin/app/models/secret"
	api "gitmin.com/gitmin/app/modules/structs"
)

// ToSecret converts Secret to API format
func ToSecret(secret *secret_model.Secret) *api.Secret {
	result := &api.Secret{
		Name: secret.Name,
	}

	return result
}

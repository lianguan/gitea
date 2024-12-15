// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package convert

import (
	secret_model "code.gitea.io/gitea/app/models/secret"
	api "code.gitea.io/gitea/app/modules/structs"
)

// ToSecret converts Secret to API format
func ToSecret(secret *secret_model.Secret) *api.Secret {
	result := &api.Secret{
		Name: secret.Name,
	}

	return result
}

	auth_model "code.gitea.io/gitea/app/models/auth"
	"code.gitea.io/gitea/app/models/unittest"
	user_model "code.gitea.io/gitea/app/models/user"
	assert.EqualValues(t, "2", resp.Header().Get("X-Total"))
	assert.EqualValues(t, "3", resp.Header().Get("X-Total"))
	assert.Empty(t, apiData)
	assert.EqualValues(t, "1", resp.Header().Get("X-Total"))
	assert.EqualValues(t, "1", resp.Header().Get("X-Total"))
// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"code.gitea.io/gitea/app/models/db"
	repo_model "code.gitea.io/gitea/app/models/repo"
	"code.gitea.io/gitea/app/models/unittest"
	user_model "code.gitea.io/gitea/app/models/user"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/indexer/stats"
	"code.gitea.io/gitea/modules/queue"
	repo_service "code.gitea.io/gitea/app/services/repository"
	files_service "code.gitea.io/gitea/app/services/repository/files"
	"code.gitea.io/gitea/tests"

	"github.com/stretchr/testify/assert"
)

func TestLinguist(t *testing.T) {
	onGiteaRun(t, func(t *testing.T, _ *url.URL) {
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

		cppContent := "#include <iostream>\nint main() {\nstd::cout << \"Hello Gitea!\";\nreturn 0;\n}"
		pyContent := "print(\"Hello Gitea!\")"
		phpContent := "<?php\necho 'Hallo Welt';\n?>"
		lockContent := "# This file is automatically @generated by Poetry 1.7.1 and should not be changed by hand."
		mdContent := "markdown"

		cases := []struct {
			GitAttributesContent  string
			FilesToAdd            []*files_service.ChangeRepoFile
			ExpectedLanguageOrder []string
		}{
			// case 0
			{
				ExpectedLanguageOrder: []string{},
			},
			// case 1
			{
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
					{
						TreePath:      "python.py",
						ContentReader: strings.NewReader(pyContent),
					},
					{
						TreePath:      "php.php",
						ContentReader: strings.NewReader(phpContent),
					},
				},
				ExpectedLanguageOrder: []string{"C++", "PHP", "Python"},
			},
			// case 2
			{
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      ".cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
					{
						TreePath:      "python.py",
						ContentReader: strings.NewReader(pyContent),
					},
					{
						TreePath:      "vendor/php.php",
						ContentReader: strings.NewReader(phpContent),
					},
				},
				ExpectedLanguageOrder: []string{"Python"},
			},
			// case 3
			{
				GitAttributesContent: "*.cpp linguist-language=Go",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
				},
				ExpectedLanguageOrder: []string{"Go"},
			},
			// case 4
			{
				GitAttributesContent: "*.cpp gitlab-language=Go?parent=json",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
				},
				ExpectedLanguageOrder: []string{"Go"},
			},
			// case 5
			{
				GitAttributesContent: "*.cpp linguist-language=HTML gitlab-language=Go?parent=json",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
				},
				ExpectedLanguageOrder: []string{"HTML"},
			},
			// case 6
			{
				GitAttributesContent: "vendor/** linguist-vendored=false",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "vendor/php.php",
						ContentReader: strings.NewReader(phpContent),
					},
				},
				ExpectedLanguageOrder: []string{"PHP"},
			},
			// case 7
			{
				GitAttributesContent: "*.cpp linguist-vendored=true\n*.py linguist-vendored\nvendor/** -linguist-vendored",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
					{
						TreePath:      "python.py",
						ContentReader: strings.NewReader(pyContent),
					},
					{
						TreePath:      "vendor/php.php",
						ContentReader: strings.NewReader(phpContent),
					},
				},
				ExpectedLanguageOrder: []string{"PHP"},
			},
			// case 8
			{
				GitAttributesContent: "poetry.lock linguist-language=Go",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "poetry.lock",
						ContentReader: strings.NewReader(lockContent),
					},
				},
				ExpectedLanguageOrder: []string{"Go"},
			},
			// case 9
			{
				GitAttributesContent: "poetry.lock linguist-generated=false",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "poetry.lock",
						ContentReader: strings.NewReader(lockContent),
					},
				},
				ExpectedLanguageOrder: []string{"TOML"},
			},
			// case 10
			{
				GitAttributesContent: "*.cpp -linguist-detectable",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
				},
				ExpectedLanguageOrder: []string{},
			},
			// case 11
			{
				GitAttributesContent: "*.md linguist-detectable",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "test.md",
						ContentReader: strings.NewReader(mdContent),
					},
				},
				ExpectedLanguageOrder: []string{"Markdown"},
			},
			// case 12
			{
				GitAttributesContent: "test2.md linguist-detectable",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "cplusplus.cpp",
						ContentReader: strings.NewReader(cppContent),
					},
					{
						TreePath:      "test.md",
						ContentReader: strings.NewReader(mdContent),
					},
					{
						TreePath:      "test2.md",
						ContentReader: strings.NewReader(mdContent),
					},
				},
				ExpectedLanguageOrder: []string{"C++", "Markdown"},
			},
			// case 13
			{
				GitAttributesContent: "README.md linguist-documentation=false",
				FilesToAdd: []*files_service.ChangeRepoFile{
					{
						TreePath:      "README.md",
						ContentReader: strings.NewReader(mdContent),
					},
				},
				ExpectedLanguageOrder: []string{"Markdown"},
			},
		}

		for i, c := range cases {
			t.Run("Case-"+strconv.Itoa(i), func(t *testing.T) {
				defer tests.PrintCurrentTest(t)()
				repo, err := repo_service.CreateRepository(db.DefaultContext, user, user, repo_service.CreateRepoOptions{
					Name: "linguist-test-" + strconv.Itoa(i),
				})
				assert.NoError(t, err)

				files := []*files_service.ChangeRepoFile{
					{
						TreePath:      ".gitattributes",
						ContentReader: strings.NewReader(c.GitAttributesContent),
					},
				}
				files = append(files, c.FilesToAdd...)
				for _, f := range files {
					f.Operation = "create"
				}

				_, err = files_service.ChangeRepoFiles(git.DefaultContext, repo, user, &files_service.ChangeRepoFilesOptions{
					Files:     files,
					OldBranch: repo.DefaultBranch,
					NewBranch: repo.DefaultBranch,
				})
				assert.NoError(t, err)

				assert.NoError(t, stats.UpdateRepoIndexer(repo))
				assert.NoError(t, queue.GetManager().FlushAll(context.Background(), 10*time.Second))

				stats, err := repo_model.GetTopLanguageStats(db.DefaultContext, repo, len(c.FilesToAdd))
				assert.NoError(t, err)

				languages := make([]string, 0, len(stats))
				for _, s := range stats {
					languages = append(languages, s.Language)
				}
				assert.Equal(t, c.ExpectedLanguageOrder, languages, "case %d: unexpected language stats", i)
			})
		}
	})
}

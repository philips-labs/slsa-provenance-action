package cli_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	gh "github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
	"github.com/philips-labs/slsa-provenance-action/lib/github"
)

const (
	githubContext = `{
		"token": "***",
		"job": "generate-provenance",
		"ref": "refs/heads/temp/dump-context",
		"sha": "c4f679f131dfb7f810fd411ac9475549d1c393df",
		"repository": "philips-labs/slsa-provenance-action",
		"repository_owner": "philips-labs",
		"repositoryUrl": "git://github.com/philips-labs/slsa-provenance-action.git",
		"run_id": "1332651620",
		"run_number": "91",
		"retention_days": "90",
		"run_attempt": "1",
		"actor": "John Doe",
		"workflow": "Integration test file provenance",
		"head_ref": "",
		"base_ref": "",
		"event_name": "push",
		"event": {
		  "after": "c4f679f131dfb7f810fd411ac9475549d1c393df",
		  "base_ref": null,
		  "before": "715b4daa0f750f420635ee488ef37a2433608438",
		  "commits": [
			{
			  "author": {
				"email": "john.doe@philips.com",
				"name": "John Doe",
				"username": "john-doe"
			  },
			  "committer": {
				"email": "noreply@github.com",
				"name": "GitHub",
				"username": "web-flow"
			  },
			  "distinct": true,
			  "id": "c4f679f131dfb7f810fd411ac9475549d1c393df",
			  "message": "Update example-local.yml",
			  "timestamp": "2021-10-12T12:18:06+02:00",
			  "tree_id": "a4dda43e9a101031dc6cd14def2d6e34ef9b4d92",
			  "url": "https://github.com/philips-labs/slsa-provenance-action/commit/c4f679f131dfb7f810fd411ac9475549d1c393df"
			}
		  ],
		  "compare": "https://github.com/philips-labs/slsa-provenance-action/compare/715b4daa0f75...c4f679f131df",
		  "created": false,
		  "deleted": false,
		  "enterprise": {
			"avatar_url": "https://avatars.githubusercontent.com/b/1244?v=4",
			"created_at": "2019-11-07T05:37:39Z",
			"description": "",
			"html_url": "https://github.com/enterprises/royal-philips",
			"id": 1244,
			"name": "Royal Philips",
			"node_id": "MDEwOkVudGVycHJpc2UxMjQ0",
			"slug": "royal-philips",
			"updated_at": "2020-12-16T12:30:18Z",
			"website_url": "https://www.philips.com"
		  },
		  "forced": false,
		  "head_commit": {
			"author": {
			  "email": "john.doe@philips.com",
			  "name": "John Doe",
			  "username": "john-doe"
			},
			"committer": {
			  "email": "noreply@github.com",
			  "name": "GitHub",
			  "username": "web-flow"
			},
			"distinct": true,
			"id": "c4f679f131dfb7f810fd411ac9475549d1c393df",
			"message": "Update example-local.yml",
			"timestamp": "2021-10-12T12:18:06+02:00",
			"tree_id": "a4dda43e9a101031dc6cd14def2d6e34ef9b4d92",
			"url": "https://github.com/philips-labs/slsa-provenance-action/commit/c4f679f131dfb7f810fd411ac9475549d1c393df"
		  },
		  "organization": {
			"avatar_url": "https://avatars.githubusercontent.com/u/58286953?v=4",
			"description": "Philips Labs - Projects in development",
			"events_url": "https://api.github.com/orgs/philips-labs/events",
			"hooks_url": "https://api.github.com/orgs/philips-labs/hooks",
			"id": 58286953,
			"issues_url": "https://api.github.com/orgs/philips-labs/issues",
			"login": "philips-labs",
			"members_url": "https://api.github.com/orgs/philips-labs/members{/member}",
			"node_id": "MDEyOk9yZ2FuaXphdGlvbjU4Mjg2OTUz",
			"public_members_url": "https://api.github.com/orgs/philips-labs/public_members{/member}",
			"repos_url": "https://api.github.com/orgs/philips-labs/repos",
			"url": "https://api.github.com/orgs/philips-labs"
		  },
		  "pusher": {
			"email": "john.doe@philips.com",
			"name": "john-doe"
		  },
		  "ref": "refs/heads/temp/dump-context",
		  "repository": {
			"allow_forking": true,
			"archive_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/{archive_format}{/ref}",
			"archived": false,
			"assignees_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/assignees{/user}",
			"blobs_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/git/blobs{/sha}",
			"branches_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/branches{/branch}",
			"clone_url": "https://github.com/philips-labs/slsa-provenance-action.git",
			"collaborators_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/collaborators{/collaborator}",
			"comments_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/comments{/number}",
			"commits_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/commits{/sha}",
			"compare_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/compare/{base}...{head}",
			"contents_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/contents/{+path}",
			"contributors_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/contributors",
			"created_at": 1631537642,
			"default_branch": "main",
			"deployments_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/deployments",
			"description": "Github Action implementation of SLSA Provenance Generation of level 1",
			"disabled": false,
			"downloads_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/downloads",
			"events_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/events",
			"fork": false,
			"forks": 2,
			"forks_count": 2,
			"forks_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/forks",
			"full_name": "philips-labs/slsa-provenance-action",
			"git_commits_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/git/commits{/sha}",
			"git_refs_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/git/refs{/sha}",
			"git_tags_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/git/tags{/sha}",
			"git_url": "git://github.com/philips-labs/slsa-provenance-action.git",
			"has_downloads": true,
			"has_issues": true,
			"has_pages": false,
			"has_projects": true,
			"has_wiki": true,
			"homepage": "",
			"hooks_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/hooks",
			"html_url": "https://github.com/philips-labs/slsa-provenance-action",
			"id": 405972862,
			"is_template": false,
			"issue_comment_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/issues/comments{/number}",
			"issue_events_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/issues/events{/number}",
			"issues_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/issues{/number}",
			"keys_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/keys{/key_id}",
			"labels_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/labels{/name}",
			"language": "Go",
			"languages_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/languages",
			"license": {
			  "key": "mit",
			  "name": "MIT License",
			  "node_id": "MDc6TGljZW5zZTEz",
			  "spdx_id": "MIT",
			  "url": "https://api.github.com/licenses/mit"
			},
			"master_branch": "main",
			"merges_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/merges",
			"milestones_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/milestones{/number}",
			"mirror_url": null,
			"name": "slsa-provenance-action",
			"node_id": "MDEwOlJlcG9zaXRvcnk0MDU5NzI4NjI=",
			"notifications_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/notifications{?since,all,participating}",
			"open_issues": 11,
			"open_issues_count": 11,
			"organization": "philips-labs",
			"owner": {
			  "avatar_url": "https://avatars.githubusercontent.com/u/58286953?v=4",
			  "email": "software-program-cto@philips.com",
			  "events_url": "https://api.github.com/users/philips-labs/events{/privacy}",
			  "followers_url": "https://api.github.com/users/philips-labs/followers",
			  "following_url": "https://api.github.com/users/philips-labs/following{/other_user}",
			  "gists_url": "https://api.github.com/users/philips-labs/gists{/gist_id}",
			  "gravatar_id": "",
			  "html_url": "https://github.com/philips-labs",
			  "id": 58286953,
			  "login": "philips-labs",
			  "name": "philips-labs",
			  "node_id": "MDEyOk9yZ2FuaXphdGlvbjU4Mjg2OTUz",
			  "organizations_url": "https://api.github.com/users/philips-labs/orgs",
			  "received_events_url": "https://api.github.com/users/philips-labs/received_events",
			  "repos_url": "https://api.github.com/users/philips-labs/repos",
			  "site_admin": false,
			  "starred_url": "https://api.github.com/users/philips-labs/starred{/owner}{/repo}",
			  "subscriptions_url": "https://api.github.com/users/philips-labs/subscriptions",
			  "type": "Organization",
			  "url": "https://api.github.com/users/philips-labs"
			},
			"private": false,
			"pulls_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/pulls{/number}",
			"pushed_at": 1634033886,
			"releases_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/releases{/id}",
			"size": 76,
			"ssh_url": "git@github.com:philips-labs/slsa-provenance-action.git",
			"stargazers": 1,
			"stargazers_count": 1,
			"stargazers_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/stargazers",
			"statuses_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/statuses/{sha}",
			"subscribers_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/subscribers",
			"subscription_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/subscription",
			"svn_url": "https://github.com/philips-labs/slsa-provenance-action",
			"tags_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/tags",
			"teams_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/teams",
			"topics": [
			  "hacktoberfest"
			],
			"trees_url": "https://api.github.com/repos/philips-labs/slsa-provenance-action/git/trees{/sha}",
			"updated_at": "2021-10-11T14:10:37Z",
			"url": "https://github.com/philips-labs/slsa-provenance-action",
			"visibility": "public",
			"watchers": 1,
			"watchers_count": 1
		  },
		  "sender": {
			"avatar_url": "https://avatars.githubusercontent.com/u/15904543?v=4",
			"events_url": "https://api.github.com/users/john-doe/events{/privacy}",
			"followers_url": "https://api.github.com/users/john-doe/followers",
			"following_url": "https://api.github.com/users/john-doe/following{/other_user}",
			"gists_url": "https://api.github.com/users/john-doe/gists{/gist_id}",
			"gravatar_id": "",
			"html_url": "https://github.com/john-doe",
			"id": 15904543,
			"login": "john-doe",
			"node_id": "MDQ6VXNlcjE1OTA0NTQz",
			"organizations_url": "https://api.github.com/users/john-doe/orgs",
			"received_events_url": "https://api.github.com/users/john-doe/received_events",
			"repos_url": "https://api.github.com/users/john-doe/repos",
			"site_admin": false,
			"starred_url": "https://api.github.com/users/john-doe/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/john-doe/subscriptions",
			"type": "User",
			"url": "https://api.github.com/users/john-doe"
		  }
		},
		"server_url": "https://github.com",
		"api_url": "https://api.github.com",
		"graphql_url": "https://api.github.com/graphql",
		"ref_protected": false,
		"ref_type": "branch",
		"workspace": "/home/runner/work/slsa-provenance-action/slsa-provenance-action",
		"action": "__self",
		"event_path": "/home/runner/work/_temp/_github_workflow/event.json",
		"action_repository": "",
		"action_ref": "",
		"path": "/home/runner/work/_temp/_runner_file_commands/add_path_779d6e30-d262-4e4a-bcdf-bf652ff08e12",
		"env": "/home/runner/work/_temp/_runner_file_commands/set_env_779d6e30-d262-4e4a-bcdf-bf652ff08e12"
	  }`
	runnerContext = `{
		"os": "Linux",
		"name": "Hosted Agent",
		"tool_cache": "/opt/hostedtoolcache",
		"temp": "/home/runner/work/_temp",
		"workspace": "/home/runner/work/slsa-provenance-action"
	  }`
)

func TestGenerateCliOptions(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	provenanceFile := path.Join(rootDir, "bin/unittest.provenance")

	testCases := []struct {
		name      string
		err       error
		arguments []string
	}{
		{
			name:      "without commandline flags",
			err:       cli.RequiredFlagError("-artifact_path"),
			arguments: make([]string, 0),
		},
		{
			name: "only providing artifact_path",
			err:  cli.RequiredFlagError("-github_context"),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
			},
		},
		{
			name: "without runner_context",
			err:  cli.RequiredFlagError("-runner_context"),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
			},
		},
		{
			name: "invalid artifact_path",
			err:  fmt.Errorf("failed to generate provenance: resource path not found: [provided=non-existing-folder/unknown-file]"),
			arguments: []string{
				"-artifact_path",
				"non-existing-folder/unknown-file",
				"-github_context",
				githubContext,
				"-runner_context",
				runnerContext,
			},
		},
		{
			name: "all arguments explicit",
			err:  nil,
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

			sb := strings.Builder{}

			cli := cli.Generate(&sb)
			err := cli.ParseAndRun(context.Background(), tc.arguments)
			defer func() {
				_ = os.Remove(provenanceFile)
			}()

			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
				assert.Contains(sb.String(), "Saving provenance to")
				if assert.FileExists(provenanceFile) {
					content, err := os.ReadFile(provenanceFile)
					assert.NoError(err)
					assert.Greater(len(content), 1)
				}
			}
		})
	}
}

func TestProvenenaceGitHubRelease(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	artifactPath := path.Join(rootDir, "gh-release-test")
	provenanceFile := path.Join(artifactPath, "unittest.provenance")

	ctx := context.Background()
	owner, repo := "philips-labs", "slsa-provenance-action"
	oauthClient := github.NewOAuth2Client(ctx, func() string { return os.Getenv("GITHUB_TOKEN") })
	client := github.NewReleaseClient(oauthClient)

	releaseID, err := createGitHubRelease(
		ctx,
		client,
		owner,
		repo,
		"v0.0.0-generate-test",
		path.Join(rootDir, "bin", "slsa-provenance"),
		path.Join(rootDir, "README.md"),
	)
	assert.NoError(err)

	defer func() {
		_ = os.RemoveAll(artifactPath)
		_, err = client.Repositories.DeleteRelease(ctx, owner, repo, releaseID)
	}()

	sb := strings.Builder{}
	cli := cli.Generate(&sb)
	err = cli.ParseAndRun(
		context.Background(),
		[]string{
			"-artifact_path",
			artifactPath,
			"-github_context",
			githubContext,
			"-output_path",
			provenanceFile,
			"-runner_context",
			runnerContext,
			"-tag_name",
			"v0.0.0-generate-test",
		},
	)
	assert.NoError(err)
	assert.Contains(sb.String(), "Saving provenance to")
	if assert.FileExists(provenanceFile) {
		content, err := os.ReadFile(provenanceFile)
		assert.NoError(err)
		assert.Greater(len(content), 1)
	}
}

func createGitHubRelease(ctx context.Context, client *github.ReleaseClient, owner, repo, version string, assets ...string) (int64, error) {
	rel, _, err := client.Repositories.CreateRelease(
		ctx,
		owner,
		repo,
		&gh.RepositoryRelease{TagName: stringPointer(version), Name: stringPointer(version), Draft: boolPointer(true), Prerelease: boolPointer(true)},
	)
	if err != nil {
		return 0, err
	}

	for _, a := range assets {
		asset, err := os.Open(a)
		if err != nil {
			return 0, err
		}
		client.AddProvenanceToRelease(ctx, owner, repo, rel.GetID(), asset)
	}

	return rel.GetID(), nil
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

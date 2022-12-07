package github_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/pkg/github"
	"github.com/philips-labs/slsa-provenance-action/pkg/intoto"
)

const (
	pushGitHubEvent = `{
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
  }`
)

func TestGenerateProvenance(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	os.Setenv("GITHUB_ACTIONS", "true")

	repoURL := "https://github.com/philips-labs/slsa-provenance-action"

	gh := github.Context{
		RunID:           "1029384756",
		RepositoryOwner: "philips-labs",
		Repository:      "philips-labs/slsa-provenance-action",
		Event:           []byte(pushGitHubEvent),
		EventName:       "push",
		ActionPath:      ".github/workflows/build.yml",
		SHA:             "849fb987efc0c0fc72e26a38f63f0c00225132be",
	}
	materials := []intoto.Item{
		{URI: "git+" + repoURL, Digest: intoto.DigestSet{"sha1": gh.SHA}},
	}

	runner := github.RunnerContext{}
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")

	artifactPath := path.Join(rootDir, "bin")
	fps := intoto.NewFilePathSubjecter(artifactPath)

	env := github.Environment{
		Context: &gh,
		Runner:  &runner,
	}
	stmt, err := env.GenerateProvenanceStatement(ctx, fps)
	if !assert.NoError(err) {
		return
	}

	binaryName := "slsa-provenance"
	binaryPath := path.Join(artifactPath, binaryName)

	assert.Len(stmt.Subject, 1)
	assertSubject(assert, stmt.Subject, binaryName, binaryPath)

	assert.Equal(intoto.SlsaPredicateType, stmt.PredicateType)
	assert.Equal(intoto.StatementType, stmt.Type)

	predicate := stmt.Predicate
	assert.Equal(github.BuildType, predicate.BuildType)
	assert.Equal(fmt.Sprintf("%s%s", repoURL, github.HostedIDSuffix), predicate.ID)
	assert.Equal(materials, predicate.Materials)
	assert.Equal(fmt.Sprintf("%s%s", repoURL, github.HostedIDSuffix), predicate.Builder.ID)

	assertMetadata(assert, predicate.Metadata, gh, repoURL)
	assertInvocation(assert, predicate.Invocation)
}

func TestGenerateProvenanceFromGitHubRelease(t *testing.T) {
	if tokenRetriever() == "" {
		t.Skip("skipping as GITHUB_TOKEN environment variable isn't set")
	}
	assert := assert.New(t)

	ctx := context.Background()
	os.Setenv("GITHUB_ACTIONS", "true")

	repoURL := "https://github.com/philips-labs/slsa-provenance-action"

	ghContext := github.Context{
		RunID:           "1029384756",
		RepositoryOwner: "philips-labs",
		Repository:      "philips-labs/slsa-provenance-action",
		Event:           []byte(pushGitHubEvent),
		EventName:       "push",
		ActionPath:      ".github/workflows/build.yml",
		SHA:             "849fb987efc0c0fc72e26a38f63f0c00225132be",
	}
	materials := []intoto.Item{
		{URI: "git+" + repoURL, Digest: intoto.DigestSet{"sha1": ghContext.SHA}},
	}

	runner := github.RunnerContext{}
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	artifactPath := path.Join(rootDir, "release-assets")
	fps := intoto.NewFilePathSubjecter(artifactPath)

	tc := github.NewOAuth2Client(ctx, tokenRetriever)
	client := github.NewReleaseClient(tc)

	version := fmt.Sprintf("v0.0.0-rel-test-%d", time.Now().UnixNano())
	releaseId, err := createGitHubRelease(
		ctx,
		client,
		owner,
		repo,
		version,
		path.Join(rootDir, "bin", "slsa-provenance"),
		path.Join(rootDir, "README.md"),
	)
	if !assert.NoError(err) {
		return
	}
	defer func() {
		_ = os.RemoveAll(artifactPath)
		_, err := client.Repositories.DeleteRelease(ctx, owner, repo, releaseId)
		assert.NoError(err)
	}()

	env := github.NewReleaseEnvironment(ghContext, runner, version, client, artifactPath)
	stmt, err := env.GenerateProvenanceStatement(ctx, fps)
	if !assert.NoError(err) {
		return
	}

	binaryName := "slsa-provenance"
	binaryPath := path.Join(artifactPath, binaryName)
	readmeName := "README.md"
	readmePath := path.Join(artifactPath, readmeName)

	assert.Len(stmt.Subject, 2)
	assertSubject(assert, stmt.Subject, binaryName, binaryPath)
	assertSubject(assert, stmt.Subject, readmeName, readmePath)

	assert.Equal(intoto.SlsaPredicateType, stmt.PredicateType)
	assert.Equal(intoto.StatementType, stmt.Type)

	predicate := stmt.Predicate
	assert.Equal(fmt.Sprintf("%s%s", repoURL, github.HostedIDSuffix), predicate.ID)
	assert.Equal(materials, predicate.Materials)
	assert.Equal(fmt.Sprintf("%s%s", repoURL, github.HostedIDSuffix), predicate.Builder.ID)
	assert.Equal(github.BuildType, predicate.BuildType)

	assertMetadata(assert, predicate.Metadata, ghContext, repoURL)
	assertInvocation(assert, predicate.Invocation)

	stmtPath := path.Join(artifactPath, "provenance.json")

	err = env.PersistProvenanceStatement(ctx, stmt, stmtPath)
	assert.NoError(err)
}

func TestGenerateProvenanceFromGitHubReleaseErrors(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	os.Setenv("GITHUB_ACTIONS", "true")

	ghContext := github.Context{
		RunID:           "1029384756",
		RepositoryOwner: "philips-labs",
		Repository:      "philips-labs/slsa-provenance-action",
		Event:           []byte(pushGitHubEvent),
		EventName:       "push",
		SHA:             "849fb987efc0c0fc72e26a38f63f0c00225132be",
	}

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	client := github.NewReleaseClient(nil)

	version := fmt.Sprintf("v0.0.0-rel-test-%d", time.Now().UnixNano())

	env := github.NewReleaseEnvironment(ghContext, github.RunnerContext{}, version, client, rootDir)

	fps := intoto.NewFilePathSubjecter(rootDir)
	stmt, err := env.GenerateProvenanceStatement(ctx, fps)
	assert.EqualError(err, "artifactPath has to be an empty directory")
	assert.Nil(stmt)

	fps = intoto.NewFilePathSubjecter(path.Join(rootDir, "README.md"))
	env = github.NewReleaseEnvironment(ghContext, github.RunnerContext{}, version, client, path.Join(rootDir, "README.md"))
	stmt, err = env.GenerateProvenanceStatement(ctx, fps)
	assert.EqualError(err, fmt.Sprintf("mkdir %s: not a directory", path.Join(rootDir, "README.md")))
	assert.Nil(stmt)
}

func assertInvocation(assert *assert.Assertions, recipe intoto.Invocation) {
	assert.Equal(".github/workflows/build.yml", recipe.ConfigSource.EntryPoint)
	assert.Nil(recipe.Environment)
	assert.Nil(recipe.Parameters)
}

func assertMetadata(assert *assert.Assertions, meta intoto.Metadata, gh github.Context, repoURL string) {
	bft, err := time.Parse(time.RFC3339, meta.BuildFinishedOn)
	assert.NoError(err)
	assert.WithinDuration(time.Now().UTC(), bft, 1200*time.Millisecond)
	assert.Equal(fmt.Sprintf("%s/%s/%s", repoURL, "actions/runs", gh.RunID), meta.BuildInvocationID)
	assert.Equal(true, meta.Completeness.Parameters)
	assert.Equal(false, meta.Completeness.Environment)
	assert.Equal(false, meta.Completeness.Materials)
	assert.Equal(false, meta.Reproducible)
}

func assertSubject(assert *assert.Assertions, subject []intoto.Subject, binaryName, binaryPath string) {
	binary, err := os.ReadFile(binaryPath)
	if !assert.NoError(err) {
		return
	}

	shaHex := intoto.ShaSum256HexEncoded(binary)
	assert.Contains(subject, intoto.Subject{Name: binaryName, Digest: intoto.DigestSet{"sha256": shaHex}})
}

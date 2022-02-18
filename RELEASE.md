# Release procedures

## Automated release procedure

To make a new release you can make use of the following `make` task.

```bash
make gh-release NEW_VERSION=v0.7.2 OLD_VERSION=v0.5.0 DESCRIPTION="A test release to see how it works"
```

`NEW_VERSION` the version that you want to release.
`OLD_VERSION` the current version you wish to replace in the markdown and yaml files.
`DESCRIPTION` the annotation used when tagging the release.

Visit <https://github.com/philips-labs/slsa-provenance-action/releases>.
Edit the release and save it to publish to GitHub Marketplace.

> :warning: **NOTE:** when you need to test some changes in `.goreleaser.yml`, also apply the changes to `.goreleaser.draft.yml`. Then make sure your new `tag` ends with `-draft` (e.g.: `v0.7.2-draft`) to make a draft release to not notify our consumers when testing updates to the release process.

### ⚠ Important alert for MacOS users ⚠

On MacOS `sed` has different behaviour and therefore doesn't work out of the box.
A workaround to make it work is to install gnu-sed and alias it in your bashrc/zshrc:

```bash
brew install gnu-sed
echo "alias sed=gsed" >> ~/.zshrc
```

## Manual release procedure

1. Upgrade version number in all repository files, find & replace previous version number with new version number.
1. Commit the changed files.
1. Tag the new commit using `git tag -sam "What is this release about?" v0.1.0`.
1. Push the tag to remote using `git push origin v0.1.0`
1. Wait for the release workflow to finish, then push the main branch using `git push`
1. Visit <https://github.com/philips-labs/slsa-provenance-action/releases>.
1. Edit the release and save it to publish to GitHub Marketplace.

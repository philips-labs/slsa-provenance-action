# Automated release procedure

A make script has been created to automate the manual procedure.
Execute the following command:
```bash
make gh-release NEW_VERSION=v0.6.0 OLD_VERSION=v0.5.0 DESCRIPTION="A test release to see how it works"
```

`NEW_VERSION` is the version that you want to release.
`OLD_VERSION` is the previous version you wish to overwrite in the markdown and yaml files.
`DESCRIPTION` is the description to use in the annotation of the tag and commit description.

# Manual release procedure

1. Upgrade version number in all repository files, find & replace previous version number with new version number.
1. Commit the changed files.
1. Tag the new commit using `git tag -sam "What is this release about?" v0.1.0`.
1. Push the tag to remote using `git push v0.1.0`
1. Wait for the release workflow to finish, then push the main branch using `git push`


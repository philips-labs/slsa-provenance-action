# Release

1. Upgrade version number in all repository files, find & replace previous version number with new version number.
1. Commit the changed files.
1. Tag the new commit using `git tag -sam "What is this release about?" v0.1.0`.

## Experimental

1. Push the tag to remote using `git push v0.1.0`
1. Wait for the release workflow to finish, then push the main branch using `git push`

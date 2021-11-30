# Contributing slsa-provenance-action

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to SLSA-provenance-action, which is hosted at <https://github.com/philips-labs/slsa-provenance-action>. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Styleguides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line
- Sign-off your commits with ``git commit -s -m "Normal Commit Message here"``, this will add ``Signed-off-by: Random J Developer <random@developer.example.org>`` at the end of the commit.

### Code

This repository has a `.editorconfig`, please ensure to follow this `.editorconfig` styles to prevent unnecessary `diffs` on the codebase. For your convenience you can choose to install a `editorconfig` plugin in the IDEA of your choice.

In summary:

- Go files
  - indented with Tabs
  - Tabwidth of 4
- The Makefile
  - indented with Tabs
  - Tabwidth of 4
- Other files
  - indented with spaces
  - tabwidth of 2
- Files end with a newline
- Whitespace is trimmed from the end of a line

### Update README.md

Update README.md when input parameters change.

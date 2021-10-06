# Contributing to gothermostat

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to SLSA-provenance-action, which is hosted at <https://github.com/philips-labs/slsa-provenance-action>. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Styleguides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line
- Consider starting the commit message with an applicable emoji:
  - :art: when improving the format/structure of the code
  - :racehorse: when improving performance
  - :non-potable_water: when plugging memory leaks
  - :memo: when writing docs
  - :penguin: when fixing something on Linux
  - :apple: when fixing something on macOS
  - :checkered_flag: when fixing something on Windows
  - :bug: when fixing a bug
  - :fire: when removing code or files
  - :green_heart: when fixing the CI build
  - :white_check_mark: when adding tests
  - :lock: when dealing with security
  - :arrow_up: when upgrading dependencies
  - :arrow_down: when downgrading dependencies
  - :shirt: when removing linter warnings

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

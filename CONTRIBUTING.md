# Contributing to enval

- [Contributing to enval](#contributing-to-enval)
  * [Contributing to the code](#contributing-to-the-code)
    + [Set up your development environment](#set-up-your-development-environment)
    + [Building](#building)
    + [Running tests](#running-tests)
  * [Contributing with more tool specs](#contributing-with-more-tool-specs)

You can contribute to `enval` either by proposing changes in the code (fixes, new features, improvements, etc) or by adding support for more tools to the existing ones at [tool-specs](./tool-specs).

## Contributing to the code

To contribute to the `enval` code you just need to:

1. [Create an issue](https://github.com/Adhara-Tech/enval/issues), please search first is there already an issue for your same problem.
2. Clone the repo and start coding!
3. Open Pull Request (please, update the `[Unreleased]` section of the [CHANGELOG](./CHANGELOG.md) if necessary)

### Set up your development environment

Refer to [./.enval.manifest.yaml](./.enval.manifest.yaml) to find the tools you need to build.

### Building

Run `make build` will compile and build the binary executables and output them in the `bin/` folder

### Running tests

Run `make test` to run the tests.

## Contributing with more tool specs

If you just want to add support for more tools as defaults in enval, you just need to open a PR adding a `<tool_name>.tool.yaml` file inside [tool-specs/](./tool-specs) folder, following the steps described at [How to define a new tool or override an existing one](https://github.com/Adhara-Tech/enval#how-to-define-a-new-tool-or-override-an-existing-one).

A test is mandatory for the new supported tool. For that you just need to add a sample of the output of the version command of the new tool to the [./pkg/manifestchecker/testdata/](./pkg/manifestchecker/testdata) folder and then add the correspondant entry in the `TestToolsManager_Validate` table test in [./pkg/manifestchecker/tool_checker_manager_test.go](./pkg/manifestchecker/tool_checker_manager_test.go).

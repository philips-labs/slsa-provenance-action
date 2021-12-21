<div id="top"></div>

<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

</div>

<br />
<div align="center">
  <a href="https://github.com/philips-labs/slsa-provenance-action">
    <img src="https://slsa.dev/images/levelBadge1.svg" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">SLSA Provenance GitHub Action</h3>

  <p align="center">
    Github Action implementation of SLSA Provenance Generation of level 1
    <br>
    <a href="https://github.com/philips-labs/slsa-provenance-action/issues">Report Bug</a>
    ·
    <a href="https://github.com/philips-labs/slsa-provenance-action/issues">Request Feature</a>
  </p>
</div>

<!-- ABOUT THE PROJECT -->
## About This Project

This GitHub action implements the level 1 requirements of the [SLSA framework](https://slsa.dev/). By using this GitHub Action it is possible to easily generate the provenance file for different artifact types.
Different artifact types include, but not limited to:

- Files
- Push event (Docker Hub, trigger different workflow, etc)

While there are no integrity guarantees on the produced provenance at L1,
publishing artifact provenance in a common format opens up opportunities for
automated analysis and auditing. Additionally, moving build definitions into
source control and onto well-supported, secure build systems represents a marked
improvement from the ecosystem's current state.

This is not an official GitHub Action set up and maintained by the SLSA team. This GitHub Action is built for research purposes by Philips Research. It is heavily inspired by the original [Provenance Action example](https://github.com/slsa-framework/github-actions-demo) built by SLSA.

<p align="right">(<a href="#top">back to top</a>)</p>

## Background

[SLSA](https://github.com/slsa-framework/slsa) is a framework intended to codify
and promote secure software supply-chain practices. SLSA helps trace software
artifacts (e.g. binaries) back to the build and source control systems that
produced them using in-toto's
[Attestation](https://github.com/in-toto/attestation/blob/main/spec/README.md)
metadata format.

### Built With

- [SLSA Framework](https://github.com/slsa-framework/slsa/)
- [Golang](https://golang.org/)
- [GitHub Actions](https://github.com/features/actions)

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

Get started quickly by reading the information below.

### Prerequisites

Ensure you have the following installed:

- Golang
- Docker

#### Recommendations

The following IDE is recommended when working on this codebase:

- [VSCode](https://code.visualstudio.com/)

### Local Installation

1. Clone the repo.

   ```sh
   git clone git@github.com:philips-labs/slsa-provenance-action.git
   ```

1. Build the binary.

   ```sh
   make build
   ```

1. Execute the binary.

   ```sh
   ./bin/slsa-provenance help
   ```

### Docker Image

Our Docker images are available at both GitHub Container Registry (ghcr) and Docker Hub.

**Docker Hub**
See all available images [here.](https://hub.docker.com/r/philipssoftware/slsa-provenance/tags)
Run the Docker image by doing:

```sh
docker run philipssoftware/slsa-provenance:v0.4.0
```

**GitHub Container Registry**
See all available images [here.](https://github.com/philips-labs/slsa-provenance-action/pkgs/container/slsa-provenance)
Run the Docker image by doing:

```sh
docker run ghcr.io/philips-labs/slsa-provenance:v0.4.0
```

The Docker image includes the working binary that can be executed by using the ``slsa-provenance`` command.

<p align="right">(<a href="#top">back to top</a>)</p>

## Usage

The easiest way to use this action is to add the following into your workflow file. Additional configuration might be necessary to fit your usecase.

<details>
  <summary>GitHub Releases</summary>

  Add the following part in your workflow file:

  See [ci workflow](.github/workflows/ci.yaml) for a full example using GitHub releases.

  > :warning: **NOTE:** this job depends on a release job that publishes the release assets to a GitHub Release.

  ```yaml
  provenance:
    name: provenance
    needs: [release]
    runs-on: ubuntu-20.04

    steps:
      - name: Generate provenance for Release
        uses: philips-labs/slsa-provenance-action@v0.4.0
        with:
          command: generate
          subcommand: files
          arguments: --artifact-path release-assets --output-path 'provenance.json' --tag-name: ${{ github.ref_name }}
        env:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
  ```

</details>

<details>
  <summary>GitHub artifacts</summary>

  Add the following part in your workflow file:

  See [example workflow](.github/workflows/example-publish.yaml) for a full example using GitHub artifacts.

  ```yaml
  generate-provenance:
    name: Generate build provenance
    runs-on: ubuntu-latest
    steps:
      - name: Download build artifact
        uses: actions/download-artifact@v2
        with:
          path: artifact/

      - name: Download extra materials for provenance
        uses: actions/download-artifact@v2
        with:
          name: extra-materials
          path: extra-materials/

      - name: Generate provenance
        uses: philips-labs/SLSA-Provenance-Action@8c78a6b34703824b9561a26b1ae5893beea9a332
        with:
          artifact_path: artifact/
          extra_materials: extra-materials/file1.json extra-materials/some-more.json

      - name: Upload provenance
        uses: actions/upload-artifact@v2
        with:
          path: provenance.json
  ```

</details>

### Description

An action to generate SLSA build provenance for an artifact

### Inputs

| parameter | description | required | default |
| - | - | - | - |
| subcommand | The subcommand to use when generating provenance | `false` | files |
| artifact_path | path to artifact or directory of artifacts | `true` |  |
| output_path | path to write build provenance file | `false` | provenance.json |
| github_context | internal (do not set): the "github" context object in json | `true` | ${{ toJSON(github) }} |
| runner_context | internal (do not set): the "runner" context object in json | `true` | ${{ toJSON(runner) }} |
| tag_name | The github release to generate provenance on.\n (if set the artifacts will be downloaded from the release and the provenance will be added as an additional release asset.) | `false` |  |
| extra_materials | paths to files containing SLSA v0.1 formatted materials (JSON array) in to include in the provenance | `false` |  |

### Available commands

- Generate
  Command that generates the provenance.
- Help
  Command that prints helpful information about what commands and arguments can be used.
- Version
  Command that prints version information about the binary that is used.

<p align="right">(<a href="#top">back to top</a>)</p>

## Contributing

If you have a suggestion that would make this project better, please fork the repository and create a pull request. You can also simply open an issue with the tag "enhancement".

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please refer to the [Contributing Guidelines](/CONTRIBUTING.md) for all the guidelines.

<p align="right">(<a href="#top">back to top</a>)</p>

## License

Distributed under the MIT License. See [LICENSE](/LICENSE.md) for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

## Contact

*Powered by Philips SWAT Eindhoven*

- [Brend Smits](https://github.com/Brend-Smits) - brend.smits@philips.com
- [Marco Franssen](https://github.com/marcofranssen)
- [Jeroen Knoops](https://github.com/JeroenKnoops)
- [Annie Jovitha](https://github.com/AnnieJovitha)

<p align="right">(<a href="#top">back to top</a>)</p>

## Acknowledgments

This project is inspired by:

- [SLSA Framework](https://slsa.dev/)
- [SLSA GitHub Action Example](https://raw.githubusercontent.com/slsa-framework/github-actions-demo)

<p align="right">(<a href="#top">back to top</a>)</p>

[contributors-shield]: https://img.shields.io/github/contributors/philips-labs/slsa-provenance-action.svg?style=for-the-badge
[contributors-url]: https://github.com/philips-labs/slsa-provenance-action/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/philips-labs/slsa-provenance-action.svg?style=for-the-badge
[forks-url]: https://github.com/philips-labs/slsa-provenance-action/network/members
[stars-shield]: https://img.shields.io/github/stars/philips-labs/slsa-provenance-action.svg?style=for-the-badge
[stars-url]: https://github.com/philips-labs/slsa-provenance-action/stargazers
[issues-shield]: https://img.shields.io/github/issues/philips-labs/slsa-provenance-action.svg?style=for-the-badge
[issues-url]: https://github.com/philips-labs/slsa-provenance-action/issues
[license-shield]: https://img.shields.io/github/license/philips-labs/slsa-provenance-action.svg?style=for-the-badge
[license-url]: https://github.com/philips-labs/slsa-provenance-action/blob/main/LICENSE.md

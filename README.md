<p align="center">
  <a href="https://pkg.go.dev/github.com/crazy-max/gomod-updates/pkg/gomodupdates"><img src="https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white&style=flat-square" alt="PkgGoDev"></a>
  <a href="https://github.com/crazy-max/gomod-updates/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/gomod-updates.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/gomod-updates/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/gomod-updates/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/gomod-updates/actions?workflow=build"><img src="https://img.shields.io/github/actions/workflow/status/crazy-max/gomod-updates/build.yml?label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/gomod-updates/"><img src="https://img.shields.io/docker/stars/crazymax/gomod-updates.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/gomod-updates/"><img src="https://img.shields.io/docker/pulls/crazymax/gomod-updates.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/gomod-updates"><img src="https://goreportcard.com/badge/github.com/crazy-max/gomod-updates?style=flat-square" alt="Go Report"></a>
  <a href="https://codecov.io/gh/crazy-max/gomod-updates"><img src="https://img.shields.io/codecov/c/github/crazy-max/gomod-updates?logo=codecov&style=flat-square" alt="Codecov"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

Report available Go module updates, including major-version candidates.

## Installation

```console
$ go install github.com/crazy-max/gomod-updates/cmd/gomod-updates@latest
```

## Usage

Run in a Go module:

```console
$ gomod-updates --update --direct
+---------------------------+---------+-------------+--------+------------------+
| Module                    | Version | New Version | Direct | Valid Timestamps |
+---------------------------+---------+-------------+--------+------------------+
| github.com/example/module | v1.0.0  | v1.1.0      | true   | true             |
+---------------------------+---------+-------------+--------+------------------+
```

Check major-version module path candidates:

```console
$ gomod-updates --update --direct --major
+---------------------------+---------+-------------------------------------+--------+------------------+
| Module                    | Version | New Version                         | Direct | Valid Timestamps |
+---------------------------+---------+-------------------------------------+--------+------------------+
| github.com/example/module | v1.0.0  | github.com/example/module/v2@v2.0.0 | true   | true             |
+---------------------------+---------+-------------------------------------+--------+------------------+
```

You can also pipe `go list` output, like `go-mod-outdated`:

```console
$ go list -mod=mod -u -m -json all | gomod-updates --update --direct --major
```

To output a Markdown table:

```console
$ gomod-updates --update --direct --major --format markdown
```

To fail when updates are found:

```console
$ gomod-updates --update --direct --major --ci
```

### Flags

```console
  -h, --help
        Show context-sensitive help.
  --update
        List only modules with updates.
  --direct
        List only direct modules.
  --major
        Check for major-version module path candidates.
  --ci
        Non-zero exit code when at least one outdated dependency was found.
  --format string
        Output format (default,markdown). (default "default")
  --mod string
        Module download mode for go list calls. (default "mod")
  --version
        Print version information.
```

## License

MIT. See `LICENSE` for more details.

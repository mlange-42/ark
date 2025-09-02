+++
title = 'Ark'
type = 'docs'
description = 'Ark is an archetype-based Entity Component System (ECS) for Go.'
[params]
next = "/quickstart"
+++
{{< html >}}
<img src="images/logo-light.svg" alt="Ark" class="light" style="width: 100%; max-width: 680px; margin:24px auto 36px auto;"/>
<img src="images/logo-dark.svg" alt="Ark" class="dark" style="width: 100%; max-width: 680px; margin:24px auto 36px auto;"/>

<div style="width 100%; text-align: center;">
<a href="https://github.com/mlange-42/ark/actions/workflows/tests.yml" style="display:inline-block">
<img alt="Test status" src="https://img.shields.io/github/actions/workflow/status/mlange-42/ark/tests.yml?branch=main&label=Tests&logo=github" style="margin:0;"></img></a>

<a href="https://badge.coveralls.io/github/mlange-42/ark?branch=main" style="display:inline-block">
<img alt="Coverage Status" src="https://img.shields.io/coverallsCoverage/github/mlange-42/ark?logo=coveralls" style="margin:0;"></img></a>

<a href="https://goreportcard.com/report/github.com/mlange-42/ark" style="display:inline-block">
<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mlange-42/ark" style="margin:0;"></img></a>

<a href="https://mlange-42.github.io/ark/" style="display:inline-block">
<img alt="User Guide" src="https://img.shields.io/badge/user_guide-%23007D9C?logo=go&logoColor=white&labelColor=gray" style="margin:0;"></img></a>

<a href="https://pkg.go.dev/github.com/mlange-42/ark" style="display:inline-block">
<img alt="Go Reference" src="https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray" style="margin:0;"></img></a>

<a href="https://github.com/mlange-42/ark" style="display:inline-block">
<img alt="GitHub" src="https://img.shields.io/badge/github-repo-blue?logo=github" style="margin:0;"></img></a>

<a href="https://doi.org/10.5281/zenodo.14994239" style="display:inline-block">
<img alt="Zenodo DOI: 10.5281/zenodo.14994239" src="https://img.shields.io/badge/10.5281%2Fzenodo.14994239-blue?label=doi" style="margin:0;"></img></a>

<a href="https://github.com/mlange-42/ark/blob/main/LICENSE-MIT" style="display:inline-block">
<img alt="MIT license" src="https://img.shields.io/badge/MIT-brightgreen?label=license" style="margin:0;"></img></a>

<a href="https://github.com/mlange-42/ark/blob/main/LICENSE-APACHE" style="display:inline-block">
<img alt="Apache license" src="https://img.shields.io/badge/Apache%202.0-brightgreen?label=license" style="margin:0;"></img></a>

<a href="https://github.com/avelino/awesome-go" style="display:inline-block">
<img alt="Mentioned in Awesome Go" src="https://awesome.re/mentioned-badge.svg" style="margin:0;"></img></a>
</div>
{{< /html >}}

Ark is an archetype-based [Entity Component System](https://en.wikipedia.org/wiki/Entity_component_system) for [Go](https://go.dev/). It is the successor of [Arche](https://github.com/mlange-42/arche).

## Ark's Features

- Designed for performance and highly optimized. See the [Benchmarks](https://mlange-42.github.io/ark/benchmarks/).
- Well-documented, type-safe [API](https://pkg.go.dev/github.com/mlange-42/ark), and a comprehensive [User guide](https://mlange-42.github.io/ark/).
- [Entity relationships](https://mlange-42.github.io/ark/relations/) as a first-class feature.
- Fast [batch operations](https://mlange-42.github.io/ark/batch/) for mass manipulation.
- No systems. Just queries. Use your own structure (or the [Tools](https://github.com/mlange-42/ark#tools)).
- World serialization and deserialization with [ark-serde](https://github.com/mlange-42/ark-serde).
- Zero [dependencies](https://github.com/mlange-42/ark/blob/main/go.mod), &approx;100% [test coverage](https://badge.coveralls.io/github/mlange-42/ark?branch=main).

## Cite as

Lange, M. & contributors (2025): Ark &ndash; An archetype-based Entity Component System for Go. DOI: [10.5281/zenodo.14994239](https://doi.org/10.5281/zenodo.14994239),  GitHub repository: https://github.com/mlange-42/ark

## Contributing

Open an [issue](https://github.com/mlange-42/ark/issues)
or start a [discussion](https://github.com/mlange-42/ark/discussions)
in the {{< repo "" "GitHub repository" >}}
if you have questions, feedback, feature ideas or want to report a bug.
Pull requests are welcome.

## License

Ark and all its sources and documentation are distributed under the [MIT license](https://github.com/mlange-42/ark/blob/main/LICENSE-MIT) and the [Apache 2.0 license](https://github.com/mlange-42/ark/blob/main/LICENSE-APACHE), as your options.

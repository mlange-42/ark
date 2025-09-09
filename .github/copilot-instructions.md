# Project Overview

Ark is an archetype-based Entity Component System (ECS) for Go.

## Folder Structure

- `/ecs`: The actual Ark ecs code.
- `/ecs/generate`: Code generation templates for component mappers and filters, for internal use only.
- `/ecs/stats`: Structs representing world statistics.
- `/examples`: Stand-alone example programs demonstrating various features of Ark.
- `/docs`: Source files for the user guide at https://mlange-42.github.io/ark/, using Hugo.
- `/benchmark`: Benchmark tests for performance evaluation.

## Testing guidelines

- Unit tests are required, and are required to pass before PRs can be merged.
- Code coverage should not decrease.
- Code coverage should finally reach 100%.
- Code coverage is measured with [coveralls.io](https://coveralls.io/github/mlange-42/ark?branch=main).
- Avoid external dependencies in tests; use only Go's standard library for assertions.

## Documentation guidelines

- All user-facing code must be documented.
- All user-facing changes must be documented in CHANGELOG.md.

## Libraries and Frameworks

- The library itself has no dependencies
- The documentation website is built with [Hugo](https://gohugo.io/).
- The CI/CD pipeline uses GitHub Actions.

# Contributing to ControlPlane AI

Thank you for considering contributing to ControlPlane AI. Every contribution matters — whether it's a bug report, documentation fix, feature suggestion, or code change.

This document outlines how to contribute effectively across all ControlPlane AI projects.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)
- [License](#license)

## Code of Conduct

All contributors must follow our [Code of Conduct](./CODE_OF_CONDUCT.md). Be respectful, constructive, and professional.

## Getting Started

1. Browse [open issues](https://github.com/ravikumarve/ControlPlane-AI/issues) to find something to work on
2. Read the [PRD.md](./PRD.md) and [ARCHITECTURE.md](./ARCHITECTURE.md) for project context
3. Check the [PRODUCT-ROADMAP.md](./PRODUCT-ROADMAP.md) for upcoming priorities
4. Join the discussion in GitHub Issues before starting significant work

### First-Time Contributors

Look for issues tagged `good-first-issue` or `help-wanted`. These are curated to be approachable with clear scope.

## How to Contribute

### Reporting Bugs

Open an issue with the following template:

```markdown
**Description**: Clear description of the bug
**Steps to Reproduce**: Minimal reproduction steps
**Expected Behavior**: What should happen
**Actual Behavior**: What actually happens
**Environment**: OS, version, deployment method
```

### Suggesting Features

Open an issue with the following template:

```markdown
**Problem**: What problem does this solve?
**Solution**: Proposed implementation approach
**Alternatives**: Other approaches considered
**Priority**: How critical is this?
```

### Writing Documentation

- Documentation lives in project-specific `docs/` directories
- Uses Markdown with consistent heading hierarchy
- Include code examples where applicable
- Update the project README documentation table when adding new docs

### Submitting Code Changes

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes following the [Development Workflow](#development-workflow)
4. Submit a pull request

## Development Workflow

### Branch Naming

- `feat/<description>` — new features
- `fix/<description>` — bug fixes
- `docs/<description>` — documentation changes
- `refactor/<description>` — code refactoring
- `chore/<description>` — maintenance tasks

### Code Style

- **Go**: Follow `gofmt` and `go vet` standards
- **Python**: Follow PEP 8 with `black` formatting
- **TypeScript/JavaScript**: Follow project ESLint configuration
- **Markdown**: Use linters for consistent formatting

### Testing

- Write tests for all new functionality
- Ensure existing tests pass before submitting
- E2E tests are mandatory for security-related changes
- Run the full test suite locally before pushing

### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
```

Types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `security`

Examples:
```
feat(proxy): add rate limiting by identity
fix(audit): handle concurrent log writes correctly
docs(readme): add quickstart section
```

## Pull Request Process

1. **Title**: Must follow conventional commits format
2. **Description**: Explain what and why, not just how
3. **Related Issues**: Link to any related issues with `Closes #123`
4. **Checklist**: Verify all items before submission
5. **Review**: At least one maintainer review required for merge
6. **CI**: All CI checks must pass
7. **Size**: Prefer small, focused PRs over large monolithic changes

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests pass (`go test ./...` or equivalent)
- [ ] New code has test coverage
- [ ] Documentation updated where relevant
- [ ] CHANGELOG updated (if applicable)
- [ ] CI passes
- [ ] Reviewed for security implications

## Release Process

ControlPlane AI follows [Semantic Versioning](https://semver.org/):

- **Patch** (`1.0.x`): Bug fixes, documentation
- **Minor** (`1.x.0`): New features, backward-compatible
- **Major** (`x.0.0`): Breaking changes

Releases are managed by the maintainer. See [GOVERNANCE.md](./GOVERNANCE.md) for details.

## License

By contributing, you agree that your contributions will be licensed under the MIT License. See [LICENSE](../LICENSE).

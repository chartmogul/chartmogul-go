# Releasing chartmogul-go

## Prerequisites

- You must have push access to the repository
- `git`, `gh`, and `jq` must be installed
- Tags matching `v*` are protected by GitHub tag protection rulesets
- Releases are immutable once published (GitHub repository setting)

## Release Process

Run the release script from the repository root:

```sh
bin/release.sh <patch|minor|major>
```

The script will:

1. Verify prerequisites and that CI is green on `v4`
2. Show any open PRs targeting `v4` and ask for confirmation
3. Show PRs merged since the last tag and ask for confirmation
4. Calculate the new version based on the last tag
5. Tag the latest commit on `v4` and push the tag
6. Wait for the [release workflow](.github/workflows/release.yml) to complete, which will:
   - Run the full test suite across Go 1.21, 1.22, 1.23, and 1.24
   - Create a GitHub Release with auto-generated release notes
7. Print links to the GitHub Release and pkg.go.dev

## Changelog

Release notes are auto-generated from merged PR titles by the [release workflow](.github/workflows/release.yml). To ensure useful changelogs:

- Use clear, descriptive PR titles (e.g., "Add bulk import endpoints")
- Prefix breaking changes with `BREAKING:` so they stand out in release notes
- After the release is created, review and edit the notes on the [Releases page](https://github.com/chartmogul/chartmogul-go/releases) if needed

## Pre-release Versions

For pre-release versions, use a semver pre-release suffix:

```sh
git tag vX.Y.Z-rc1
git push origin vX.Y.Z-rc1
```

These will be automatically marked as pre-releases on GitHub.

## Security

### Repository Protections

- **Immutable releases**: Once a GitHub Release is published, its tag cannot be moved or deleted, and release assets cannot be modified
- **Tag protection rulesets**: `v*` tags cannot be deleted or force-pushed

### Go Module Security

- Go modules are fetched directly from the Git repository and verified against the [Go checksum database](https://sum.golang.org) - there is no separate package registry to compromise
- The `go.sum` file in downstream projects contains cryptographic hashes (SHA-256) for all dependencies, ensuring reproducible and tamper-evident builds
- The checksum database provides a global, append-only transparency log that prevents after-the-fact modification of published module versions

### What This Protects Against

- A compromised maintainer account cannot modify or delete existing releases
- Tags cannot be moved to point to different commits after publication
- The Go checksum database provides an independent immutability layer beyond GitHub - once a version is fetched by any user, its hash is recorded in the transparency log and cannot be changed
- Unlike package registries that rely on re-upload prevention alone, Go's checksum database cryptographically verifies that every user gets identical module contents

### Repository Settings (Admin)

These settings must be configured by a repository admin:

1. **Immutable Releases**: Settings > General > Releases > Enable "Immutable releases"
2. **Tag Protection Ruleset**: Settings > Rules > Rulesets > New ruleset targeting tags matching `v*` with deletion, force-push, and update prevention

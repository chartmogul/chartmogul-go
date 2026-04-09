# Releasing chartmogul-go

## Prerequisites

- You must have push access to the repository
- Tags matching `v*` are protected by GitHub tag protection rulesets
- Releases are immutable once published (GitHub repository setting)

## Release Process

1. Ensure all changes are merged to the `v4` branch
2. Ensure CI is green on the `v4` branch
3. Create and push a version tag:
   ```sh
   git tag v4.X.Y
   git push origin v4.X.Y
   ```
4. The [release workflow](.github/workflows/release.yml) will automatically create a GitHub Release with auto-generated release notes
5. Verify the release appears at https://github.com/chartmogul/chartmogul-go/releases
6. The Go module proxy (proxy.golang.org) will cache the new version on first fetch

## Changelog

Release notes are auto-generated from merged PR titles by the [release workflow](.github/workflows/release.yml). To ensure useful changelogs:

- Use clear, descriptive PR titles (e.g., "Add External ID field to Contact model")
- Prefix breaking changes with `BREAKING:` so they stand out in release notes
- After the release is created, review and edit the notes on the [Releases page](https://github.com/chartmogul/chartmogul-go/releases) if needed

## Pre-release Versions

For pre-release versions, use a semver pre-release suffix:

```sh
git tag v4.X.Y-rc1
git push origin v4.X.Y-rc1
```

These will be automatically marked as pre-releases on GitHub.

## Security

### Repository Protections

- **Immutable releases**: Once a GitHub Release is published, its tag cannot be moved or deleted, and release assets cannot be modified
- **Tag protection rulesets**: `v*` tags cannot be deleted or force-pushed

### Go Module Proxy

- [proxy.golang.org](https://proxy.golang.org) caches module versions permanently on first fetch
- [sum.golang.org](https://sum.golang.org) records cryptographic hashes in a tamper-evident Merkle tree
- Users running `go get` with default settings automatically verify against both services

### What This Protects Against

- A compromised maintainer account cannot modify or delete existing releases
- Tags cannot be moved to point to different commits after publication
- The Go checksum database provides an independent verification layer beyond GitHub

### Repository Settings (Admin)

These settings must be configured by a repository admin:

1. **Immutable Releases**: Settings > General > Releases > Enable "Immutable releases"
2. **Tag Protection Ruleset**: Settings > Rules > Rulesets > New ruleset targeting tags matching `v*` with deletion and force-push prevention

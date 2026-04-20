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

1. Determine the target branch from the last tag: patch/minor releases land on the current `v{major}` branch, major releases land on `v{major+1}`. If the `v{major+1}` branch does not yet exist, the script opens a [setup PR](#major-releases) that updates the module path and stops - merge that PR, complete the admin checklist in its description, then re-run `bin/release.sh major`
2. Verify prerequisites and that CI is green on the target branch
3. Show any open PRs targeting that branch and ask for confirmation
4. Show PRs merged since the last tag and ask for confirmation
5. Calculate the new version based on the last tag
6. Tag the latest commit on the target branch and push the tag
7. Wait for the [release workflow](.github/workflows/release.yml) to complete, which will:
   - Run the full test suite across Go 1.21, 1.22, 1.23, and 1.24
   - Create a GitHub Release with auto-generated release notes
8. Print links to the GitHub Release and pkg.go.dev

## Changelog

Release notes are auto-generated from merged PR titles by the [release workflow](.github/workflows/release.yml). To ensure useful changelogs:

- Use clear, descriptive PR titles (e.g., "Add bulk import endpoints")
- Prefix breaking changes with `BREAKING:` so they stand out in release notes
- After the release is created, review and edit the notes on the [Releases page](https://github.com/chartmogul/chartmogul-go/releases) if needed

## Major Releases

Go encodes the major version in the module path (e.g. `github.com/chartmogul/chartmogul-go/v5`), so a major bump requires a new `v{major+1}` branch with an updated `go.mod`.

Running `bin/release.sh major` when that branch does not exist will:

1. Push a new `v{major+1}` branch from the current major branch
2. Open a setup PR that rewrites the module path in `go.mod` and all `.go`/`.md` files
3. Exit with a pointer to the PR

The setup PR description contains a checklist of admin tasks that must be completed after merge:

- Switch the default branch to `v{major+1}`
- Update `.github/workflows/test.yml` so `branches` filters include `v{major+1}`
- Extend the branch protection ruleset to cover `v{major+1}`
- Verify the `v*` tag ruleset still applies
- Update README install snippet and badges

Once the setup PR is merged and the admin tasks are done, re-run `bin/release.sh major` to tag `v{major+1}.0.0`.

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

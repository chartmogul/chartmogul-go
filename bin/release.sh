#!/usr/bin/env bash
set -euo pipefail

# --- Usage -------------------------------------------------------------------

usage() {
  echo "Usage: bin/release.sh <patch|minor|major>"
  exit 1
}

[[ $# -eq 1 ]] || usage

BUMP_TYPE="$1"
case "$BUMP_TYPE" in
  patch|minor|major) ;;
  *) usage ;;
esac

# --- Prerequisites -----------------------------------------------------------

for cmd in git gh jq; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "Error: '$cmd' is required but not found on PATH." >&2
    exit 1
  fi
done

REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)

# --- Determine release branch ------------------------------------------------
#
# Go module versioning uses major-version branches (v4, v5, ...). Patch and
# minor releases land on the current major branch; major releases land on a
# new v{major+1} branch, which must already exist (it requires a module path
# change in go.mod and so cannot be created by this script).

LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [[ -z "$LAST_TAG" ]]; then
  echo "Error: No existing tag found - cannot determine current major version." >&2
  exit 1
fi

CURRENT_VERSION="${LAST_TAG#v}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

if [[ "$BUMP_TYPE" == "major" ]]; then
  TARGET_MAJOR=$((MAJOR + 1))
else
  TARGET_MAJOR="$MAJOR"
fi

DEFAULT_BRANCH="v${TARGET_MAJOR}"
SOURCE_BRANCH="v${MAJOR}"

# --- Bootstrap new major branch (one-time setup for major releases) ----------
#
# A major release targets a new v{N+1} branch, which must have the module path
# in go.mod and all Go imports updated. If that branch doesn't exist yet, the
# script opens a setup PR with those changes and stops - the PR must be merged
# and the branch made the default before a tag can be pushed.

if ! git ls-remote --exit-code --heads origin "$DEFAULT_BRANCH" &>/dev/null; then
  if [[ "$BUMP_TYPE" != "major" ]]; then
    echo "Error: Branch '${DEFAULT_BRANCH}' does not exist on origin." >&2
    exit 1
  fi

  echo ""
  echo "Branch '${DEFAULT_BRANCH}' does not exist. A major release requires a new"
  echo "major-version branch with the module path updated."
  echo ""
  echo "This script will:"
  echo "  1. Push a new '${DEFAULT_BRANCH}' branch from origin/${SOURCE_BRANCH}"
  echo "  2. Open a PR that updates the module path from /${SOURCE_BRANCH} to /${DEFAULT_BRANCH} in go.mod and all .go/.md files"
  echo "  3. Exit. You then need to merge the PR, have admins apply branch protections, make '${DEFAULT_BRANCH}' the default branch, then re-run 'bin/release.sh major'."
  echo ""
  read -rp "Proceed? [y/N] " CONFIRM
  if [[ ! "$CONFIRM" =~ ^[yY]$ ]]; then
    echo "Aborted."
    exit 0
  fi

  STASHED=false
  if ! git diff --quiet || ! git diff --cached --quiet; then
    git stash --include-untracked --quiet
    STASHED=true
  fi
  restore_stash() {
    if [[ "$STASHED" == true ]]; then
      git stash pop --quiet 2>/dev/null || true
      STASHED=false
    fi
  }
  trap restore_stash EXIT

  git fetch origin "$SOURCE_BRANCH"

  # Create the new major branch on origin, pointing at source HEAD.
  git push origin "origin/${SOURCE_BRANCH}:refs/heads/${DEFAULT_BRANCH}"
  echo "Created '${DEFAULT_BRANCH}' on origin."

  SETUP_BRANCH="setup/${DEFAULT_BRANCH}-module-path"
  git checkout -B "$SETUP_BRANCH" "origin/${SOURCE_BRANCH}"

  OLD_PATH="github.com/chartmogul/chartmogul-go/${SOURCE_BRANCH}"
  NEW_PATH="github.com/chartmogul/chartmogul-go/${DEFAULT_BRANCH}"

  # Update module path in go.mod and all references in .go / .md files.
  # Using find with -print0 so paths with spaces are handled; LC_ALL=C keeps sed portable.
  while IFS= read -r -d '' f; do
    LC_ALL=C sed -i.bak "s|${OLD_PATH}|${NEW_PATH}|g" "$f" && rm -f "${f}.bak"
  done < <(find . \( -name "*.go" -o -name "*.md" -o -name "go.mod" \) -not -path "./vendor/*" -not -path "./.git/*" -print0)

  if git diff --quiet; then
    echo "Error: No references to '${OLD_PATH}' found - nothing to update." >&2
    exit 1
  fi

  git add -A
  git commit -m "Update module path from /${SOURCE_BRANCH} to /${DEFAULT_BRANCH}"
  git push -u origin "$SETUP_BRANCH"

  PR_BODY=$(cat <<EOF
## Summary

Bootstrap the \`${DEFAULT_BRANCH}\` branch for the upcoming ${DEFAULT_BRANCH}.0.0 major release. Updates the module path from \`${OLD_PATH}\` to \`${NEW_PATH}\` in \`go.mod\` and all Go imports and documentation references.

This PR is opened automatically by \`bin/release.sh major\` when the \`${DEFAULT_BRANCH}\` branch does not yet exist.

## Admin checklist (to be done after merge)

- [ ] **Update default branch**: Settings > General > Default branch > switch from \`${SOURCE_BRANCH}\` to \`${DEFAULT_BRANCH}\`
- [ ] **Branch protection ruleset**: Settings > Rules > Rulesets > add/extend ruleset to target \`${DEFAULT_BRANCH}\` with require-PR, require-status-checks, force-push prevention, and deletion prevention
- [ ] **Tag protection ruleset**: confirm the existing \`v*\` tag ruleset still applies (it should, but verify after the default branch change)
- [ ] **pkg.go.dev**: after the first \`${DEFAULT_BRANCH}.0.0\` tag is pushed, visit https://pkg.go.dev/${NEW_PATH} to trigger indexing

## Next steps

After this PR is merged and the admin checklist above is complete, re-run:

\`\`\`sh
bin/release.sh major
\`\`\`

which will tag \`${DEFAULT_BRANCH}.0.0\` on the \`${DEFAULT_BRANCH}\` branch.
EOF
)

  PR_URL=$(gh pr create \
    --base "$DEFAULT_BRANCH" \
    --head "$SETUP_BRANCH" \
    --title "Update module path to /${DEFAULT_BRANCH}" \
    --body "$PR_BODY" \
    | tail -1)

  echo ""
  echo "Setup PR created: ${PR_URL}"
  echo ""
  echo "Next steps:"
  echo "  1. Review and merge the PR"
  echo "  2. Complete the admin checklist in the PR description"
  echo "  3. Re-run 'bin/release.sh major' to tag ${DEFAULT_BRANCH}.0.0"

  git checkout "$SOURCE_BRANCH" 2>/dev/null || true
  restore_stash
  trap - EXIT
  exit 0
fi

# --- Check CI is green -------------------------------------------------------

echo "Checking CI status on ${DEFAULT_BRANCH}..."
LAST_RUN=$(gh run list --branch "$DEFAULT_BRANCH" --workflow test.yml --limit 1 --json conclusion,url)
CONCLUSION=$(echo "$LAST_RUN" | jq -r '.[0].conclusion // empty')
RUN_URL=$(echo "$LAST_RUN" | jq -r '.[0].url // empty')

if [[ "$CONCLUSION" != "success" ]]; then
  echo "Error: Latest CI run on ${DEFAULT_BRANCH} is not green (status: ${CONCLUSION:-unknown})." >&2
  [[ -n "$RUN_URL" ]] && echo "  $RUN_URL" >&2
  exit 1
fi
echo "CI is green."

# --- Show open PRs -----------------------------------------------------------

OPEN_PRS=$(gh pr list --base "$DEFAULT_BRANCH" --json number,title,url)
PR_COUNT=$(echo "$OPEN_PRS" | jq 'length')

if [[ "$PR_COUNT" -gt 0 ]]; then
  echo ""
  echo "There are $PR_COUNT open PR(s) targeting ${DEFAULT_BRANCH}:"
  echo "$OPEN_PRS" | jq -r '.[] | "  #\(.number) \(.title)\n    \(.url)"'
  echo ""
  read -rp "Continue releasing? [y/N] " CONFIRM
  if [[ ! "$CONFIRM" =~ ^[yY]$ ]]; then
    echo "Aborted."
    exit 0
  fi
fi

# --- Show PRs included in this release ----------------------------------------

echo ""
echo "PRs merged since ${LAST_TAG}:"
MERGED_PRS=$(gh pr list --base "$DEFAULT_BRANCH" --state merged --search "merged:>=$(git log -1 --format=%aI "$LAST_TAG")" --json number,title,url)

MERGED_COUNT=$(echo "$MERGED_PRS" | jq 'length')
if [[ "$MERGED_COUNT" -eq 0 ]]; then
  echo "  (none)"
else
  echo "$MERGED_PRS" | jq -r '.[] | "  #\(.number) \(.title)\n    \(.url)"'
fi

echo ""
read -rp "Release these changes as ${BUMP_TYPE}? [y/N] " CONFIRM
if [[ ! "$CONFIRM" =~ ^[yY]$ ]]; then
  echo "Aborted."
  exit 0
fi

# --- Calculate new version ----------------------------------------------------

case "$BUMP_TYPE" in
  major) NEW_VERSION="${TARGET_MAJOR}.0.0" ;;
  minor) NEW_VERSION="${MAJOR}.$((MINOR + 1)).0" ;;
  patch) NEW_VERSION="${MAJOR}.${MINOR}.$((PATCH + 1))" ;;
esac

TAG="v${NEW_VERSION}"

echo ""
echo "Tagging: ${LAST_TAG} -> ${TAG}"

# --- Ensure we're on latest default branch ------------------------------------

git checkout "$DEFAULT_BRANCH"
git pull

# --- Tag and push -------------------------------------------------------------

git tag "$TAG"
git push origin "$TAG"

echo "Tag ${TAG} pushed."
echo "Waiting for release workflow..."

# --- Poll for release CI ------------------------------------------------------

interrupt_ci_wait() {
  echo ""
  echo "Interrupted. The tag ${TAG} has been pushed and the release workflow is running."
  echo "  Workflow: https://github.com/${REPO}/actions/workflows/release.yml"
  echo "  GitHub:   https://github.com/${REPO}/releases/tag/${TAG}"
  exit 0
}

trap interrupt_ci_wait INT

sleep 5 # give GitHub a moment to register the run
while true; do
  RUN=$(gh run list --branch "$TAG" --workflow release.yml --limit 1 --json status,conclusion,url)
  STATUS=$(echo "$RUN" | jq -r '.[0].status // empty')
  RUN_CONCLUSION=$(echo "$RUN" | jq -r '.[0].conclusion // empty')
  RELEASE_RUN_URL=$(echo "$RUN" | jq -r '.[0].url // empty')

  if [[ "$STATUS" == "completed" ]]; then
    echo ""
    if [[ "$RUN_CONCLUSION" == "success" ]]; then
      echo "Release workflow completed successfully."
    else
      echo "Release workflow finished with status: ${RUN_CONCLUSION}" >&2
    fi
    [[ -n "$RELEASE_RUN_URL" ]] && echo "  $RELEASE_RUN_URL"
    break
  fi
  printf "."
  sleep 10
done

trap - INT

# --- Summary ------------------------------------------------------------------

echo ""
echo "Release ${TAG} complete!"
echo "  GitHub:  https://github.com/${REPO}/releases/tag/${TAG}"
echo "  pkg.go:  https://pkg.go.dev/github.com/${REPO}/${DEFAULT_BRANCH}@${TAG}"

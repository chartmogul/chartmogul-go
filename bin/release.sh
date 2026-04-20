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

if ! git ls-remote --exit-code --heads origin "$DEFAULT_BRANCH" &>/dev/null; then
  echo "Error: Branch '${DEFAULT_BRANCH}' does not exist on origin." >&2
  if [[ "$BUMP_TYPE" == "major" ]]; then
    echo "  For a major release, create the '${DEFAULT_BRANCH}' branch first with the updated" >&2
    echo "  module path in go.mod (github.com/chartmogul/chartmogul-go/${DEFAULT_BRANCH})." >&2
  fi
  exit 1
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

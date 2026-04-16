#!/usr/bin/env bash
set -euo pipefail

DEFAULT_BRANCH="v4"

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

LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [[ -n "$LAST_TAG" ]]; then
  echo ""
  echo "PRs merged since ${LAST_TAG}:"
  MERGED_PRS=$(gh pr list --base "$DEFAULT_BRANCH" --state merged --search "merged:>=$(git log -1 --format=%aI "$LAST_TAG")" --json number,title,url)
else
  echo ""
  echo "PRs merged (no previous tag found, showing recent):"
  MERGED_PRS=$(gh pr list --base "$DEFAULT_BRANCH" --state merged --limit 10 --json number,title,url)
fi

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

CURRENT_VERSION="${LAST_TAG#v}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

case "$BUMP_TYPE" in
  major) NEW_VERSION="$((MAJOR + 1)).0.0" ;;
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

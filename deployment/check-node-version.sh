#!/usr/bin/env bash
# The app-builder image installs Node.js before the repository is available, so
# it cannot read the pin from package.json. Keep its literal in sync with
# devEngines.runtime, which is the single source of truth everywhere else.
set -euo pipefail

cd "$(dirname "$0")/.."

dockerfile=deployment/app-builder/Dockerfile
pinned=$(node -p "require('./package.json').devEngines.runtime.version.replace(/^[^0-9]*/, '')")
actual=$(sed -n 's/^ENV NODE_VERSION \(.*\)$/\1/p' "$dockerfile")

if [ "$actual" != "$pinned" ]; then
  echo "$dockerfile pins Node.js $actual, but package.json pins $pinned" >&2
  exit 1
fi

echo "Node.js $pinned is pinned consistently."

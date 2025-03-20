#!/bin/bash
set -e

# Copy hooks to the correct directory
mkdir -p .git/hooks
cp .githooks/* .git/hooks/
chmod +x .git/hooks/*

git config core.hooksPath .git/hooks

echo "Git hooks installed!"
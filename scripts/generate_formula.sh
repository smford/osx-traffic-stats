#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-1.0.0}"
TAG_NAME="${2:-v${VERSION}}"
DIST_DIR="${3:-bin}"
OUT_FILE="${4:-${DIST_DIR}/osx-traffic-stats.rb}"

# Clean version (strip leading v if present)
CLEAN_VERSION="${VERSION#v}"

get_sha() {
  local file="$1"
  if [ -f "$file" ]; then
    if command -v sha256sum >/dev/null 2>&1; then
      sha256sum "$file" | cut -d' ' -f1
    elif command -v shasum >/dev/null 2>&1; then
      shasum -a 256 "$file" | cut -d' ' -f1
    fi
  else
    echo "REPLACE_WITH_SHA256"
  fi
}

# Look for universal release archive in DIST_DIR
TAR_FILE="${DIST_DIR}/osx-traffic-stats-${TAG_NAME}-darwin-universal.tar.gz"
if [ ! -f "$TAR_FILE" ]; then
  TAR_FILE="${DIST_DIR}/osx-traffic-stats-v${CLEAN_VERSION}-darwin-universal.tar.gz"
fi

SHA_DARWIN_UNIVERSAL=$(get_sha "$TAR_FILE")

mkdir -p "$(dirname "$OUT_FILE")"

cat <<EOF > "$OUT_FILE"
# typed: false
# frozen_string_literal: true

# This formula was auto-generated for osx-traffic-stats (https://github.com/smford/osx-traffic-stats).
class OsxTrafficStats < Formula
  desc "Real-time macOS menu bar network traffic monitor"
  homepage "https://github.com/smford/osx-traffic-stats"
  version "${CLEAN_VERSION}"
  license "MIT"
  depends_on :macos

  on_macos do
    on_arm do
      url "https://github.com/smford/osx-traffic-stats/releases/download/v#{version}/osx-traffic-stats-v#{version}-darwin-universal.tar.gz"
      sha256 "${SHA_DARWIN_UNIVERSAL}"
    end
    on_intel do
      url "https://github.com/smford/osx-traffic-stats/releases/download/v#{version}/osx-traffic-stats-v#{version}-darwin-universal.tar.gz"
      sha256 "${SHA_DARWIN_UNIVERSAL}"
    end
  end

  def install
    bin.install "osx-traffic-stats"
  end

  test do
    assert_match "OSX Traffic Stats", shell_output("#{bin}/osx-traffic-stats -v")
  end
end
EOF

echo "Generated Homebrew formula at ${OUT_FILE}"

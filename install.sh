#!/bin/sh
set -e

REPO="emoral435/time-broker"
BIN_DIR="${TIME_BROKER_DIR:-$HOME/.time-broker/bin}"
LINK_DIR="${LINK_DIR:-}"

if [ -z "$LINK_DIR" ]; then
  case ":$PATH:" in
    *":$HOME/.local/bin:"*) LINK_DIR="$HOME/.local/bin" ;;
    *) LINK_DIR="/usr/local/bin" ;;
  esac
fi

BIN_PATH="$BIN_DIR/time-broker"
TB_PATH="$BIN_DIR/tb"
LINK_PATH="$LINK_DIR/time-broker"
TB_LINK_PATH="$LINK_DIR/tb"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  darwin|linux) ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"

if [ -z "$VERSION" ]; then
  echo "No release found. Installing from source via 'go install'..."
  if command -v go >/dev/null 2>&1; then
    go install "github.com/${REPO}@latest"
    echo "Installed. Ensure $(go env GOPATH)/bin is in your PATH."
    exit 0
  else
    echo "Go is not installed and no pre-built release is available."
    echo "Install Go from https://go.dev/dl/ and try again."
    exit 1
  fi
fi

FILENAME="time-broker-${VERSION}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading time-broker ${VERSION} for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o "${TMPDIR}/${FILENAME}"
tar xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"

if ! mkdir -p "$BIN_DIR"; then
  echo "Could not create install directory: $BIN_DIR"
  exit 1
fi

mv "${TMPDIR}/time-broker" "$BIN_PATH"
chmod 755 "$BIN_PATH" 2>/dev/null || true

ln -sf "time-broker" "$TB_PATH"

resolve_path() {
  (cd "$1" 2>/dev/null && pwd -P)
}

REAL_BIN_DIR="$(resolve_path "$BIN_DIR")"
REAL_LINK_DIR="$(resolve_path "$LINK_DIR" 2>/dev/null || echo "")"

if [ -n "$REAL_BIN_DIR" ] && [ "$REAL_BIN_DIR" = "$REAL_LINK_DIR" ]; then
  echo "Install dir and link dir resolve to the same path; skipping symlink."
else
  if [ -w "$LINK_DIR" ] || (mkdir -p "$LINK_DIR" 2>/dev/null && [ -w "$LINK_DIR" ]); then
    rm -f "$LINK_PATH" "$TB_LINK_PATH"
    ln -s "$BIN_PATH" "$LINK_PATH"
    ln -s "$TB_PATH" "$TB_LINK_PATH"
  else
    echo "Linking to ${LINK_DIR} (requires sudo)..."
    sudo mkdir -p "$LINK_DIR"
    sudo rm -f "$LINK_PATH" "$TB_LINK_PATH"
    sudo ln -s "$BIN_PATH" "$LINK_PATH"
    sudo ln -s "$TB_PATH" "$TB_LINK_PATH"
  fi
fi

echo "time-broker ${VERSION} installed to ${BIN_DIR}"
echo "Commands: time-broker, tb"

case ":$PATH:" in
  *":$LINK_DIR:"*) ;;
  *) echo "Add ${LINK_DIR} to your PATH and restart your terminal." ;;
esac

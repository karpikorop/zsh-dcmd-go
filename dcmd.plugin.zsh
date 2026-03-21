# Zsh wrapper for the Go-based dcmd
DCMD_PLUGIN_DIR="${0:A:h}"
DCMD_BIN="${DCMD_PLUGIN_DIR}/dcmd"
DCMD_REPO="karpikorop/zsh-dcmd-go"

_dcmd_download_binary() {
    local os arch binary_name url

    case "$(uname -s)" in
        Darwin) os="darwin" ;;
        Linux)  os="linux" ;;
        *)      return 1 ;;
    esac

    case "$(uname -m)" in
        x86_64)       arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *)            return 1 ;;
    esac

    binary_name="dcmd-${os}-${arch}"
    url="https://github.com/${DCMD_REPO}/releases/latest/download/${binary_name}"

    echo "dcmd: downloading binary for ${os}/${arch}..."
    if curl -fsSL -o "$DCMD_BIN" "$url"; then
        chmod +x "$DCMD_BIN"
        echo "dcmd: binary installed successfully."
        return 0
    else
        rm -f "$DCMD_BIN"
        return 1
    fi
}

_dcmd_build_from_source() {
    if ! command -v go &>/dev/null; then
        return 1
    fi

    if [[ ! -f "${DCMD_PLUGIN_DIR}/cmd/dcmd/main.go" ]]; then
        return 1
    fi

    echo "dcmd: building from source..."
    (cd "$DCMD_PLUGIN_DIR" && go build -o dcmd ./cmd/dcmd)
}

_dcmd_ensure_binary() {
    if [[ -f "$DCMD_BIN" ]]; then
        return 0
    fi

    if _dcmd_download_binary; then
        return 0
    fi

    if _dcmd_build_from_source; then
        return 0
    fi

    echo "dcmd: could not install binary." >&2
    echo "  Install a pre-built binary from: https://github.com/${DCMD_REPO}/releases" >&2
    echo "  Or install Go and re-source this plugin to build from source." >&2
    return 1
}

dcmd() {
    if ! _dcmd_ensure_binary; then
        return 1
    fi

    local cmd_to_run
    cmd_to_run=$("$DCMD_BIN" "$@")

    if [[ -n "$cmd_to_run" ]]; then
        print -z "$cmd_to_run"
    fi
}

# dcmd: the binary writes the chosen command to DCMD_BUFFER_FILE. This function puts it on the prompt for review.
dcmd() {
    local buffer_file exit_code
    buffer_file=$(mktemp "${TMPDIR:-/tmp}/dcmd.XXXXXX") || return 1

    {
        DCMD_BUFFER_FILE="$buffer_file" __DCMD_BIN__ "$@"
        exit_code=$?

        if [[ -s "$buffer_file" ]]; then
            print -z -- "$(<"$buffer_file")"
        fi
    } always {
        rm -f -- "$buffer_file"
    }

    return $exit_code
}

#!/bin/bash
set -euo pipefail

readonly THIS_SCRIPT_DIR=$(
    cd $(dirname $0)
    pwd
)

readonly GOMD5SUM="${THIS_SCRIPT_DIR}/bin/gomd5sum"

function test_standard() {
    local tmpdir
    tmpdir=$(mktemp -d "${THIS_SCRIPT_DIR}/XXXXXX")
    tmpA=$(mktemp "${tmpdir}/XXXXXX")
    tmpB=$(mktemp "${tmpdir}/XXXXXX")
    tmpC=$(mktemp "${tmpdir}/XXXXXX")

    local got
    got=$(${GOMD5SUM} "${tmpA}" "${tmpB}" "${tmpC}")

    local want
    want=$(md5sum "${tmpA}" "${tmpB}" "${tmpC}")

    if [[ "${got}" != "${want}" ]]; then
        echo "[ERRO] There is a difference between got and want."
    else
        echo "[INFO] OK"
    fi
}

function main() {
    test_standard
}

main "$@"

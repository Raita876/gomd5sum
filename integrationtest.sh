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

    if ! diff -u <(${GOMD5SUM} "${tmpA}" "${tmpB}" "${tmpC}") <(md5sum "${tmpA}" "${tmpB}" "${tmpC}"); then
        echo "[test_standard] NG: There is a difference between got and want."
        exit 1
    else
        echo "[test_standard] OK"
    fi

    rm -r "${tmpdir}"
}

function main() {
    test_standard
}

main "$@"

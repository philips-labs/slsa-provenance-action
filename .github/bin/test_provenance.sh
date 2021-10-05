#!/usr/bin/env bash

FIRST_ARTIFACT=$1
SECOND_ARTIFACT=$2

cp "$FIRST_ARTIFACT" /tmp/first_artifact
cp "$SECOND_ARTIFACT" /tmp/second_artifact

diff_artifacts() {
    (grep -v "buildInvocationId" | grep -vF "buildFinishedOn") > /tmp/first_artifact
    (grep -v "buildInvocationId" | grep -vF "buildFinishedOn") > /tmp/second_artifact
    cat /tmp/first_artifact
    echo "--------------"
    cat /tmp/second_artifact
    echo "--------------"
    diff -f /tmp/first_artifact /tmp/second_artifact > /dev/null
    local exit_code=$?
    rm -rf /tmp/{first_artifact,second_artifact}
    exit $exit_code
}

diff_artifacts
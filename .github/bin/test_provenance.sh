#!/usr/bin/env bash

FIRST_ARTIFACT=$1
SECOND_ARTIFACT=$2

diff_artifacts() {
    (grep -v "buildInvocationId" | grep -vF "buildFinishedOn1") > /tmp/"$FIRST_ARTIFACT" < "$FIRST_ARTIFACT"
    (grep -v "buildInvocationId" | grep -vF "buildFinishedOn1")  > /tmp/"$SECOND_ARTIFACT" < "$SECOND_ARTIFACT"
    cat /tmp/"$FIRST_ARTIFACT"
    echo "--------------"
    cat /tmp/"$SECOND_ARTIFACT"
    echo "--------------"
    diff -f /tmp/"$FIRST_ARTIFACT" /tmp/"$SECOND_ARTIFACT" > /dev/null
    local exit_code=$?
    rm -rf /tmp/{"$FIRST_ARTIFACT","$SECOND_ARTIFACT"}
    exit $exit_code
}

diff_artifacts
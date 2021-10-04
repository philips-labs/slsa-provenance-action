#!/bin/bash

FIRST_ARTIFACT=$1
SECOND_ARTIFACT=$2
# Comma seperated line of text of attributes to skip, examples are: buildFinishedOn,sha1
# SKIP_LINES=$3

diff_artifacts() {
    comm -3 $FIRST_ARTIFACT $SECOND_ARTIFACT > difference.txt
    grep -vF "buildFinishedOn" difference.txt > difference2.txt
    grep -vF "sha1" difference2.txt > difference3.txt
    cat difference3.txt
    # cat difference.txt | grep -vF "buildFinishedOn" | grep -vF "sha1" | grep -vF "buildInvocationId" | grep -vE "^c" | grep -vE "^."
}

diff_artifacts
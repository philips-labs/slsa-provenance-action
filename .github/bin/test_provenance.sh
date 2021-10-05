#!/usr/bin/env bash

FIRST_ARTIFACT=$1
SECOND_ARTIFACT=$2

diff_artifacts() {
    local temp_file1=/tmp/file1
    local temp_file2=/tmp/file2
    (grep -v "buildInvocationId" | grep -v "buildFinishedOn" | grep -v "sha1") >$temp_file1 <"$FIRST_ARTIFACT"
    (grep -v "buildInvocationId" | grep -v "buildFinishedOn" | grep -v "sha1") >$temp_file2 <"$SECOND_ARTIFACT"
    echo "-----step3---------"
    cat $temp_file1
    echo "-----step4---------"
    cat $temp_file2
    echo "-----step5---------"
    diff -wf $temp_file1 $temp_file2 >difference.txt #>/dev/null
    cat difference.txt
    local exit_code=$?
    rm -rf /tmp/{file1,file2}
    exit $exit_code
}

diff_artifacts

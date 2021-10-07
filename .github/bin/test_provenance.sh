#!/usr/bin/env bash


diff_artifacts() {
  if [ ! -f "$1" ] || [ ! -f "$2" ] ; then
    echo "Please provide two files to compare" >&2
    exit 1
  fi

  local temp_file1=/tmp/file1
  local temp_file2=/tmp/file2
  (grep -v "buildInvocationId" | grep -v "buildFinishedOn" | grep -v "sha1") > $temp_file1 < "$FIRST_ARTIFACT"
  (grep -v "buildInvocationId" | grep -v "buildFinishedOn" | grep -v "sha1") > $temp_file2 < "$SECOND_ARTIFACT"

  diff -wf $temp_file1 $temp_file2 > /dev/null
  local exit_code=$?
  rm -rf /tmp/{file1,file2}
  exit $exit_code
}

diff_artifacts "$1" "$2"

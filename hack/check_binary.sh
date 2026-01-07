#!/usr/bin/env bash


binary="$1"
version_arg="$2"
version="$3"

function not_empty() {
  if [ -z "$2" ]; then
    echo "$1 is empty"
    exit 1
  fi

  return 0
}

not_empty "binary" "$binary"
not_empty "version_arg" "$version_arg"
not_empty "version" "$version"

full_path="$(pwd)/bin/${binary}"

if [ ! -x "$full_path" ]; then
  echo "$binary_full_path not exists or not executable"
  exit 1
fi

if ! "$full_path" "$version_arg" | grep -q "$version" ; then
  echo "$full_path version not match ${version}. Version is $("$full_path" "$version_arg")"
  exit 1
fi

exit 0

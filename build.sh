#!/usr/bin/env bash

build() {
  suffix=".go"

  for f in "$@"
  do
    binary_name=${f%$suffix}
    echo "= Building $binary_name"
    go build -o "build/dicam-$binary_name" -ldflags "-s" $f
  done

  printf "= Built done!\n\n"
  ls -lah build/*
}

cd src/
build *.go

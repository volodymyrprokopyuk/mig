#!/usr/bin/env fish

function mig_lint
  echo -n "=> linting mig "
  golangci-lint run ./... || return 1
end

function mig_build
  echo "=> updating mig"
  go get -u ./... && go mod tidy
  echo "=> building mig"
  CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/mig ./cmd/mig
end

function main
  set cmd $argv[1]
  set sub $argv[2]
  switch $cmd
  case lint
    mig_lint
  case build
    mig_build
  case '*'
    echo "unknown command $cmd" && return 1
  end
end

main $argv

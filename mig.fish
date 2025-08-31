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

function mig_tag
  set git_status (git status --porcelain)
  if test "$git_status" != ""
    echo "uncommitted changes"
    return 1
  end
  echo "=> tagging $MIG_VERSION"
  git tag $MIG_VERSION
  # git push --tags
  # git tag -d 0.0.0
  # git push origin --delete tag 0.0.0
end

function main
  set cmd $argv[1]
  set sub $argv[2]
  switch $cmd
  case lint
    mig_lint
  case build
    mig_build
  case tag
    mig_tag
  case '*'
    echo "unknown command $cmd" && return 1
  end
end

main $argv

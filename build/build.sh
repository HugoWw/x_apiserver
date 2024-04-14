#!/bin/bash

######################################################
# $Name:        build.sh
# $Version:     v0.1
# $Function:    build program script
# $Author:      gouhuan
# $Create Date: 2024-3-24
# $Description: shell
######################################################

## constant variable
readonly APPVersion="github.com/HugoWw/x_apiserver/cmd/x_apiserver/options.APPVersion"
readonly GitCommit="github.com/HugoWw/x_apiserver/cmd/x_apiserver/options.GitCommit"
readonly BuildDate="github.com/HugoWw/x_apiserver/cmd/x_apiserver/options.BuildDate"

# Declare variables
CURRENT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd -P)
SWAG_GEN=${CURRENT_DIR}/swagger-gen.go
SWAG_PATH=${CURRENT_DIR}/
MOUDLE=$1
PARAMETE_TAG=$2
REG_REPO="registry.harbor.com/x-apiserver"

# Log Func
log() {
  local type="$1"
  shift
  # accept argument string or stdin
  local text="$*"
  if [ "$#" -eq 0 ]; then text="$(cat)"; fi
  local dt
  dt="$(date "+%Y-%m-%d %H:%M:%S")"
  printf '%s [%s] [build.sh]: %s\n' "$dt" "$type" "$text"
}
info() {
  log INFO "$@"
}
warn() {
  log WARN "$@" >&2
}
error() {
  log ERROR "$@" >&2
  exit 1
}

# Usage info
Usage() {
  cat <<EOF
Execute build.sh script build programs.

usage: ${0} [OPTIONS]

The following flags are required.

    --build-apiserver           Only build apiserver program
    --build-apiserver-img  tag  Build apiserver image(if 'tag' is empty,the default is latest)
EOF
  exit 1
}

BuildSwagDoc(){
  go run ${SWAG_GEN} --doc-dir ${SWAG_PATH}
}

BuildSwagImage(){
  local tag
  local image_tag

  if [ -z ${PARAMETE_TAG} ];then
    tag="latest"
  else
    tag=${PARAMETE_TAG}
  fi

  image_tag=${REG_REPO}/x-apiserver-swagger:${tag}
  cd ${CURRENT_DIR}/
  info "======== Current Dir: $(pwd), Start build Swagger ApiServer image: ${image_tag}"


  docker build -t ${image_tag} -f ./Dockerfile.swag .
  if [ $? = 0 ];then
       info "====== Build Swagger ApiServer image Successfully! Current Dir: $(pwd)"
  else
       error "====== Build Swagger ApiServer image Failed! Current Dir: $(pwd)"
  fi
}

BuildApiServer() {
  local cmd_go
  local build_flag
  local build_time
  local git_commit

  git_commit=git-$(git rev-parse --short HEAD)
  git_branch=$(git rev-parse --abbrev-ref HEAD) # git version below 2.22
  build_time=$(date +'%Y-%m-%d-%H:%M:%S')
  cmd_go=${CURRENT_DIR}/apiserver/x_apiserver
  build_flag="-X ${GitCommit}=${git_commit} -X ${BuildDate}=${build_time} -X ${APPVersion}=${git_branch}"

  cd ${CURRENT_DIR}/../cmd/x_apiserver
  info "======== Current Dir: $(pwd), Start build X-APIServer......."

  # 编译x-apiserver二进制
  export CGO_ENABLED="1";export GO111MODULE=on;export GOPROXY='https://goproxy.cn,direct';export GOSUMDB=off;export GONOSUMDB=off;go mod tidy
  export CGO_ENABLED="1";export GO111MODULE=on;export GOPROXY='https://goproxy.cn,direct';export GOSUMDB=off;export GONOSUMDB=off;go build -ldflags="${build_flag}" -o "${cmd_go}"
  if [ $? = 0 ]; then
    cd ${CURRENT_DIR}/
    info "====== Build X-APIServer Successfully! Return Script Dir,Current Dir: $(pwd)"
  else
    cd ${CURRENT_DIR}/
    error "====== Build X-APIServer Failed! Return Script Dir,Current Dir: $(pwd)"
  fi
}

BuildApiServerImage() {
  local tag
  local image_tag

  if [ -z ${PARAMETE_TAG} ]; then
    tag="latest"
  else
    tag=${PARAMETE_TAG}
  fi

  image_tag=${REG_REPO}/x-apiserver:${tag}
  cd ${CURRENT_DIR}/apiserver
  info "======== Current Dir: $(pwd), Start build x-apiserver docker image: ${image_tag}"

  git_commit=git-$(git rev-parse --short HEAD)
  docker build --build-arg GIT_REV=${git_commit} -t ${image_tag} -f ./Dockerfile .
  if [ $? = 0 ]; then
    info "====== Build x-apiserver docker image Successfully! Current Dir: $(pwd)"
  else
    error "====== Build x-apiserver docker image Failed! Current Dir: $(pwd)"
  fi
}


BuildAllImage(){
  BuildApiServer
  BuildApiServerImage
  BuildSwagDoc
  BuildSwagImage
}


main() {
  case ${MOUDLE} in
  --build-apiserver)
    BuildApiServer
    BuildSwagDoc
    ;;
  --build-apiserver-img)
    BuildApiServerImage
    BuildSwagImage
    ;;
  *)
    Usage
    ;;
  esac
}

main

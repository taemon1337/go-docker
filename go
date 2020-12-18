#!/bin/bash
OWD=$(pwd)
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

GOPATH=$HOME/code/go
WORKDIR=/go/src
VOLS="-v ${HOME}/.cache:/go/cache"
DOCKUSER="--user $(id -u):$(id -g)"
DOCK="-it --rm"
NET="--net host"
ENV="-e USER=gopher -e HOME=$HOME -e GOCACHE=/go/cache"
TEST=""
DEBUG=""
CMD=""

# if not in docker group, run sudo docker instead of docker
if $(groups $USER | grep docker) ; then 
  DOCKER="docker"
else
  DOCKER="docker"
fi

_echo() {
  if [ -n "$DEBUG" ]; then
    echo "[DEBUG] ${@}"
  fi
}

_main() {
  GOLANG_VERSION="${GOLANG_VERSION:-1.15}"
  GOLANG="${GOLANG:-golang:$GOLANG_VERSION}"

  if [ -n "$WORKDIR" ]; then DOCK="$DOCK -w $WORKDIR" ; fi
  if [ -n "$TEST" ]; then DEBUG=1; fi # if test then debug as well
  if [ -n "$GOLANG_ALPINE" ]; then GOLANG="$GOLANG-alpine"; fi
  if [ -n "$GOPATH" ]; then ENV+=" -e GOPATH=/go" ; VOLS="$VOLS -v $GOPATH:/go:rw" ; fi
  if [ -n "$GOOS" ]; then ENV+=" -e GOOS=$GOOS" ; fi

  local gocmd="$CMD"
  if [ "$gocmd" == "" ]; then gocmd="go ${@}" ; fi

  local cmd="$DOCKER run $ENV $DOCK $NET $VOLS $GOLANG $gocmd"
  _echo "$cmd"

  if [ -n "$TEST" ]; then
    read -p "Do you want to execute the above command? " -n 1 -r
    if [[ $REPLY =~ "^[Yy]$" ]] ; then
      return 1
    fi
  fi

  $cmd
  return $?
}

# parse args
args=()
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -v|--verbose|--debug)
      DEBUG="1"
      shift
      ;;
    --golang)
      GOLANG="$2"
      shift
      shift
      ;;
    --golang-version)
      GOLANG_VERSION="$2"
      shift
      shift
      ;;
    --gopath)
      GOPATH="$2"
      shift
      shift
      ;;
    --goos)
      GOOS="$2"
      shift
      shift
      ;;
    --alpine|--golang-alpine)
      GOLANG_ALPINE="1"
      shift
      shift
      ;;
    -t|--dry-run|--test)
      TEST="1"
      shift
      shift
      ;;
    -n|--net)
      NET="$2"
      shift
      shift
      ;;
    -e|--env)
      ENV="$ENV -e $2"
      shift
      shift
      ;;
    -u|--user)
      DOCKUSER="$2"
      shift
      shift
      ;;
    -v|--volume)
      VOLS="$VOLS -v $2"
      shift
      shift
      ;;
    -w|--workdir)
      WORKDIR="$2"
      shift
      shift
      ;;
    --sh|--shell|--bash)
      CMD="bash"
      shift
      ;;
    -c|--cmd)
      CMD="$2"
      shift
      shift
      ;;
    *)
      args+=("$1")
      shift
      ;;
  esac
done

_main "${args[@]}"
cd $OWD


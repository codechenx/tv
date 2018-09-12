#!/usr/bin/env bash


# This script will install tv to the directory you're in. To install
# somewhere else (e.g. /usr/local/bin) just move tv binary into it


function githubLatestTag {
    finalUrl=`curl https://github.com/$1/releases/latest -s -L -I -o /dev/null -w '%{url_effective}'`
    echo "${finalUrl##*v}"
}

UNAME=$(uname)

platform=""

if [ "$UNAME" == "Linux" ] ; then
	platform="linux_amd64"
elif [ "$UNAME" == "Darwin" ] ; then
	 platform="darwin_amd64"
elif [[ "$UNAME" == CYGWIN* || "$UNAME" == MINGW* ]] ; then
	platform="windows_amd64.exe"
fi

echo "Detected platform: $platform"


TAG=`githubLatestTag codechenx/tv`

echo "Downloading https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$platform""
curl -L "https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$platform"" > tv
chmod +x tv

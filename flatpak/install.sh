#!/usr/bin/env bash


function githubLatestTag {
    finalUrl=`curl https://github.com/$1/releases/latest -s -L -I -o /dev/null -w '%{url_effective}'`
    echo "${finalUrl##*v}"
}

UNAME=$(uname)
ARCH=$(uname -m)
platform=""

if [ "$UNAME" == "Linux" ] ; then
	 platform="linux"
elif [ "$UNAME" == "Darwin" ] ; then
	 platform="darwin"
fi

if [ "$ARCH" == "x86_64" ] ; then
	architecture="amd64"
elif [ "$ARCH" == "armv7l" ] ; then
	 architecture="armv7"
elif [ "$ARCH" == "arm64" ] ; then
	 architecture="arm64"
else
	 architecture="386"
fi

echo "Detected platform: $platform"_"$architecture"


TAG=`githubLatestTag codechenx/tv`

echo "Downloading https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$TAG"_"$platform"_"$architecture".tar.gz"
curl -L "https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$TAG"_"$platform"_"$architecture".tar.gz" > tv_"$TAG"_"$platform"_"$architecture".tar.gz
tar -zxf tv_"$TAG"_"$platform"_"$architecture".tar.gz
rm -f tv_"$TAG"_"$platform"_"$architecture".tar.gz
chmod +x tv
rm -r LICENSE README.md
echo "#################################################################################
This script has downloaded tv binary file to current directory
you need to move tv binary file to any directory which is in the environment variable PATH"

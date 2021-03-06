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

FILE=$PWD/tv
if test -f "$FILE"; then
    echo "cant't download tv, because $FILE exists."
    exit 1
fi

TAG=`githubLatestTag codechenx/tv`
tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
echo "Downloading https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$TAG"_"$platform"_"$architecture".tar.gz"
curl -L "https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$TAG"_"$platform"_"$architecture".tar.gz" > $tmp_dir/tv_"$TAG"_"$platform"_"$architecture".tar.gz
tar -zxf $tmp_dir/tv_"$TAG"_"$platform"_"$architecture".tar.gz -C $tmp_dir
mv $tmp_dir/tv $PWD
chmod +x tv
echo "#################################################################################
This script has downloaded tv binary file to current directory
you need to move tv binary file to any directory which is in the environment variable PATH"
exit 0

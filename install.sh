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
curl -L "https://github.com/codechenx/tv/releases/download/v$TAG/tv_"$TAG"_"$platform"_"$architecture".tar.gz" > tv.tar.gz
tar -zxf tv.tar.gz
chmod +x tv

echo "\033[33mThis script will download tv binary file to your current directory
you need to run sudo cp tv /usr/local/bin/ 
or copy tv binary file to any directory which is in the environment variable PATH"

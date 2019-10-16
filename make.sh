#!/bin/sh
version=$1
if [[ -z "$version" ]]; then
  echo "usage: $0 <version> [-install] [-build]"
  echo "version must be in the format: major.minor.build"
  exit 1
fi

version=(${version//./ })
versionmajor=${version[0]}
versionminor=${version[1]}
versionbuild=${version[2]}

if [ -z "$versionmajor" ]
then
    echo "version error: major cannot be empty"
    exit 1
fi

if [ -z "$versionminor" ]
then
    echo "version error: minor cannot be empty"
    exit 1
fi

if [ -z "$versionbuild" ]
then
    echo "version error: build cannot be empty"
    exit 1
fi

version=$(<version/template)
version=${version/const Major int = 0/const Major int = $versionmajor}
version=${version/const Minor int = 0/const Minor int = $versionminor}
version=${version/const Build int = 0/const Build int = $versionbuild}
echo "$version" > version/version.go

echo "emits version set to ${versionmajor}.${versionminor}.${versionbuild}"

if [[ $* == *-build* ]]; then
    rm -rf bin/${versionmajor}.${versionminor}.${versionbuild}
    for platform in darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm64 windows/386 windows/amd64
    do
        platform=(${platform//\// })
        goos=${platform[0]}
        goarch=${platform[1]}
        file=emits
        if [ $goos = "windows" ]; then
            file+=".exe"
        fi
        echo "building emits version ${versionmajor}.${versionminor}.${versionbuild} ${goos}/${goarch}"
        GOOS=${goos} GOARCH=${goarch} go build -o bin/${versionmajor}.${versionminor}.${versionbuild}/${goos}/${goarch}/${file} bitbucket.org/emits-io/emits
        #if [ $? -eq 0 ]; then
        #fi
    done
fi

if [[ $* == *-install* ]]; then
    echo "installing emits version ${versionmajor}.${versionminor}.${versionbuild}"
    go build && go install
    emits version
fi

exit 0
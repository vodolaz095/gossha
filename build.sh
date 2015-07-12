#!/bin/sh

semver="1.0.2"
ver="$semver.`git log --format='%h' -n 1`.`uname`.`uname -m`"
subver="Build #`git log --format='%h' -n 1` on `hostname`.`uname`.`uname -m` on `date`"
archiv="build/gossha.`uname`.`uname -m`"

echo "Starting build $ver... on `date`"

echo "We have `go version`!"
echo "Installing dependencies..."
godep restore
echo "Dependencies installed!"

echo "Clearing distibution files..."
rm build/gossha -f
rm build/*.zip -f
rm build/*.tar.gz -f
rm build/*.tar.bz2 -f
echo "Distribution cleared!"

echo "Trying to engrave version..."
rm -f ver.go
echo "package gossha

const SUBVERSION = \"$subver\"
const VERSION = \"$ver\"">ver.go
echo "Version engraved!"

echo "Performing unit tests..."
if go test -v; then
	echo "Test complete!"

	echo "Building application..."
	if go build -o "build/gossha" app/gossha.go ; then
		echo "Build complete!"
		zip $archiv.$semver.zip build/gossha
		zip $archiv.$semver.zip README.md
		zip $archiv.$semver.zip README_RU.md
		zip $archiv.$semver.zip CHANGELOG.md
		zip $archiv.$semver.zip homedir/ -r
		zip $archiv.$semver.zip systemd/ -r
		echo "Creating .zip ok!"
		tar -czvf $archiv.$semver.tar.gz  build/gossha README.md README_RU.md CHANGELOG.md homedir/ systemd/
		echo "Creating .tar.gz ok!"
		tar -cjvf $archiv.$semver.tar.bz2 build/gossha README.md README_RU.md CHANGELOG.md homedir/ systemd/
		rm build/gossha -f
	else
		echo "Build failed!"
		exit 1
	fi
	git checkout ver.go
	echo "Build $ver complete on `date`!"
	exit 0
else
	echo "Unit tests for $ver/$subver failed!!!"
	exit 0
fi

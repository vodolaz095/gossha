export semver=1.0.3
export arch=$(shell uname).$(shell uname -m)
export gittip=$(shell git log --format='%h' -n 1)
export ver=$(semver).$(gittip).$(arch)
export subver=$(shell hostname).$(arch) on $(space)$(shell date)
export archiv=build/gossha.$(arch).$(semver)

clear:
	git checkout ver.go
	rm -f build/gossha
	rm -f build/*.zip
	rm -f build/*.tar.gz
	rm -f build/*.tar.bz2
	rm -f build/*.txt
	rm -f build/*.txt.sig

clean: clear

engrave:
	rm -f ver.go
	touch ver.go
	echo "package gossha" >> ver.go
	echo "" >>ver.go
	echo "//VERSION constant is engraved on build process" >>ver.go
	echo "const VERSION = \"$(ver)\"" >> ver.go
	echo "" >>ver.go
	echo "//SUBVERSION constant is engraved on build process" >>ver.go
	echo "const SUBVERSION = \"$(subver)\"" >>ver.go
	echo "" >>ver.go

deps:
	go get -u github.com/tools/godep
	godep restore

test: deps
	go get -u github.com/golang/lint/golint
	gofmt  -w=true -s=true -l=true ./..
	golint ./...
	go test -v


build: clear engrave deps test
	go build -o "build/gossha" app/gossha.go

pack: build
	zip $(archiv).zip  build/gossha README.md README_RU.md CHANGELOG.md homedir/ systemd/ -r
	tar -czvf $(archiv).tar.gz  build/gossha README.md README_RU.md CHANGELOG.md homedir/ systemd/
	tar -cjvf $(archiv).tar.bz2 build/gossha README.md README_RU.md CHANGELOG.md homedir/ systemd/

sign:
	rm build/*.txt -f
	rm build/*.txt.sig -f
	find build/ -name gossha.* -exec md5sum {} + > build/md5sum.txt
	gpg2 -a --output build/md5sum.txt.sig  --detach-sig build/md5sum.txt
	gpg2 --verify build/md5sum.txt.sig build/md5sum.txt
	find build/ -name gossha.* -exec sha1sum {} + > build/sha1sum.txt
	gpg2 -a --output build/sha1sum.txt.sig --detach-sig build/sha1sum.txt
	gpg2 --verify build/sha1sum.txt.sig build/sha1sum.txt
	@echo ""
	@echo ""
	@echo "MD5 hashes"
	@echo "========================"
	@cat build/md5sum.txt
	@echo ""
	@echo ""
	@echo "SHA1 hashes"
	@echo "========================"
	@cat build/sha1sum.txt
	@echo ""
	@echo ""
	@echo "*.sig files are signed with my GPG key of \`994C6375\`"


all: build

install: build
	su -c 'mv build/gossha /usr/bin/gossha'

uninstall:
	su -c 'rm /usr/bin/gossha -f'

remove: uninstall


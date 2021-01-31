export app=gossha
export semver=2.0.0
export arch=$(shell uname)-$(shell uname -m)
export gittip=$(shell git log --format='%h' -n 1)
export ver=$(semver).$(gittip).$(arch)
export archiv=build/$(app)-$(arch)-v$(semver)


all: build

deps:
	go mod download
	go mod verify
	go mod tidy

check: deps lint
# TODO

# Install git hooks
install_git_hooks:
	git config --global core.hooksPath .githooks/

lint:
	gofmt  -w=true -s=true -l=true ./
	golint ./..
	go vet

test: check

build: clean deps check
	go build -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/$(app)" main.go
	upx build/$(app)

build_podman:
	podman build -t reg.vodolaz095.life/gossha:$(ver) .

build_without_tests: clean deps
	go build -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/$(app)" main.go
	upx build/$(app)

dist: build
	zip $(archiv).zip  build/$(app) README.md README_RU.md CHANGELOG.md homedir/ contrib/ -r
	tar -czvf $(archiv).tar.gz  build/$(app) README.md README_RU.md CHANGELOG.md homedir/ contrib/
	tar -cjvf $(archiv).tar.bz2 build/$(app) README.md README_RU.md CHANGELOG.md homedir/ contrib/

start:
	go run main.go

sign: dist
	rm build/*.txt -f
	rm build/*.txt.sig -f
	find build/ -name gossha-* -exec md5sum {} + > build/md5sum.txt
	gpg2 -a --output build/md5sum.txt.sig  --detach-sig build/md5sum.txt
	gpg2 --verify build/md5sum.txt.sig build/md5sum.txt
	find build/ -name gossha-* -exec sha1sum {} + > build/sha1sum.txt
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

clean:
	rm -rf build/$(app)
	rm -rf build/*.zip
	rm -rf build/*.tar.gz
	rm -rf build/*.tar.bz2
	rm -rf build/*.txt
	rm -rf build/*.txt.sig
	rm -rf build/id_rsa
	rm -rf build/id_rsa.pub

test: check

install: build
	su -c 'cp -f build/$(app) /usr/bin/'

uninstall:
	su -c 'rm -rf /usr/bin/gossha'

keys:
	rm -f build/id_rsa
	rm -f build/id_rsa.pub
	ssh-keygen -N '' -f build/id_rsa

.PHONY: build

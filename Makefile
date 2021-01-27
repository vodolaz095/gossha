export semver=2.0.0
export arch=$(shell uname)-$(shell uname -m)
export gittip=$(shell git log --format='%h' -n 1)
export ver=$(semver).$(gittip).$(arch)
export archiv=build/gossha-$(arch)-v$(semver)


all: build

deps:
	go get -u github.com/golang/lint/golint

check:
	gofmt  -w=true -s=true -l=true ./
	golint ./..
	go vet

test: check

build: clean deps check
	go build -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

build_mysql_only: clean deps check
	go build -tags "mysql nosqlite3" -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

build_mysql_sqlite3: clean deps check
	go build -tags "mysql" -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

build_postgress_only: clean deps check
	go build -tags "postgresql nosqlite3" -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

build_mysql_postgress: clean deps check
	go build -tags "mysql postgresql nosqlite3" -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

build_postgress_sqlite3: clean deps check
	go build -tags "postgresql" -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

build_mysql_postgress_sqlite3: clean deps check
	go build -tags "postgresql mysql" -ldflags "-X github.com/vodolaz095/gossha/version.Version=$(ver)" -o "build/gossha" main.go

dist: build_mysql_postgress_sqlite3
	zip $(archiv).zip  build/gossha README.md README_RU.md CHANGELOG.md homedir/ contrib/ -r
	tar -czvf $(archiv).tar.gz  build/gossha README.md README_RU.md CHANGELOG.md homedir/ contrib/
	tar -cjvf $(archiv).tar.bz2 build/gossha README.md README_RU.md CHANGELOG.md homedir/ contrib/

docker: build_mysql_postgress_sqlite3 keys
	systemctl start docker
	docker build -t gossha .

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
	rm -rf build/gossha
	rm -rf build/*.zip
	rm -rf build/*.tar.gz
	rm -rf build/*.tar.bz2
	rm -rf build/*.txt
	rm -rf build/*.txt.sig
	rm -rf build/id_rsa
	rm -rf build/id_rsa.pub

test: check

install: build
	su -c 'cp -f build/gossha /usr/bin/'

uninstall:
	su -c 'rm -rf /usr/bin/gossha'

keys:
	rm -f build/id_rsa
	rm -f build/id_rsa.pub
	ssh-keygen -N '' -f build/id_rsa

.PHONY: build

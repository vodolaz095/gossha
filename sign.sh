#!/bin/sh

find build/ -name gossha.* -exec md5sum {} + > build/md5sum.txt
find build/ -name gossha.* -exec sha1sum {} + > build/sha1sum.txt

gpg2 -a --output build/md5sum.txt.sig --detach-sig build/md5sum.txt
gpg2 -a --output build/sha1sum.txt.sig --detach-sig build/md5sum.txt

#!/bin/sh

rm build/*.txt -f
rm build/*.txt.sig -f

find build/ -name gossha.* -exec md5sum {} + > build/md5sum.txt
gpg2 -a --output build/md5sum.txt.sig  --detach-sig build/md5sum.txt
gpg2 --verify build/md5sum.txt.sig build/md5sum.txt

sleep 1

find build/ -name gossha.* -exec sha1sum {} + > build/sha1sum.txt
gpg2 -a --output build/sha1sum.txt.sig --detach-sig build/sha1sum.txt
gpg2 --verify build/sha1sum.txt.sig build/sha1sum.txt

sleep 1
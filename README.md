GoSSHa
=================================

[![Build Status](https://travis-ci.org/vodolaz095/gossha.svg)](https://travis-ci.org/vodolaz095/gossha)
[![GoDoc](https://godoc.org/github.com/vodolaz095/gossha?status.svg)](https://godoc.org/github.com/vodolaz095/gossha)

[Руководство пользователя на Русском языке](https://github.com/vodolaz095/gossha/blob/master/README_RU.md)

Cross-platform ssh-server based chat program, with data persisted into
relational databases of MySQL, PostgreSQL or Sqlite3.
Public channel (with persisted messages) and private message (not stored) are supported.
Application has serious custom scripting and hacking potential.

Use case - devops chat with possibility to run scripts from chat, without SSH access to server.


Main addvantages
=================================

1) [Secure SHell protocol](https://en.wikipedia.org/wiki/Secure_Shell) is used to make all communications safe and secure.

2) Users' profiles and messages are stored in relational databases in easy to manipulate format, so we can use 3rd party applications to work with them.

3) Users can be authorized by passwords or private keys.

4) We can start application listening on few addresses and ports on the same time. For example, listeing on `192.168.1.2:2222` on local area network, and on `193.41.32.25:27015` for uplink connections.

5) Users can execute scripts defined by admin on behalf of users running the GoSSHa server.

6) Application can run scripts after each public or private message, with senders name, ip, message exported as environment variables. See `homedir/` folder for examples

7) Application is created in `Go` language, and can be build on many environments and architectures - `Linux`, `Microsoft Windows`, `MacOs`.

Usage
=================================

Firstly, you can create admin account by calling `$ gossha root [username]`

```shell

	[vodolaz095@rhel ~]$ gossha root admin
	  ____      ____ ____  _   _
	 / ___| ___/ ___/ ___|| | | | __ _
	| |  _ / _ \___ \___ \| |_| |/ _` |
	| |_| | (_) |__) |__) |  _  | (_| |
	 \____|\___/____/____/|_| |_|\__,_|


	Persistent SSH based chat for the ones, who cares.
	Build: 1.24.1.b06789e.Linux.x86_64
	Version: Build #b06789e on rhel.Linux.x86_64 on Sun Jun 28 01:10:39 MSK 2015

	Console commands avaible:
	 $ gossha ban [username] - delete user and all his/her messages
	 $ gossha passwd [username] - create/update ordinary user by name and password
	 $ gossha root [username] - create/update root user by name and password

	Empty argument - start in server mode

	Enter password:
	User admin is created and/or new password is set!

```

Than you can login using any of [SSH clients](https://en.wikipedia.org/wiki/Comparison_of_SSH_clients)

For example, like this

```shell

		$ ssh admin@localhost -p 27015

```

Than you can import you private ssh key to be used instead of password by using
the `\k` command.

```

	[vodolaz095@rhel ~]$ ssh admin@localhost -p 27015
	Host key fingerprint is 3d:63:45:c4:82:03:ca:99:80:49:03:8e:f2:d8:3a:bb
	+--[ RSA 2048]----+
	|+=.   .. . oo    |
	|= .o +  o ...    |
	|o.  =    . ..    |
	|.+       . .     |
	|. o     S =      |
	| .       . o     |
	|o                |
	| o               |
	|E.               |
	+-----------------+

	admin@localhost's password:
	GoSSHa - very secure chat.
	Build #1.24.1.b06789e.Linux.x86_64
	Version: Build #b06789e on rhel.Linux.x86_64 on Sun Jun 28 01:10:39 MSK 2015
	Commands avaible:
	 \b - (B)an user (you need to have `root` permissions!)
	 \e - Close current session
	 \exit - Close current session
	 \f - (F)orgot localy available SSH key used for authorising your logins via this client
	 \h - (H)elp, show this screen
	 \i - Print (I)nformation about yourself
	 \k - Use locally available SSH (K)eys to authorise your logins on this server
	 \passwd - Changes current user password
	 \q - Close current session
	 \quit - Close current session
	 \r - (R)egister new user (you need to have `root` permissions!)
	 \rr - (R)egister new (r)oot user (you need to have `root` permissions!)
	 \w - List users, (W)ho are active on this server
	 \x - E(X)ecutes custom user script from home directory
	 all other input is treated as message, that you send to server


	[admin@localhost.localdomain(127.0.0.1) x]{14:14:56}:hello!!!
	[admin@localhost.localdomain(127.0.0.1) *]{02:24:04}:\k
	Importing public key...
	Key imported succesefully!
	[admin@localhost.localdomain(127.0.0.1) *]{02:24:04}:

```
Ordinary messages are colored in `white`, system messages - `green`, private
messages - `blue`.
To send private message, type `@`, than username (`TAB` autocompletion works) to
whom you want to send private message of the record. Private messages are not
stored in the database, and they disapear, when user logouts and logins.


Configuration parameters
=================================
Application can be configured in few wayes (ordered by priority).

1) By starting application with flags defined.

2) By environment variables

3) By JSON object values in config file loaded from `/etc/gossha/gossha.json`

4) By JSON object values in config file loaded from `$HOME/.gossha/gossha.json`

This is example config file provided with application:

```javascript

		{
		  "port": 27015,
		  "debug": false,
		  "driver": "sqlite3",
		  "connectionString": "/home/vodolaz095/.gossha/gossha.db",
		  "sshPublicKeyPath": "/home/vodolaz095/.ssh/id_rsa.pub",
		  "sshPrivateKeyPath": "/home/vodolaz095/.ssh/id_rsa",
		  "homedir": "/home/vodolaz095/.gossha",
		  "executeOnMessage": "",
		  "executeOnPrivateMessage": ""
		}

```

***Port*** (integer) for application to listein on `0.0.0.0` address (all interfaces). The
default value is `27015`, it can be set by `--port=27015` flag, or via `GOSSHA_PORT=27015`
environment value.


***Debug*** (boolean) toggle mode with usage of more verbose output to stdout and start [pprof](https://golang.org/pkg/net/http/pprof/)
server on `localhost:3000` port for debugging/benchmarking purposes.
Can be enabled by `--debug=true` flag, or via `GOSSHA_DEBUG=true` environment value.

***Driver*** and ***connectionString*** sets the connection to database.
We can use [sqlite3](https://github.com/mattn/go-sqlite3),
[MySQL](https://github.com/go-sql-driver/mysql) (`MariaDB` in compatibility mode),
[PostgreSQL](https://github.com/lib/pq) databases via appropriate drivers.

Possible pairs of values are

```
   	--driver=sqlite3 --connectionString=/var/lib/gossha/gossha.db

   	--driver=sqlite3 --connectionString=:memory:

   	--driver=mysql --connectionString='user:password@/dbname?charset=utf8&parseTime=True&loc=Local'

   	--driver=postgres --connectionString='user=gorm dbname=gorm sslmode=disable'

   	--driver=postgres --connectionString='postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full'

```
by default, the `sqlite3` driver is used with database stored at `$HOME/.gossha/gossha.db`.
We can load `driver` from `--driver=sqlite3` flag or `GOSSHA_DRIVER=sqlite3` environment value
We can load `connectionString` from `--connectionString=:memory:` flag or `GOSSHA_CONNECTIONSTRING=:memory:` environment value

***SshPublicKeyPath*** points to Public Key to be used by SSH server, default value is `$HOME/.ssh/id_rsa.pub`.
Can be set via `--sshPublicKeyPath=/home/myusername/.ssh/id_rsa.pub` flag or `GOSSHA_SSHPUBLICKEYPATH=/home/myusername/.ssh/id_rsa.pub` environment value.

***sshPrivateKeyPath*** points to Private Key to be used by SSH server, default value is `$HOME/.ssh/id_rsa.pub`.
Can be set via `--sshPrivateKeyPath=/home/myusername/.ssh/id_rsa` flag or `GOSSHA_SSHPRIVATEKEYPATH=/home/myusername/.ssh/id_rsa` environment value.

***Homedir*** is path to directory containing user's scripts to be executed via `\x` command in chat. It is worth notice,
that this scripts have to be executable files, like the examples, provided in `homedir/scripts` directory of
the distribution or repo. The username, ip and other data is populated from environment values used for scripts.
We can make this executable files in any language - shell, binaries, nodejs files, php scripts.
Can be set via `--homedir=/home/myusername/.gossha` flag or `GOSSHA_HOMEDIR=/home/myusername/.gossha` environment value.

***executeOnMessage*** is path to executable to be started on each message.
We can make this executable files in any language - shell, binaries, nodejs files, php scripts.
See `homedir/afterMessage` for shell example.
Can be set via `--executeOnMessage=/home/myusername/.gossha/afterMessage` flag or `GOSSHA_EXECUTEAFTERMESSAGE=/home/myusername/.gossha/afterMessage` environment value.

***executeOnPrivateMessage*** is path to executable to be started on each message.
We can make this executable files in any language - shell, binaries, nodejs files, php scripts.
See `homedir/afterPrivateMessage` for shell example.
Can be set via `--executeOnPrivateMessage=/home/myusername/.gossha/afterPrivateMessage` flag or `GOSSHA_EXECUTEAFTERPRIVATEMESSAGE=/home/myusername/.gossha/afterPrivateMessage` environment value.


Building from sources
=================================
I assume you have one of popular `Linux` distros, i don't care about other OSes.

1) [Install Go language](http://golang.org/doc/install) and it's [environment](http://golang.org/doc/code.html#GOPATH) properly. At least `1.4.2` version.

2) Verify you have [GNU Make](https://www.gnu.org/software/make/) at least of
4.0 version.

3) Clone code from repository in appropriate place

```shell

	$ cd $GOPATH/src/github.com
	$ mkdir vodolaz095
	$ cd vodolaz095
	$ git clone ssh://git@github.org/vodolaz095/gossha.git

```

3) Try to build 

```shell

	$ make

```
The binary file will be created in `build/gossha`

4) Try to install globaly (root password will be asked!) 

```shell

	$ make install

```
This step results in binary generated and placed in `/usr/bin/gossha`.
Also you can uninstall binaries by (root password will be asked!)

```shell

	$ make uninstall

```

5) By default, when you run the application first time, the directory 
with databases, configs and scripts will be created in `$HOME/.gossha/`



Installation via prebuild binaries
=================================

You can get compiled binaries from here
[https://github.com/vodolaz095/gossha/releases](https://github.com/vodolaz095/gossha/releases)

You can verify the signatures via `GPG` or `GPG2`. It have to be something like this:

```shell

		[vodolaz095@vodolaz095 build]$ gpg2 --verify md5sum.txt.sig md5sum.txt
		gpg: Signature made Mon 29 Jun 2015 02:44:13 AM MSK using RSA key ID 994C6375
		gpg: Good signature from "Anatoliy Ostroumov <ostroumov095@gmail.com>" [ultimate]
		gpg:                 aka "[jpeg image of size 2756]" [ultimate]
		gpg:                 aka "[jpeg image of size 3725]" [ultimate]


```

with RSA key ID of `994C6375`!


License
=================================
The MIT License (MIT)

Copyright (c) 2015 Ostroumov Anatolij ostroumov095(at)gmail(dot)com et al.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

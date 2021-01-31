# v 2.1.0
Improve password hashing. Introduce golang modules.

- Introduce golang [modules](https://golang.org/ref/mod) for dependency managment.
- Passwords are now hashed using [Argon2](https://github.com/alexedwards/argon2id) algorithm.
- Code can be compiled using Golang 1.14.13 now on centos8 machine.


# v 2.0.0
Newer dependencies. Code strongly refactored.
Breaking changes:

- configuration is stored in `.toml` file with tons of comments
- code simplification
- more unit tests
- possibility to build application with different database drivers being used. For example, build with only MySQL support, or build with MySQL, SQLite3 and PostgreSQL support
- deprecated overriding configuration by flags - it made online documentation very unfriendly, and using file/environment is easier.

Many typos and small issues fixed



# v 1.1.6
Newer dependencies. Build with Go 1.5.3

# v 1.1.5
Ability to set password for users from shell by calling

```shell

	$ gossha passwd user password
	$ gossha root user password

```

Many typos fixed. Dockerfile provided.

# v 1.1.4
Use other library to conceal password input for creating users
Newer dependencies. Small fixes.

# v 1.1.3
More recent dependencies. Build with Go 1.5.1.

# v 1.1.2
More recent dependencies, build with go1.4.2 linux/amd64

# v 1.1.1
More recent dependencies

# v 1.1.0
Tons of smallfixes, console interface is refactored. Console commands of `gossha list`,`gossha log`,`gossha dumpcfg` are added.
More verbose error reporter with link to bug-tracker.

# v 1.0.4
Newer crypto and gorm libs. More standart and automated `Makefile` behaviour. Readme updated.

# v 1.0.3
Newer crypto and gorm libs

# v 1.0.2
Newer crypto libs

# v 1.0.1
Code style, fixes in continious integration scripts.
Proper build environment for Raspberry Pi v1.
Removed some debug comments

# v 1.0.0
First release candidate

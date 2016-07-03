GoSSHa
=================================

[![Build Status](https://travis-ci.org/vodolaz095/gossha.svg)](https://travis-ci.org/vodolaz095/gossha)
[![GoDoc](https://godoc.org/github.com/vodolaz095/gossha?status.svg)](https://godoc.org/github.com/vodolaz095/gossha)


Кросплатформенный чат на основе SSH протокола.
Данные хранятся в реляционных базах данных MySQL, PostgreSQL or Sqlite3.
Поддерживаются публичные (хранятся в базе данных) и личные  сообщения.
Программа обладает большим потенциалом по использованию пользовательских скриптов
и автоматизации.

Пример использования - чат для системных администраторов на сервере, с возможносью
запускать скрипты по командам пользователей и интеграция с внешними базами данных.

Основные особенности
=================================

1) Используется надежный и защищённый протокол [Secure SHell protocol](https://ru.wikipedia.org/wiki/SSH)
для передачи всей информации.

2) Профили пользователей и сообщения хранятся в релационных базах данный.

3) Пользователь может авторизироваться с помощью пароля или приватного ключа SSH

4) Один экземпляр приложения может быть запущен сразу на нескольких портах и интерфейсах.
Например, мы может запустить приложение на `192.168.1.2:2222` в локальной сети и на `193.41.32.25:27015` для доступа из интернета.

5) Администратор сервера может создать скрипты, которые пользователи могут запускать на сервере от имени процесса программы.

6) Программа может запускать скрипты после каждого публичного или личного сообщения. См. примеры в директории `homedir/`

7) Программа написана на языке `Go`, и её можно скомпилировать не только на `Linux`, но также и на `Microsoft Windows`, `MacOs`.

Использование
=================================

Сначала нам надо создать аккаунт администратора - `$ gossha root [username]`

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

И теперь мы можем использовать какой-либо из [SSH clients](https://en.wikipedia.org/wiki/Comparison_of_SSH_clients)
для доступа к серверу, например:

```shell

		$ ssh admin@localhost -p 27015

```

После авторизации по имени и паролю, мы можем импортировать личный ключ с помощью команды `\k` в чате:

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
	Commands available:
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

Обычные сообщения выделены белым цветом, системные - зелёным, личные - синим.
Для отправки личных сообщений, наберите символ крякозабла `@`, потом имя пользователя
(автодополнение по  нажатию клавиши `TAB` работает), потом сообщение.
Личные сообщения `НЕ СОХРАНЯЮТСЯ` в базе данных, они могут быть доставлены только когда
получатель соединён с сервером, и приватные сообщения пропадают, когда пользователь
завершает сеанс.


Настройка сервера
=================================
Установить параметры программы можно с помощью приведённых ниже методов, отсортированных
по важности (более высокий метод имеет приемущество)


1) Запустив программу с флагами командной строки

2) Установив переменные окружения

3) Ввести параметры в конфигурационный файл, загруженный из `/etc/gossha/gossha.toml`

4) Ввести параметры в конфигурационный файл, загруженный из `$HOME/.gossha/gossha.toml`

Вот пример конфигурации:

```toml

# Automatically generated config file for GoSSHa - SSH powered chat
# Place it either in
#   /etc/gossha/gossha.toml
# or
#   ~/.gossha/gossha.toml
#

# Enable debug
Debug=true

# On what port to listen for all interfaces (like for 0.0.0.0 address)
Port = 27015

# What addresses to bind to
BindTo = ["127.0.0.1:27014"]

#Setting database connection - various possible combinations are shown

#SQLite3 with database in local file
#Driver = "sqlite3"
#ConnectionString = "/home/vodolaz095/.gossha/gossha.db"

#SQLite3 with database in memory
#Driver = "sqlite3"
#ConnectionString = ":memory:"

#MySQL database
#Driver = "mysql"
#ConnectionString = "username:password@hostname/database?charset=utf8&parseTime=True&loc=Local"

#PostgreSQL database. 1st variant
#Driver = "postgres"
#ConnectionString ="user=gorm dbname=gorm sslmode=disable"

#PostgreSQL database. 2nd variant
#Driver="postgres"
#ConnectionString="postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")"

#This database connection setting are being used
Driver = "sqlite3"
ConnectionString = "/home/vodolaz095/.gossha/gossha.db"


#Path to SSH Public key
SshPublicKeyPath = "/home/vodolaz095/.ssh/id_rsa.pub"
#Path to SSH Private key
SshPrivateKeyPath = "/home/vodolaz095/.ssh/id_rsa"

#Directory to search for custom scripts
Homedir = "/home/vodolaz095/.gossha"

#Script to be executed on each message
ExecuteOnMessage=""

#Script to be execute on each private message
ExecuteOnPrivateMessage=""


```

***Port*** (целое положительное число) порт, который слушает приложение на  `0.0.0.0` адресе (Все интерфейсы).
Значение по умолчанию - `27015`. Этот параметр можно задать установив переенную окружения `GOSSHA_PORT=27015`ы.


***Debug*** (булевое значение) включить режим отладки, совместивый с [pprof](https://golang.org/pkg/net/http/pprof/).
Получить доступ к отладочной информации можно по адресу [http://localhost:6060](http://localhost:6060).
Также данный режим можно включить, используя переменную окружения `GOSSHA_DEBUG=true`

***Driver*** и ***connectionString*** позволяют настроись соединение с базой данных
с помощью соответствующих драйверов:
[sqlite3](https://github.com/mattn/go-sqlite3),
[MySQL](https://github.com/go-sql-driver/mysql) (`MariaDB` в режиме совместимости),
[PostgreSQL](https://github.com/lib/pq).

Возможные комбинации параметров -

```
   	--driver=sqlite3 --connectionString=/var/lib/gossha/gossha.db

   	--driver=sqlite3 --connectionString=:memory:

   	--driver=mysql --connectionString='user:password@/dbname?charset=utf8&parseTime=True&loc=Local'

   	--driver=postgres --connectionString='user=gorm dbname=gorm sslmode=disable'

   	--driver=postgres --connectionString='postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full'

```
По умолчанию используется драйвер базы данных `sqlite3` с базой, хранящейся в `$HOME/.gossha/gossha.db`.
Параметр `driver` можно задать переменной окружения `GOSSHA_DRIVER=sqlite3`
Параметр  `connectionString` можно задать переменной окружения `GOSSHA_CONNECTIONSTRING=:memory:`.

***SshPublicKeyPath*** указывает на путь к файлу публичного ключа, используемого SSH сервером.
Значение по умолчанию `$HOME/.ssh/id_rsa.pub`.
Параметр можно задать переменной окружения `GOSSHA_SSHPUBLICKEYPATH=/home/myusername/.ssh/id_rsa.pub`.

***sshPrivateKeyPath***  указывает на путь к файлу личного ключа, используемого SSH сервером.
Значение по умолчанию `$HOME/.ssh/id_rsa.pub`.
Параметр можно задать переменной окружения  `GOSSHA_SSHPRIVATEKEYPATH=/home/myusername/.ssh/id_rsa`.

***Homedir*** путь к директории, содержащей исполняемые файлы, которые можно запустить используя команду
 `\x` в чате. Эти исполняемые файлы могут быть бинарными файлами, шелл скриптами, исполняемыми файлами,
как примеры в директории  `homedir/scripts`. Имя пользователя, IP адресс, и другие данные
устанавливаются как переменные окружения сервера.
Параметр можно задать переменной окружения `GOSSHA_HOMEDIR=/home/myusername/.gossha`.

***executeOnMessage*** путь к исполняемому файлу, запускаемому после доставки каждого публичного
сообщения. Исполняемый файл может быть шелл скриптом, бинарным, nodejs или PHP скриптом.
См. пример `homedir/afterMessage`.
Параметр можно задать переменной окружения `GOSSHA_EXECUTEAFTERMESSAGE=/home/myusername/.gossha/afterMessage`.

***executeOnPrivateMessage*** путь к исполняемому файлу, запускаемому после доставки каждого личного
сообщения. Исполняемый файл может быть шелл скриптом, бинарным, nodejs или PHP скриптом.
См. пример `homedir/afterPrivateMessage`.
Параметр можно задать переменной окружения `GOSSHA_EXECUTEAFTERPRIVATEMESSAGE=/home/myusername/.gossha/afterPrivateMessage` environment value.


Компилирование из исходных кодов
=================================
Предполагается, что используется какой-либо из популярных дистрибутивов `Linux`

1) [Установите язык программирования Go](http://golang.org/doc/install) и его [среду исполнения](http://golang.org/doc/code.html#GOPATH). Версия языка программирования Go должна быть новее `1.4.2`.

2) Убедитесь, что на Вашем компьютере установлен [GNU Make](https://www.gnu.org/software/make/) как минимум 4.0 версии.

3) Загрузите исходные коды программы из репозитория в соответствующую директорию в $GOPATH

```shell

	$ go get github.com/vodolaz095/gossha

```

3) Помолившись, попробуйте скомпилировать приложение.

```shell

	$ cd $GOPATH/src/github.com/vodolaz095/gossha
	$ make

```

и бинарный файл должен появится в `build/gossha`.


4) Чтобы установить программу глобально, запустите (будет запрошен пароль суперпользователя):

```shell

	$ make install

```
В результате в директории `/usr/bin/gossha` будет создан исполняемый файл программы.
Также можно удалить программу из системы, воспользовавшись командой  (будет запрошен пароль суперпользователя):

```shell

	$ make uninstall

```

5) По умолчанию, после запуска программы директория с базой данных и конфигурационными файлами будет создана в 
домашней директории пользователя в `$HOME/.gossha/`.


Установка из скомпилированных бинарных файлов
=================================

Вы можете скачать бинарные файлы с этого адреса 
[https://github.com/vodolaz095/gossha/releases](https://github.com/vodolaz095/gossha/releases)

При подтверждении подписи вывод должен быть примерно такой:

```shell

		[vodolaz095@vodolaz095 build]$ gpg2 --verify md5sum.txt.sig md5sum.txt
		gpg: Signature made Mon 29 Jun 2015 02:44:13 AM MSK using RSA key ID 994C6375
		gpg: Good signature from "Anatoliy Ostroumov <ostroumov095@gmail.com>" [ultimate]
		gpg:                 aka "[jpeg image of size 2756]" [ultimate]
		gpg:                 aka "[jpeg image of size 3725]" [ultimate]


```

Я использую ключ №`994C6375`!


License
=================================
The MIT License (MIT)

Copyright (c) 2015 Ostroumov Anatolij ostroumov095(at)gmail(dot)com et al.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

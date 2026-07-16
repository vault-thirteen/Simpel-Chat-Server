# Simpel Chat Server

This is a simple chat server written in pure Go programming language. It tries 
to use only simple technologies and algorithms. The main purpose of this 
product is to allow freedom of speech for those who need it, as all the 
conscious living beings in the Universe have the right to speak freely in their 
life. 

The name `Simpel` is not a mistake. This name is intentional. `Simpel` is the 
word `Simple` in its normal, unspoiled form.

This product is a back-end part of the chat, i.e. a chat server.

## Functionality

The chat consists of chat rooms and users which may use these rooms. The chat 
provides users with access to chat rooms in which they may communicate using 
textual messages. The main idea of the chat is following. Each user may use 
only a single room at the same time. This principle copies every one's real 
life, when each physical person can be located only in one physical room at 
one moment of time.

Users may have three different access levels: 
- simple users;
- moderators;
- administrators. 

Simple users are the majority of users, who enter or leave rooms and use chat 
to communicate with other users. In order to read messages from a room or write 
message into a room a user must enter a particular room. As it was mentioned 
earlier, a user may join only a single room at the same time, i.e. you can 
not be in two different places at the same time. This means that before joining 
another room, you must leave your current room, as in real life, as in simple 
web chats of past time and as it is done in TeamSpeak communication system. 

Further description requires you to know that chat rooms can be either public 
or private. Public rooms are open for any user. Private rooms allow only 
selected users to join. 

Moderators are those users who can change rules for rooms. Each room can have 
its own separate set of moderators, independent of other rooms. 

Administrators are those who can create chat rooms and appoint moderators to 
each private room.

## Technologies

The server uses simple, and sometime even primitive, algorithms, methods and 
technologies. Text files are used for application settings, JSON files contain 
Chat's configuration. Room settings, users and sessions are stored in a simple 
MySQL database. Messages are stored in memory (RAM) for fast access. 

In order to limit memory usage, message size and message count in each chat 
room are limited. Room count and sessions number are also limited for the same 
reason. Each parameter is configurable for those who have a machine powerful 
enough to hold everything in memory.

Communication between front-ends, or clients, and the server is done using the 
network protocol of the HTTP family (HTTPS, TLS) with the help of the original 
RPC protocol, based on JSON textual format. In the server's core, an ORM system 
named GORM helps the system to communicate with MySQL database.

Authorisation is super simple and is performed using two factors – user's 
password and a verification code, sent via a common electronic mail or simply 
e-mail.

Note that while all the messages are stored in RAM, they are not persistent 
across server's restarts. This feature makes communication very fast and does 
not make you waste free space on your storage drives, however the downside of 
it does not allow to save messages to persistent storage. This is done on 
purpose. 

## Installation

1. Prepare your SSL certificates. 
 
A script to create self-signed certificates is available in the `script` 
folder.

2. Build the project using the `build.bat` script.

3. Copy files and folders from the created `_BUILD_` directory to your place.

4. Get the executable binary file using the following command.

> go install github.com/vault-thirteen/Simpel-Chat-Server/src@latest

5. Replace the executable file created by the build script with the file 
received with `go install` command.

6. Say "Thank you" to the developers of Go language for not fixing old bugs in 
`go install` tool, for old bugs with versioning and many other old bugs in Go 
language.

**Important note**. Do not try to build the executable file locally. You will 
see an old versioning bug and will be unable to use the server normally.

## Usage

To start a server, compile the application and provide it with a path to the 
application configuration file. A build script `build.bat` can help you.

> server.exe /path/to/app.cfg

An example of an application configuration file can be found in the `settings` 
folder, the file is named `app.cfg`. This file uses simple format – each line 
is a separate setting.

1. 1-st line: A path to chat settings;
2. 2-nd line: A boolean flag to load system DLL files for Windows operating system;
3. 3-rd line: A boolean flag to enable colours in a console.

Both of the flags may be used only on Windows O.S, i.e. Linux users should 
have them as 'false'. 

An example of chat settings, written in JSON file, is available in the same 
folder, the file is named `chat.json`. Default settings are quite good for a 
simple small chat server.

Configuration has one special feature which increases safety: if database 
password or SMTP server password is not set in the JSON configuration file, 
the server asks for them in console during the start-up procedure.

A reminder for those who read text diagonally. This is a back-end part of the 
product. In order to use this chat as a client, a separate front-end product 
is required. Also, MySQL server is not included into the product, you must 
start it separately.

## Reason

This project was started in memory of old-school web chats which were
popular in the early 2000-ish years. All in all, this product is very
simplistic, yet quite powerful for its purpose.

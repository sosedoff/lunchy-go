# lunchy-go

A friendly wrapper for launchctl. Start your agents and go to lunch!

This is a port of original [lunchy](https://github.com/mperham/lunchy) ruby gem by Mike Perham with extra functionality.

[![Release](https://img.shields.io/github/release/sosedoff/lunchy-go.svg)](https://github.com/sosedoff/lunchy-go/releases)

## Overview

Don't you hate OSX's launchctl? You have to give it exact filenames. 
The syntax is annoying different from Linux's nice, simple init system and overly verbose. 
It's just not a very developer-friendly tool.

Lunchy aims to be that friendly tool by wrapping launchctl and providing a few 
simple operations that you perform all the time:

- ls [pattern]
- start [pattern]
- stop [pattern]
- restart [pattern]
- status, ps [pattern]
- install [file]
- show [pattern]
- edit [pattern]
- remove, rm [pattern]
- scan [path]

where pattern is just a substring that matches the agent's plist filename. 

So instead of:

```
launchctl load ~/Library/LaunchAgents/io.redis.redis-server.plist
```

you can do this:

```
lunchy start redis
```

and:

```
$ lunchy ls

com.danga.memcached
com.google.keystone.agent
com.mysql.mysqld
io.redis.redis-server
org.mongodb.mongod
```

## Install

You can install binary by running the following bash command:

```
curl -s https://raw.githubusercontent.com/sosedoff/lunchy-go/master/install.sh | bash
```

### Homebrew

Formula is available as cask. Make sure you have [homebrew-cask](https://github.com/caskroom/homebrew-cask) installed.

```
brew cask install lunchy
```

### Binary Releases

Precompiled binaries are available on Github: https://github.com/sosedoff/lunchy-go/releases

### Build from source

Or you can build it youself if you have Go 1.2+ installed:

```
git clone https://github.com/sosedoff/lunchy-go.git $GOPATH/src/lunchy
cd lunchy
go build
mv ./lunchy-go /usr/local/bin/lunchy
```

## Usage

Add a new plist:

```
# Install plist
lunchy install /usr/local/Cellar/redis/2.8.1/homebrew.mxcl.redis.plist
```

Manage services:

```
lunchy start redis
lunchy stop redis
lunchy restart redis
lunchy status redis
```

If you have multiple plists from homebrew, you can simple control all of them:

```
$ lunchy status
homebrew.mxcl.elasticsearch
homebrew.mxcl.mysql
homebrew.mxcl.postgresql
homebrew.mxcl.redis

# Will stop all processes prefixed by "homebrew"
$ lunchy stop homebrew
```

Manage plists:

```
lunchy show redis
lunchy edit redis
```

Scan directory for existing plists:

```
lunchy scan /usr/local/Cellar
```

## License

The MIT License (MIT)

Copyright (c) 2013-2015 Dan Sosedoff, <dan.sosedoff@gmail.com>
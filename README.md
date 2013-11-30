# lunchy-go

A friendly wrapper for launchctl. Start your agents and go to lunch!

This is a port of original [lunchy](https://github.com/mperham/lunchy) ruby gem by Mike Perham.

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
- status [pattern]
- install [file]
- show [pattern]
- edit [pattern]

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
> lunchy ls
com.danga.memcached
com.google.keystone.agent
com.mysql.mysqld
io.redis.redis-server
org.mongodb.mongod
```

## Install

Download a prebuilt binary for OSX from Bintray: https://bintray.com/sosedoff/generic/lunchy

## Usage

Add a new plist:

```
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

$ lunchy stop homebrew
```

Manage plists:

```
lunchy show redis
lunchy edit redis
```

## License

The MIT License (MIT)

Copyright (c) 2013 Dan Sosedoff <dan.sosedoff@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
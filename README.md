# lunchy-go

A friendly wrapper for launchctl. Start your agents and go to lunch!

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
If you don't use a unique pattern, Lunchy will warn you of this and give you a 
list of the matching items instead.

So instead of:

    launchctl load ~/Library/LaunchAgents/io.redis.redis-server.plist

you can do this:

    lunchy start redis

and:

    > lunchy ls
    com.danga.memcached
    com.google.keystone.agent
    com.mysql.mysqld
    io.redis.redis-server
    org.mongodb.mongod

This is a port of original [lunchy](https://github.com/mperham/lunchy) gem by
Mike Perham.
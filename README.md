slack-tldr [![Build Status](https://travis-ci.org/ivey/slack-tldr.svg?branch=master)](https://travis-ci.org/ivey/slack-tldr)
==========

A Slack bot to provide easy access to pinned items inline in the chat.
Groups used to the hangupsbot built-in ```/bot tldr``` may find this
an useful tool to ease transition to Slack.

Installation
------------

1. Download the installation archive from [the releases tab](https://github.com/ivey/slack-tldr/releases).
2. Untar: ```tar xvfz slack-tldr.linux.amd64.tgz```
3. Create a Slack bot and copy the token.
4. Run it: ```./slack-tldr -token 'YOURSLACKTOKEN'```
5. Optional: use slack-tldr.service to keep it running under systemd


Usage
-----

Invite the bot to the channels you want it to be active in. Since it
stores items with pins in the channel, each channel has a separate
list.

If you just want to see the current pins, ```+tldr``` will trigger the
bot to display them.

You can create a new pin with ```+tldr Important thing to remember```.

You can delete a pinned item by getting the list from ```+tldr``` and
then ```+tldr remove 2``` where 2 is the number you want to delete.


Customization
-------------

If you'd rather use something other than ```+tldr``` as the command
string, the ```-command``` flag will allow you to change it.


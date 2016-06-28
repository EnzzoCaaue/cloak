# Cloaka

Cloaka is an Automatic Account Creator for Open Tibia. 

Cloaka provides great performance and an easy to use scripting environment using **lua**. Faster and less demanding than the PHP AACs out there.

It is made using [Pigo](https://github.com/Raggaer/Pigo) web-framework

# Commands

Cloaka support console commnads (handy for development) such as

- reload template (will compile all the templates.html)
- reload monsters (will load all the monsters from monsters.xml file)
- reload config lua (will reload the config.lua file into memory)
- reload config (will reload the config.json file into memory)
- reload stages (will reload the stages.xml file into memory)
- reload items (will reload the items.xml file into memory)
- exit (will terminate cloaka)

Feel free to create your own commands

# Compatible OT distributions

Cloaka is only compatible with [The Forgotten Server](https://github.com/otland/forgottenserver). Cloaka only supports its newer version 1.0 or greater.

# Download

Currently there are no releases of **Cloak** you can however build it

# How to build

To build your custom version of Cloaka you will need to have the following dependencies:

+ Go 1.5 or greater
+ Git

Once everything is installed you can simply do (to get the newer version)

```
go get github.com/Cloakaac/cloak
```

And finally

```
go build
```
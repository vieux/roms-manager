![roms-manager logo](./logo.jpg "roms-manager's logo")

roms-manager is a simple tool to scan arcade romset .dat and gamelist.xml files.

## Why?

`clrmamepro` is complicated to use, I want something simple, I do not care about rebuilding romsets, I simply want the roms I already have to work on my system by simply hidding the incompatible ones.

The main usecase is with RetroPie/Recalbox when using `mame`/`fbneo`/`neogeo` (and optionanly enabling the `arcade` virtual system), roms-manager is used to hide incompatible games (either bad roms or different screen aspect-ratio / buttons layout) and duplicates games so you end up with only a clean list of games.

## Install

### Latest release

Go to the [release page](https://github.com/vrgl117-games/roms-manager/releases) and grab the binary.

### Development build

Install golang on your machine and then `go get github.com/vrgl117-games/roms-manager`

## How to use

roms-manager has two mains functions

### Scan

`scan` takes a database file (either `.dat` or a Mame `.xml`) and an
EmulationStation `gameslist.xml` file and an optional `catver.ini` file.

Features: 
* hide incompatible games (wrong rom size, wrong CRC)
* Only keep on games amongst an original and it's a clone(s)
* hide games using a list of keywords (bootlegs, hacks, etc...)
* if present in the database file, hide games with the incorrect aspect ratio or button layout
* support catver.ini file to hide games based on category (mature by default)

See `roms-manager scan --help` for the list of flags.

### Hide-Duplicates

`hide-duplicates` takes a list of  `gameslist.xml` files.

Feature:
* hide games already present in other `gameslist.xml` file.

## TODO

* download `.dat` files on the fly
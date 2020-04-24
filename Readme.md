# roms-manager

roms-manager is a simple tool to scan arcade romset .dat and gamelist.xml files.

## Why?

`clrmamepro` is complicated to use, I want something simple, I do not care about rebuilding romsets, I simply want the roms I already have to work on my system by simply hidding the incompatible ones.

The main usecase is with RetroPie/Recalbox when using `mame`/`fbneo`/`neogeo` (and optionanly enabling the `arcade` virtual system), roms-manager is used to hide incompatible (either bad roms or different screen aspect-ratio / buttons layout) and duplicates games so you end up with only a clean list of games.

As an example, I use `mame 0.78` and `FBA 0.2.97.44` and an arcade cabinet with a 4x4 screen and 6 buttons.

On `neogeo`, roms-manager hides 143 our of 268
On `fba_libretro`, roms-manager hides 745 out of 1837
On `mame`, roms-manager hides 179 games out of 342
Then roms-manager hides 247 games from `neogeo` and `fba_libretro` that are already present in `mame`

## Install

### Latest release

go to the [release page](https://github.com/vrgl117-games/roms-manager/releases) and grab the binary.

### Development build

Install golang on your machine and then `go get github.com/vrgl117-games/roms-manager`

## How to use

roms-manager has two mains functions

### Scan

`scan` takes a database file (either `.dat` or a Mame `.xml`) and an EmulationStation `gameslist.xml` file.

Features: 
* hide incompatible games (wrong rom size, wrong CRC)
* hide original game if a working clone is present
* hide games using a list of keywords (bootlegs, hacks, etc...)
* if present in the database file, hide games with the incorrect aspect ration or button layout

See `roms-manager scan --help` for the list of flags.

### Hide-Duplicates

`hide-duplicates` takes a "master" `gameslist.xml` file  and a list of "secondary" `gameslist.xml` file(s).

Feature:
* hide games from the secondary `gameslist.xml` file(s) that are already present in the master one.

## TODO

* add support for `catver.ini`
* download `.dat` files on the fly
# Change Log

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org/).

## [1.5.1] - 2020-10-30

* Added support for BBCode as the initial mode during startup.
  * Use either `dropbox-gif-linker bbcode` or `dropbox-gif-linker b`.

## [1.5.0] - 2020-10-29

* Added support for BBCode.
  * Enable using either the `bbcode` or `b` commands.

## [1.4.1] - 2018-11-07

* Better assortment of emojis.

## [1.4.0] - 2018-11-06

* Cached images are now automatically checked to verify that they are currently valid on Dropbox.
* Updated the emojis used.

## [1.3.0] - 2018-10-12

* Added support for `delete` and `del` commands, to be used after the loading of a gif that has a
bad Dropbox URL.

## [1.2.0] - 2018-10-12

* Moved to new `go mod` support for dependencies.

## [1.1.0] - 2018-09-16

* Migrated from SQLite to BoltDB, a pure Go key-value store, allowing cross-compiling again.

## [1.0.0] - 2018-05-14

* Binary available only for MacOS 64-Bit due to cross-compiling SQLite not being simple.

## [1.0.0-rc4] - 2018-05-11

* Fixed missing config error.
* Better default colors.

## [1.0.0-rc3] - 2018-05-11

* Fixed tag displays and ID-based lookups.

## [1.0.0-rc2] - 2018-05-11

* Bugfix for version displayed in program.

## [1.0.0-rc1] - 2018-05-11

* SQLite support.
* Better `help` commands.
* Better Dropbox support.

## [0.6.0] - 2018-05-09

* Supports version argument and displays it on boot.
  * `./dropbox-gif-linker version` -> `dropbox-gif-linker version 0.60 darwin/amd64`
* Copies to clipboard, and keeps the mode in mind.
* Basic help support.

## [0.5.0] - 2018-05-09

* Supports a configuration file, `~/.dgl.json`. See the README for details.
* Creates and retrieve DropBox links

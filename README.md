sushimaster
===========

POC: Single binary executable many commands like busybox

## Install

````
go get github.com/riywo/sushimaster
````

## Usage

Prepare a directory like below:

````
$ tree data/
data/
└── bin
    ├── bar
    └── foo
$ cat data/bin/foo
#!/bin/bash
echo foo
$ ./data/bin/foo
foo
$ ./data/bin/bar
bar
````

Then, make `sushibox` by `sushimaster` and have fun.

````
$ sushimaster data
$ ls
data    sushibox
$ file sushibox
sushibox: Mach-O 64-bit executable x86_64
$ ./sushibox foo
foo
$ ./sushibox bar
bar
````

Symlink works like `buxybox`.

````
$ ln -s sushibox foo
$ ln -s sushibox foo
$ ./foo
foo
$ ./bar
bar
````

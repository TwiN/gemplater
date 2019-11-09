# gemplater

Gemplater is a Go application that allows you to add simple templating to your projects.

It was originally created to allow some extra customizability to dotfiles repositories.


## Usage

The magic happens with the `install` subcommand, which _installs_, or _processes_ 
the target file(s) and outputs them to the destination specified or to stdout
if no destination is specified

```
gemplater install TARGET [DESTINATION] [FLAGS]
```

Examples:

```
gemplater install .profile ~/.profile
gemplater install dotfiles ~/ --quick --remember
```

Flags:

```
-h, --help       help for install
-i, --ignore     Whether to ignore missing variables. If not set, missing variables will trigger interactive mode
-q, --quick      Do not ask for value of variables that are already set. Requires -i to not be set
-r, --remember   Remember variables interactively set on one file for other files. Requires -i to not be set. Useless if TARGET is not directory
```

If you do not specify the `DESTINATION`, the output will be printed in the console.

...

default config file: `.gemplater.yml`

# gemplater

Gemplater is a Go application that allows you to add simple templating to your projects.

It was originally created to allow some extra customizability to dotfiles repositories.


## Usage


```
gemplater install FILE [DESTINATION]
```

e.g.:

```
gemplater install .profile ~/.profile
```

If you do not specify the `DESTINATION`, the output will be printed in the console.


...

default config file: `.gemplater/defaults`

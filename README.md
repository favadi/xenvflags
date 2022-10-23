# xenvflags

xenvflags is an utility that wrap around cli programs and automatically apply
extra argument from environment variables. It is best used together with
[direnv](https://direnv.net/).

## Installation

```shell
# with homebrew
$ brew install favadi/xenvflags/xenvflags

# with go
$ go get github.com/favadi/xenvflags
```

Or download the pre-built binaries in [releases
section](https://github.com/favadi/xenvflags/releases).


## Example usages

For example, we are using [shfmt](https://github.com/mvdan/sh) for both
personal and work projects.

- personal projects at ~/code/personal use 4 spaces indentation (arguments: `-i 4`)
- work projects are at ~/code/work use use 2 spaces indentation (arguments: `-i 2`)

The shfmt cli program is installed in /usr/local/bin/shfmt.

Setup:

- create a symlink to xenvflags: `ln -s /usr/local/bin/xenvflags $HOME/bin/shfmt`
- append `SHFMT_EXTRA_ARGS='-i 4'` to ~/code/personal/.envrc
- append `SHFMT_EXTRA_ARGS='-i 2'` to ~/code/work/.envrc

With [direnv hook](https://direnv.net/docs/hook.html) enabled, running `shfmt`
in personal and work projects will have the desired indentation flags append.

### Debug

```shell
# print xenvflags version and exit
$ XENVFLAGS_VERSION=true shfmt
dev

# print debug information
$ SHFMT_EXTRA_ARGS="-i 2" XENVFLAGS_DEBUG=true shfmt --help
2022/10/23 16:40:36 version: dev
2022/10/23 16:40:36 executable: /usr/local/bin/shfmt
2022/10/23 16:40:36 original arguments: [--help]
2022/10/23 16:40:36 extra arguments from env: [-i 2]
...
```

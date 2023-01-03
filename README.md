# javadoc2md

`javadoc2md` is a Javadoc transpiler which extracts documentation from Java
files, and writes Docusaurus-compatible markdown files.

For the most part, this project does not test whether the markdown generated
could be used for other purposes.

## Installation

You can install the latest version of `javadoc2md` to your GOPATH using 
the following command:

```shell
> go install github.com/dburkart/javadoc2md
```

## Usage

```
Usage of javadoc2md:
  -input string
    Input directory to transpile (default ".")
  -output string
    Output directory to receive markdown files (default ".")
  -skip-private
    Skip private definitions
```

## Limitations

Since this transpiler is written in Go, and it's operating over essentially
what is Java syntax, there are a few caveats which could result in weirdness:

  * References to standard library functions / classes / etc. do yet resolve.
    It's not clear whether I'll ever make them resolve, since I'm not sure what
    the best way to do so is yet.
  * Some bits of Java syntax are not yet understood by the parser, i.e. generics
    and the like.

Additionally, since this project is still undergoing active development, thare are
not answers to some questions yet, such as:

  * How best to lay out files on disk
  * How best to expose on-disk layout for configuration purposes
  * What sorts of metadata to include in the markdown files

## Filing Bugs

If you find something which should work but doesn't, please don't hesitate to file
an issue on this repository with a reduced test case, and the output you would
reasonably expect.

# The Go Programming Language

## Modified to emit warnings

This git commit modifies the go compiler to support emitting warnings instead of errors when it encounters unused labels, variables, and imports. Unused things do not stop compilation, testing, or production of the final binary when the `-warnunused` flag is set. Compiling without the flag makes it behave as normal, generating errors for unused things.

**All changes to the upstream go compiler are contained within this one commit.** From time to time, I'll rebase this commit onto upstream master and force-push to https://github.com/kstenerud/go master in order to keep it clean and current, and also create `goX.Y.Z-warn` tags you can build from.

* [Use Case](#use-case)
* [Usage](#usage)
* [When will this go into mainline?](#when-will-this-go-into-mainline)
* [Getting and Building the Modified Go Compiler](#getting-and-building-the-modified-go-compiler)
* [Keeping up to date with master](#keeping-up-to-date-with-master)
* [Applying this Patch to a Specific Go Release](#applying-this-patch-to-a-specific-go-release)
* [Rebasing to golang upstream](#rebasing-to-golang-upstream)
* [Helper Script](#helper-script)
* [Original Document](#original-document)



## Use Case

The use case for warnings is the exploratory development or debugging phases, where you really don't care about leaving unused things lying around for the time being (for example, temporarily commenting something out), and would rather that the compiler just got out of your way until you've got something ready to compile normally and commit.



## Usage

    go build -gcflags=-warnunused somefile.go
    go test -gcflags=-warnunused

When you're done with your exploratory/debugging phase, simply build or test without the flag:

    go build somefile.go
    go test

Compiling without the flags will fail on unused things as normal.

### Example:

example.go:

```golang
package main

import "fmt"

func main() {
    var start int = 1

    breakOuter:
    // for x := start; x < 10; x++ {
    for x := 3; x < 10; x++ {
        for y := 0; y < 10; y++ {
            result := y * 10 + x
            // fmt.Printf("Result: %d\n", result)
            // if result == 42 {
            //  fmt.Printf("Breaking")
            //  break breakOuter
            // }
        }
    }
}
```

Result of `go build -gcflags=-warnunused example.go`:

    ./example.go:8:2: Warning: label breakOuter defined and not used
    ./example.go:3:8: Warning: imported and not used: "fmt"
    ./example.go:6:6: Warning: start declared and not used
    ./example.go:12:4: Warning: result declared and not used



## When will this go into mainline?

This feature will not be put into mainline golang: https://golang.org/doc/faq#unused_variables_and_imports



## Getting and Building the Modified Go Compiler

Download the source code, install a bootstrap compiler, and build!


### 1: Download and install the source code

#### Option A: Downloaded as a tag directly from github:

* Go to https://github.com/kstenerud/go/releases
* Select the latest `goX.Y.Z-warn` tag

#### Option B: Checked out of this git repository:

* `git tag` to get a list of tags
* `git checkout goX.Y.Z-warn`

The go tree must reside at /usr/local/go (according to https://golang.org/doc/install#install)

So for example: `sudo tar -C /usr/local -xzf go1.12.6-warn.tgz`


### 2: Download and install a bootstrap compiler

You can bootstrap from one of the official distributions at https://golang.org/dl

It must reside at `$HOME/go1.4`

For example:

    wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz
    tar xf go1.12.7.linux-amd64.tar.gz
    mv go $HOME/go1.4


### 3: Copy over the bootstrap go tool

    mkdir -p /usr/local/go/bin
    cp $HOME/go1.4/bin/go /usr/local/go/bin/


### 4: Build the compiler

    cd /usr/local/go/src
    ./make.bash

This will build the go compiler binaries and put them in `/usr/local/go`



## Keeping up to date with master

In order to keep things simple, all changes from golang master are contained within this one commit. But this also means that I'll be replacing this one commit with a new one on every change.

To update to the latest `master` branch of this repo, you must replace your git head:

```
git fetch --tags
git reset --hard origin/master
```

**Warning: this will destroy all your local changes!**



## Applying this Patch to a Specific Go Release

To apply this change to a specific go release, simply checkout the tag, make a branch, cherry-pick this commit, and build.


### Release patching example

Get the tag to cherry-pick:

	$ git checkout master
    $ git log --oneline -n 1
    0e367a3c38 (HEAD -> master, origin/master, origin/HEAD) Support emitting warnings ...

Make a branch from your chosen tag and cherry-pick:

    git checkout go1.12.6
    git checkout -b go1.12.6-warn-branch
    git cherry-pick 0e367a3c38

Build the compiler:

    cd src
    ./make.bash

Optional: Run tests also: `./all.bash`

Make a new tag for your modified go compiler:

    git tag -a go1.12.6-warn



## Rebasing to golang upstream

To get the latest goodies from golang upstream while retaining this patch:

```
git pull --rebase https://github.com/golang/go.git master
git fetch --tags https://github.com/golang/go.git
```



## Helper Script

Here's a quick script to put this change into a specific go release. If using the commit from from `master` (the default) causes conflicts, try using one from a previous `warn` tag instead. For example:

    make-go-patched-branch -s go1.13.3-warn go1.13.4

make-go-patched-branch:

```bash
#!/bin/bash

set -eu

if [[ ! $(git config --get remote.origin.url) =~ .*/go\.git$ ]]; then
    echo "This command must be run from inside the golang git repository"
    exit 1
fi

GO_RELEASE_SOURCE_TAG=master

show_help() {
    echo "Usage: $0 [options] <go release tag>"
    echo "Use git tag to get a list"
    echo
    echo "Options:"
    echo " -s <tag>: Use the commit from the specified warn tag instead of $GO_RELEASE_SOURCE_TAG (example: go1.13.4-warn)"
}

while getopts "s:" o; do
    case "$o" in
        s)
            GO_RELEASE_SOURCE_TAG=$OPTARG
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done
shift $((OPTIND-1))


if [[ $# -ne 1 ]]; then
    show_help
    exit 1
fi

GO_RELEASE_ORIGINAL_TAG=$1
GO_RELEASE_NEW_TAG=${GO_RELEASE_ORIGINAL_TAG}-warn
GO_RELEASE_BRANCH=${GO_RELEASE_NEW_TAG}-branch

git checkout $GO_RELEASE_SOURCE_TAG
WARNINGS_COMMIT=$(git rev-parse HEAD)

git checkout "$GO_RELEASE_ORIGINAL_TAG"
git checkout -b "$GO_RELEASE_BRANCH"
git cherry-pick $WARNINGS_COMMIT

echo "
New branch $GO_RELEASE_BRANCH created with warnings commit cherry-picked ($(git rev-parse HEAD))

Now you can build it:
    cd src
    ./make.bash or ./all.bash

When you're happy, tag it:
    git tag -a $GO_RELEASE_NEW_TAG
"
```


## Original Document:

Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.

![Gopher image](doc/gopher/fiveyears.jpg)
*Gopher image by [Renee French][rf], licensed under [Creative Commons 3.0 Attributions license][cc3-by].*

Our canonical Git repository is located at https://go.googlesource.com/go.
There is a mirror of the repository at https://github.com/golang/go.

Unless otherwise noted, the Go source files are distributed under the
BSD-style license found in the LICENSE file.

### Download and Install

#### Binary Distributions

Official binary distributions are available at https://golang.org/dl/.

After downloading a binary release, visit https://golang.org/doc/install
or load [doc/install.html](./doc/install.html) in your web browser for installation
instructions.

#### Install From Source

If a binary distribution is not available for your combination of
operating system and architecture, visit
https://golang.org/doc/install/source or load [doc/install-source.html](./doc/install-source.html)
in your web browser for source installation instructions.

### Contributing

Go is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines:
    https://golang.org/doc/contribute.html

Note that the Go project uses the issue tracker for bug reports and
proposals only. See https://golang.org/wiki/Questions for a list of
places to ask questions about the Go language.

[rf]: https://reneefrench.blogspot.com/
[cc3-by]: https://creativecommons.org/licenses/by/3.0/

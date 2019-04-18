# The Go Programming Language

## Modified to emit warnings

This git commit modifies the go programming language to support emitting warnings instead of errors when it encounters unused labels, variables, and imports. Warnings do not stop compilation or production of the final binary.

The use case for warnings is the exploratory development or debugging phases, where you really don't care about leaving unused things lying around for the time being (for example, temporarily commenting something out), and would rather that the compiler just got out of your way until you've got something ready to compile normally and commit.

Usage:

	go build -gcflags="-warnunused" somefile.go
	go test -gcflags="-warnunused"


### Example

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

Result of `go build -gcflags="-warnunused" example.go`:

	./example.go:8:2: Warning: label breakOuter defined and not used
	./example.go:3:8: Warning: imported and not used: "fmt"
	./example.go:6:6: Warning: start declared and not used
	./example.go:12:4: Warning: result declared and not used


This feature will not go into mainline: https://golang.org/doc/faq#unused_variables_and_imports


### Applying to a Specific Go Release

To apply this change to a specific go release, simply checkout the tag, make a branch, cherry-pick this commit, and build.

Example:

	git checkout go1.12.4
	git checkout -b go1.12.4-warn
	git cherry-pick 6a176eca2f977cff401370f72f4b9f328dbb96fb
	cd src
	./make.bash

Here's a quick script to create a modified branch:

```bash
#!/bin/bash

set -eu

if [[ ! $(git config --get remote.origin.url) =~ .*/go\.git$ ]]; then
    echo "This command must be run from inside the golang git repository"
    exit 1
fi

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <go release tag>"
    echo "Use git tag to get a list"
    exit 1
fi

GO_RELEASE_ORIGINAL_TAG=$1
GO_RELEASE_NEW_TAG=${GO_RELEASE_ORIGINAL_TAG}-warn

git checkout master

WARNING_COMMIT=$(git rev-parse HEAD)

git checkout "$GO_RELEASE_ORIGINAL_TAG"
git checkout -b "$GO_RELEASE_NEW_TAG"
git cherry-pick $WARNING_COMMIT

echo
echo "New branch $GO_RELEASE_NEW_TAG created with warnings commit cherry-picked ($(git rev-parse HEAD))"
echo "Next steps:"
echo "  cd src"
echo "  ./all.bash or ./make.bash"
echo "  git tag -a $GO_RELEASE_NEW_TAG"
echo "Point GOROOT to where go was built"
echo "Modify your PATH to point to \$GOROOT/bin"
```

I'll rebase this commit onto master and force-push from time to time in order to keep it clean and current.


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

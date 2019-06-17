// Hacky way to make global flags aaccessible, because the alternative is
// re-architecting the compiler to pass down these flags to code that isn't
// designed to have metadata like this passed in.
//
// Author: Karl Stenerud
package warnings

// Treat unused labels, variables & imports as warnings.
// This value is set in cmd/compiler/internal/gc/main.go, and can be
// controlled via gcflags:
//
//     go build -gcflags="-warnunused" myfile.go
//
// If the gcflag is not set, unused things are errors.
// If the gcflag is set, unused things are warnings.
//
// Intended use case: Set this flag during development to get the compiler
// out of your way, and then compile without the flag before committing.
//
// This feature will not go into mainline: https://golang.org/doc/faq#unused_variables_and_imports
//
// Note: This flag is true by default so that the vet command only warns by default.
// This is "safe enough" since the build phase already does the same checks,
// and so the process would never normally get to the vet command anyway.
// In future, I'll try to patch vet as well, but its flag system is tied
// directly to the analyzer list so it'll get messy.
// The alternative would be to call "go test" with "--vet=off", which I think
// is an even worse solution.
var treatUnusedAsWarning bool = true

func TreatUnusedAsWarning(unusedIsWarning bool) {
	treatUnusedAsWarning = unusedIsWarning
}

func IsUnusedTreatedAsError() bool {
	return !treatUnusedAsWarning;
}

func IsUnusedTreatedAsWarning() bool {
	return treatUnusedAsWarning;
}

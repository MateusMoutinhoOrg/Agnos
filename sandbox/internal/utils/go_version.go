package utils

// GoRelease is the Go release a generated project is declared against: the
// version its go.mod names and the one its docs tell a reader to install.
//
// It is a constant rather than whatever `go env GOVERSION` answers on the
// machine that runs `start`. A `go` directive is a floor, not a pin — a newer
// toolchain builds a module naming an older release — so reading the local
// version only pins the project to the day it was scaffolded, and makes it
// unbuildable on the version its own docs promise.
const GoRelease = "1.25.0"

// GoFloor is GoRelease without its patch, the "1.25" a doc spells when it says
// which releases are new enough. It is the same number: change one and change
// the other.
const GoFloor = "1.25"

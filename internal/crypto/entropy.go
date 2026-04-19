package crypto

import (
	"crypto/rand"
	"io"
)

// entropyReader stays private so tests can deterministically exercise
// entropy failure paths without changing the public API.
var entropyReader io.Reader = rand.Reader

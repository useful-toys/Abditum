package storage

import (
	"fmt"

	"github.com/useful-toys/abditum/internal/crypto"
)

// Binary file format constants for the .abditum file format.
//
// Header layout (49 bytes total):
//
//	Offset  0: Magic         [4]byte  -- "ABDT"
//	Offset  4: Version       uint8    -- format version
//	Offset  5: Salt          [32]byte -- Argon2id salt
//	Offset 37: Nonce         [12]byte -- AES-256-GCM nonce
//	Offset 49: Payload       []byte   -- encrypted JSON (AES-256-GCM ciphertext + tag)
//
// The full 49-byte header (offsets 0-48) is used as GCM AAD.

// Magic is the 4-byte file signature identifying the .abditum format.
var Magic = [4]byte{'A', 'B', 'D', 'T'}

// Field sizes (in bytes).
const (
	MagicSize   = 4
	VersionSize = 1
	SaltSize    = 32
	NonceSize   = 12
)

// Field offsets within the header.
const (
	SaltOffset  = MagicSize + VersionSize // 5
	NonceOffset = SaltOffset + SaltSize   // 37
)

// HeaderSize is the total size of the fixed binary header in bytes.
// 4 (magic) + 1 (version) + 32 (salt) + 12 (nonce) = 49
const HeaderSize = MagicSize + VersionSize + SaltSize + NonceSize // 49

// CurrentFormatVersion is the format version written to new .abditum files.
const CurrentFormatVersion uint8 = 1

// FileMetadata holds the information used to detect external file changes.
//
// External change detection uses file size as a fast path (O(1), no read required)
// and SHA-256 of the full file as confirmation. Modification time (mtime) is
// deliberately NOT used due to filesystem and cloud sync inconsistencies
// (see arquitetura.md §7 and historico-decisoes.md §17).
type FileMetadata struct {
	Size int64
	Hash [32]byte
}

// FormatProfile holds the Argon2id parameters for a specific format version.
//
// Parameters are NOT stored in the file header -- they are derived from the
// format version via ProfileForVersion. Bumping the format version means adding
// a new FormatProfile entry and a migration path.
type FormatProfile struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// formatProfiles maps format version numbers to their Argon2id parameters.
//
// Version 1: Argon2id t=3, m=256 MiB (262144 KiB), p=4, keyLen=32
// (per formato-arquivo-abditum.md §3)
var formatProfiles = map[uint8]FormatProfile{
	1: {Time: 3, Memory: 262144, Threads: 4, KeyLen: 32},
}

// ProfileForVersion returns the FormatProfile for the given format version.
//
// Returns ErrVersionTooNew if the version is not in the known profiles map.
// The caller should use the returned profile to configure key derivation.
func ProfileForVersion(version uint8) (FormatProfile, error) {
	profile, ok := formatProfiles[version]
	if !ok {
		return FormatProfile{}, fmt.Errorf("%w: version %d", ErrVersionTooNew, version)
	}
	return profile, nil
}

// ToArgonParams converts the FormatProfile to a crypto.ArgonParams suitable
// for passing to crypto.DeriveKey.
func (p FormatProfile) ToArgonParams() crypto.ArgonParams {
	return crypto.ArgonParams{
		Time:    p.Time,
		Memory:  p.Memory,
		Threads: p.Threads,
		KeyLen:  p.KeyLen,
	}
}

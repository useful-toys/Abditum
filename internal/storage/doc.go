// Package storage implements the binary .abditum file format for Abditum password manager.
//
// # File Format
//
// The .abditum format uses a 49-byte fixed binary header followed by an
// AES-256-GCM encrypted JSON payload. The header layout is:
//
//   - Bytes  0-3:  Magic ("ABDT", 4 bytes ASCII)
//   - Byte   4:    Format version (uint8)
//   - Bytes  5-36: Argon2id salt (32 bytes)
//   - Bytes 37-48: AES-256-GCM nonce (12 bytes)
//   - Bytes 49+:   Encrypted payload (JSON via vault.SerializarCofre)
//
// The full 49-byte header is used as GCM Additional Authenticated Data (AAD),
// meaning any tampering with header bytes causes authentication failure.
//
// The canonical format specification is in formato-arquivo-abditum.md.
//
// # Atomic Writes
//
// Saves use a .tmp -> rename protocol with a .bak/.bak2 backup chain to ensure
// the vault file is never left in a partially-written state. A startup orphan
// recovery step (RecoverOrphans) cleans up any .tmp files left by a crash.
//
// # External Change Detection
//
// The package detects whether the vault file was modified externally (e.g., by
// cloud sync) using file size as a fast path and SHA-256 as confirmation.
// Modification time (mtime) is deliberately NOT used due to filesystem and
// cloud sync inconsistencies.
//
// # Format Version Migration
//
// The version-to-profile mapping (ProfileForVersion) maps format version numbers
// to Argon2id parameters. Bumping the format version adds a new profile entry
// and a migration path; Argon2id parameters are never stored in the header.
package storage

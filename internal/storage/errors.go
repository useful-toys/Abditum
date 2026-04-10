package storage

import "errors"

// Sentinel errors for the storage package.
//
// These errors are returned by storage functions to indicate specific failure
// conditions. They are designed to be checked with errors.Is().
//
// Note: crypto.ErrAuthFailed is used directly for wrong-password scenarios
// (GCM tag verification failure). It is not redeclared here.

// ErrInvalidMagic is returned when the file does not start with the "ABDT" magic bytes.
//
// This typically means the file is not a valid .abditum vault file, or it has
// been truncated before the header could be read.
var ErrInvalidMagic = errors.New("invalid file type")

// ErrVersionTooNew is returned when the format version in the header exceeds
// the highest version supported by this build.
//
// The user should upgrade Abditum to open this vault file.
var ErrVersionTooNew = errors.New("unsupported format version")

// ErrCorrupted is returned when the file passes magic and authentication checks
// but the decrypted content fails structural validation.
//
// Typical causes:
//   - Invalid JSON in the decrypted payload
//   - Missing or invalid pasta_geral in the deserialized vault
//   - pasta_geral.nome is not "Geral"
//
// This indicates the file contents are corrupted beyond cryptographic protection.
var ErrCorrupted = errors.New("file integrity check failed")

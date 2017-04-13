package crypto

import (
	"encoding/binary"

	"github.com/libeclipse/tranquil/memory"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/scrypt"
)

// DeriveSecureValues derives and returns a masterKey and rootIdentifier.
func DeriveSecureValues(masterPassword, identifier []byte, costFactor map[string]int) (*[32]byte, []byte) {
	// Allocate and protect memory for the concatenated values, and append the values to it.
	concatenatedValues := make([]byte, len(masterPassword)+len(identifier))
	memory.Protect(concatenatedValues)
	concatenatedValues = append(masterPassword, identifier...)

	// Allocate and protect memory for the output of the hash function, and put the output into it.
	rootKeySlice := make([]byte, 64)
	memory.Protect(rootKeySlice)
	rootKeySlice, _ = scrypt.Key(
		concatenatedValues,       // Input data.
		[]byte(""),               // Salt.
		1<<uint(costFactor["N"]), // Scrypt parameter N.
		costFactor["r"],          // Scrypt parameter r.
		costFactor["p"],          // Scrypt parameter p.
		64)                       // Output hash length.

	// Allocate a protected array to hold the key, and copy the key into it.
	var masterKey [32]byte
	memory.Protect(masterKey[:])
	copy(masterKey[:], rootKeySlice[0:32])

	// Slice and return respective values.
	return &masterKey, rootKeySlice[32:64]
}

// DeriveIdentifierN derives a value for derivedIdentifier for a value of `n`.
func DeriveIdentifierN(rootIdentifier []byte, n int) []byte {
	// Convert n to a byte slice.
	byteN := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteN, uint32(n))

	// Derive derivedIdentifier.
	derivedIdentifier := blake2b.Sum256(append(rootIdentifier, byteN...))

	// Return as slice instead of array.
	return derivedIdentifier[:]
}
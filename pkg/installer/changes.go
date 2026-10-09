// Package installer encodes the changes column of the installer tables: the commits of a build, newest first, stored by
// the collector as unpadded base64 of the SHA-1 bytes (a third shorter than hex).
package installer

import (
	"encoding/base64"
	"encoding/hex"
)

// EncodeCommit returns the stored form of a hex commit, false when it is not hex.
func EncodeCommit(hexCommit string) (string, bool) {
	raw, err := hex.DecodeString(hexCommit)
	if err != nil {
		return "", false
	}
	return base64.RawStdEncoding.EncodeToString(raw), true
}

// DecodeCommit returns the hex commit of a stored change, false when it is not base64.
func DecodeCommit(change string) (string, bool) {
	raw, err := base64.RawStdEncoding.DecodeString(change)
	if err != nil {
		return "", false
	}
	return hex.EncodeToString(raw), true
}

// DecodeChanges returns the commits as hex. A value that is not base64 is kept as is, so a switch of the stored format
// (e.g. to hex) shows up instead of silently corrupting the output.
func DecodeChanges(changes []string) []string {
	commits := make([]string, len(changes))
	for i, c := range changes {
		if commit, ok := DecodeCommit(c); ok {
			commits[i] = commit
		} else {
			commits[i] = c
		}
	}
	return commits
}

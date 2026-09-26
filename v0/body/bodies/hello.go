package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v0/body"
	"github.com/dejitarudemon/axidb-go-protocol/v0/err"
	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

const (
	// MaxVersionsPerOneHello is the maximum number of version bytes in one Hello body.
	MaxVersionsPerOneHello = 255
)

var _ body.Body = Hello{}

// Hello is a version-0 body listing protocol versions the sender can speak.
type Hello struct {
	// Versions lists working protocol versions. Version 0 is not a working version.
	Versions []fields.Version
}

// NewHello returns a Hello after dropping version 0, values that do not fit in one
// byte, and duplicate version numbers. The first occurrence of each version is kept.
func NewHello(versions []fields.Version) Hello {
	return Hello{
		Versions: filter(versions),
	}
}

// filter drops version 0, values greater than 255, and duplicate version numbers.
func filter(versions []fields.Version) []fields.Version {
	filtered := make([]fields.Version, 0, len(versions))
	used := make(map[fields.Version]struct{}, len(versions))

	for _, version := range versions {
		if _, ok := used[version]; !ok && version != 0 && version <= MaxVersionsPerOneHello {
			used[version] = struct{}{}
			filtered = append(filtered, version)
		}
	}

	return filtered
}

// Size returns the encoded body size in bytes.
func (h Hello) Size() int {
	return min(len(h.Versions), MaxVersionsPerOneHello) * fields.VersionFieldSize
}

// Encode writes the wire encoding of the body into buf.
func (h Hello) Encode(buf buffer.Appender) {
	for i, version := range h.Versions {
		if i >= MaxVersionsPerOneHello {
			break
		}

		version.Encode(buf)
	}
}

// IsValid reports whether the body satisfies protocol rules.
// Version 0 is rejected because it is only the Hello frame version.
// Values greater than 255 are rejected because a version occupies one byte.
// Duplicates are allowed. The list must not exceed [MaxVersionsPerOneHello].
func (h Hello) IsValid() error {
	if len(h.Versions) > MaxVersionsPerOneHello {
		return err.NewValidationError(
			"too many versions",
			"versions", len(h.Versions),
			"target", "Hello",
		)
	}

	for _, version := range h.Versions {
		if version == 0 {
			return err.NewValidationError(
				"version 0 is not a working version",
				"versions", h.Versions,
				"target", "Hello",
			)
		}

		if version > MaxVersionsPerOneHello {
			return err.NewValidationError(
				"version does not fit in one byte",
				"version", version,
				"target", "Hello",
			)
		}
	}

	return nil
}

// Common returns versions present in both h and other, in the order they appear in h.
// Version 0 and later duplicates are skipped.
func (h Hello) Common(other Hello) []fields.Version {
	offered := make(map[fields.Version]struct{}, len(other.Versions))
	for _, version := range other.Versions {
		if version != 0 {
			offered[version] = struct{}{}
		}
	}

	common := make([]fields.Version, 0)
	seen := make(map[fields.Version]struct{}, len(h.Versions))

	for _, version := range h.Versions {
		if version == 0 {
			continue
		}

		if _, ok := offered[version]; !ok {
			continue
		}

		if _, dup := seen[version]; dup {
			continue
		}

		seen[version] = struct{}{}
		common = append(common, version)
	}

	return common
}

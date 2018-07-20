package origin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Version represents software version
	Version struct {
		major  int
		minor  int
		patch  *int // workaround in case of not exists
		suffix string
	}
)

var reVersion = regexp.MustCompile("^[0-9]+")

// Major returns major version
func (v *Version) Major() int {
	return v.major
}

// Minor returns minor version
func (v *Version) Minor() int {
	return v.minor
}

// String returns version string
func (v *Version) String() string {
	f := fmt.Sprintf("%d.%d", v.major, v.minor)
	if v.patch != nil {
		f = fmt.Sprintf("%s.%d", f, *v.patch)
	}
	if v.suffix != "" {
		f = fmt.Sprintf("%s%s", f, v.suffix)
	}
	return f
}

// Equal v == a
func (v *Version) Equal(a *Version) bool {
	return v.Compare(a) == 0
}

// GreaterThan v > a
func (v *Version) GreaterThan(a *Version) bool {
	return v.Compare(a) == 1
}

// GreaterThanEqualString v >= a
func (v *Version) GreaterThanEqualString(a string) bool {
	p := MustParseVersion(a)
	return v.Compare(p) == 0 || v.Compare(p) == 1
}

// GreaterThanString v > a
func (v *Version) GreaterThanString(a string) bool {
	return v.Compare(MustParseVersion(a)) == 1
}

// LessThan v < a
func (v *Version) LessThan(a *Version) bool {
	return v.Compare(a) == -1
}

// LessThanString v < a
func (v *Version) LessThanString(a string) bool {
	return v.Compare(MustParseVersion(a)) == -1
}

// Compare compares given version
func (v *Version) Compare(a *Version) int {
	if v.major > a.major {
		return 1
	} else if v.major < a.major {
		return -1
	} else {
		if v.minor > a.minor {
			return 1
		} else if v.minor < a.minor {
			return -1
		} else {
			var vp int
			if v.patch != nil {
				vp = *v.patch
			}

			var ap int
			if a.patch != nil {
				ap = *a.patch
			}

			if vp > ap {
				return 1
			} else if vp < ap {
				return -1
			}
		}
	}
	return 0
}

// ParseVersion parses `x.x[.x][-yyy]` string
func ParseVersion(v string) (*Version, error) {
	sv := strings.Split(v, "-")
	var suffix string
	if len(sv) > 1 {
		suffix = fmt.Sprintf("-%s", sv[1])
	}

	sv = strings.Split(sv[0], ".")
	if len(sv) != 2 && len(sv) != 3 {
		return nil, errors.Errorf("version parse: unknown version:%v", v)
	}

	major, err := strconv.Atoi(sv[0])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(sv[1])
	if err != nil {
		return nil, err
	}

	var patch *int
	if len(sv) == 3 {
		match := reVersion.FindAllString(sv[2], -1)
		p, err := strconv.Atoi(match[0])
		if err != nil {
			return nil, err
		}
		patch = &p

		// e.g. extract rc1 from 1.2.3rc1
		match = reVersion.Split(sv[2], -1)
		if len(match) >= 2 {
			suffix = match[1] + suffix
		}
	}

	return &Version{
		major:  major,
		minor:  minor,
		patch:  patch,
		suffix: suffix,
	}, nil
}

// MustParseVersion got panic if parse failure
func MustParseVersion(a string) *Version {
	v, err := ParseVersion(a)
	if err != nil {
		panic(err)
	}
	return v
}

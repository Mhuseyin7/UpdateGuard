package update

import (
	"regexp"
	"strconv"
	"strings"
)

type VersionChange string

const (
	Unknown VersionChange = "UNKNOWN"
	Patch   VersionChange = "PATCH"
	Minor   VersionChange = "MINOR"
	Major   VersionChange = "MAJOR"
	Same    VersionChange = "SAME"
)

var semver = regexp.MustCompile(`^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

// ClassifySemver classifies only fully semantic tags. Registry tags that are not SemVer are unknown by design.
func ClassifySemver(current, target string) VersionChange {
	a, b := semver.FindStringSubmatch(strings.TrimSpace(current)), semver.FindStringSubmatch(strings.TrimSpace(target))
	if a == nil || b == nil {
		return Unknown
	}
	av := []int{}
	bv := []int{}
	for i := 1; i <= 3; i++ {
		x, _ := strconv.Atoi(a[i])
		y, _ := strconv.Atoi(b[i])
		av = append(av, x)
		bv = append(bv, y)
	}
	if av[0] != bv[0] {
		return Major
	}
	if av[1] != bv[1] {
		return Minor
	}
	if av[2] != bv[2] || a[0] != b[0] {
		return Patch
	}
	return Same
}

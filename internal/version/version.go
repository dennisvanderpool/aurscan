// Package version exposes aurscan's build version through a robust fallback
// chain so `--version` is meaningful no matter how the binary was produced:
//
//  1. Version/Commit/Date stamped at link time via -ldflags -X (set by the
//     Makefile from `git describe`, or by the PKGBUILD from $pkgver — this is
//     the only reliable source for AUR builds, which have no .git directory).
//  2. Go's built-in VCS stamping from debug.ReadBuildInfo() — populated
//     automatically by `go build` inside a git checkout, so even
//     `go install github.com/manticore-projects/aurscan/cmd/aurscan@latest`
//     yields a real revision.
//  3. A hardcoded fallback ("dev") when nothing else is available.
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// These are overridden at link time:
//
//	-ldflags "-X github.com/manticore-projects/aurscan/internal/version.Version=v0.2.0 ..."
var (
	Version = "" // e.g. v0.2.0 or v0.2.0-3-gabc1234-dirty
	Commit  = "" // git short SHA
	Date    = "" // build date (RFC3339 or YYYY-MM-DD)
)

const fallback = "dev"

// info resolves the effective version fields, applying the fallback chain.
func info() (ver, commit, date string, modified bool) {
	ver, commit, date = Version, Commit, Date

	// Fill any gaps from Go's embedded build info / VCS stamps.
	if bi, ok := debug.ReadBuildInfo(); ok {
		var vcsRev, vcsTime string
		var vcsMod bool
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				vcsRev = s.Value
			case "vcs.time":
				vcsTime = s.Value
			case "vcs.modified":
				vcsMod = s.Value == "true"
			}
		}
		if commit == "" && vcsRev != "" {
			commit = vcsRev
			if len(commit) > 12 {
				commit = commit[:12]
			}
		}
		if date == "" {
			date = vcsTime
		}
		// bi.Main.Version is "(devel)" for local builds but a real tag for
		// `go install ...@vX.Y.Z`; prefer it over the fallback only.
		if ver == "" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			ver = bi.Main.Version
		}
		modified = vcsMod
	}

	if ver == "" {
		ver = fallback
	}
	return ver, commit, date, modified
}

// Short returns just the version string (with a -dirty suffix if the working
// tree was modified and the stamped version doesn't already say so).
func Short() string {
	ver, _, _, modified := info()
	if modified && !hasSuffix(ver, "-dirty") {
		ver += "-dirty"
	}
	return ver
}

// String returns a full multi-field version line for `--version`.
func String() string {
	ver, commit, date, modified := info()
	if modified && !hasSuffix(ver, "-dirty") {
		ver += "-dirty"
	}
	s := "aurscan " + ver
	if commit != "" {
		s += " (" + commit + ")"
	}
	if date != "" {
		s += " built " + date
	}
	s += fmt.Sprintf("\n%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	return s
}

func hasSuffix(s, suf string) bool {
	return len(s) >= len(suf) && s[len(s)-len(suf):] == suf
}

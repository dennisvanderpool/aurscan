package yay

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildsPackages(t *testing.T) {
	tests := []struct {
		name string
		argv []string
		want bool
	}{
		{"bare yay defaults to -Syu", nil, true},
		{"search-install term", []string{"firefox"}, true},
		{"sync install", []string{"-S", "firefox"}, true},
		{"sync refresh+upgrade", []string{"-Syu"}, true},
		{"sync upgrade only", []string{"-Su"}, true},
		{"sync long", []string{"--sync", "firefox"}, true},
		{"sync download-only builds", []string{"-Sw", "foo"}, true},
		{"sync refresh only", []string{"-Sy"}, false},
		{"sync double refresh", []string{"-Syy"}, false},
		{"sync refresh with target builds", []string{"-Sy", "firefox"}, true},
		{"sync print", []string{"-Sp", "foo"}, false},
		{"sync print long", []string{"--sync", "--print", "foo"}, false},
		{"sync help", []string{"-Sh"}, false},
		{"version short", []string{"-V"}, false},
		{"version long", []string{"--version"}, false},
		{"help short", []string{"-h"}, false},
		{"help long", []string{"--help"}, false},
		{"sync search", []string{"-Ss", "foo"}, false},
		{"sync search long", []string{"--sync", "--search", "foo"}, false},
		{"sync info", []string{"-Si", "foo"}, false},
		{"sync list", []string{"-Sl"}, false},
		{"sync clean", []string{"-Sc"}, false},
		{"sync groups", []string{"-Sg"}, false},
		{"query", []string{"-Q"}, false},
		{"query explicit", []string{"-Qe"}, false},
		{"query search", []string{"-Qs", "foo"}, false},
		{"remove", []string{"-R", "foo"}, false},
		{"files", []string{"-F", "foo"}, false},
		{"getpkgbuild", []string{"-G", "foo"}, false},
		{"upgrade from file", []string{"-U", "./foo.pkg.tar.zst"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildsPackages(tt.argv); got != tt.want {
				t.Errorf("buildsPackages(%q) = %v, want %v", tt.argv, got, tt.want)
			}
		})
	}
}

func TestPackageRootFindsPKGBUILDParent(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "PKGBUILD"), []byte("pkgname=test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(dir, "pkg", "usr", "bin")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := packageRoot(nested); got != dir {
		t.Fatalf("packageRoot(%q) = %q, want %q", nested, got, dir)
	}
}

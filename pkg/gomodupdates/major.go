package gomodupdates

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

const defaultMajorLimit = 10

// VersionLister lists known versions for a module path.
type VersionLister interface {
	Versions(ctx context.Context, modulePath string) ([]string, error)
}

// GoVersionLister lists module versions by shelling out to go list.
type GoVersionLister struct {
	GoCommand string
	Mod       string
}

// Versions returns versions for modulePath using go list -m -versions.
func (l GoVersionLister) Versions(ctx context.Context, modulePath string) ([]string, error) {
	cmd := l.GoCommand
	if cmd == "" {
		cmd = "go"
	}

	args := []string{"list"}
	if l.Mod != "" {
		args = append(args, "-mod="+l.Mod)
	}
	args = append(args, "-m", "-versions", "-json", modulePath)

	out, err := exec.CommandContext(ctx, cmd, args...).Output()
	if err != nil {
		if _, ok := errors.AsType[*exec.ExitError](err); ok {
			return nil, nil
		}
		return nil, err
	}

	var mod Module
	if err := json.Unmarshal(out, &mod); err != nil {
		return nil, err
	}
	return mod.Versions, nil
}

// MajorUpdate is the latest discovered major module-path candidate.
type MajorUpdate struct {
	Path    string
	Version string
}

// String returns the module path and version together because v2+ updates often
// require changing the import path.
func (u MajorUpdate) String() string {
	if u.Path == "" || u.Version == "" {
		return ""
	}
	return u.Path + "@" + u.Version
}

// FindMajorUpdate probes conventional semantic import version paths for a newer
// major module candidate.
func FindMajorUpdate(ctx context.Context, modulePath, currentVersion string, lister VersionLister, limit int) (MajorUpdate, error) {
	if lister == nil || modulePath == "" || currentVersion == "" {
		return MajorUpdate{}, nil
	}
	if limit == 0 {
		limit = defaultMajorLimit
	}

	base, mode, start := majorBase(modulePath)
	if currentMajor := versionMajor(currentVersion); currentMajor >= start {
		start = currentMajor + 1
	}
	if start > limit {
		return MajorUpdate{}, nil
	}

	var latest MajorUpdate
	for major := start; major <= limit; major++ {
		candidate := majorPath(base, mode, major)
		versions, err := lister.Versions(ctx, candidate)
		if err != nil {
			return MajorUpdate{}, err
		}
		if len(versions) == 0 {
			break
		}
		if version := latestStable(versions); version != "" {
			latest = MajorUpdate{
				Path:    candidate,
				Version: version,
			}
		}
	}

	return latest, nil
}

type majorMode int

const (
	majorModeSlash majorMode = iota
	majorModeGopkg
)

func majorBase(modulePath string) (string, majorMode, int) {
	base, pathMajor, ok := module.SplitPathVersion(modulePath)
	if !ok {
		return modulePath, majorModeSlash, 2
	}

	mode := majorModeSlash
	if strings.HasPrefix(modulePath, "gopkg.in/") {
		mode = majorModeGopkg
	}

	if pathMajor != "" {
		major, ok := parseMajor(strings.TrimPrefix(module.PathMajorPrefix(pathMajor), "v"))
		if ok {
			return base, mode, major + 1
		}
	}

	return base, mode, 2
}

func majorPath(base string, mode majorMode, major int) string {
	if mode == majorModeGopkg {
		return base + ".v" + strconv.Itoa(major)
	}
	return base + "/v" + strconv.Itoa(major)
}

func versionMajor(version string) int {
	major, ok := parseMajor(strings.TrimPrefix(semver.Major(version), "v"))
	if !ok {
		return 0
	}
	return major
}

func parseMajor(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	major, err := strconv.Atoi(s)
	return major, err == nil
}

func latestStable(versions []string) string {
	stable := make([]string, 0, len(versions))
	for _, version := range versions {
		if semver.IsValid(version) && semver.Prerelease(version) == "" {
			stable = append(stable, version)
		}
	}
	if len(stable) == 0 {
		return ""
	}
	semver.Sort(stable)
	return stable[len(stable)-1]
}

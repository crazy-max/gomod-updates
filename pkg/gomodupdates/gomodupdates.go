package gomodupdates

import (
	"os"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

// Module holds information for one module returned by go list -m -json.
type Module struct {
	Path      string       `json:",omitempty"`
	Version   string       `json:",omitempty"`
	Versions  []string     `json:",omitempty"`
	Replace   *Module      `json:",omitempty"`
	Time      *time.Time   `json:",omitempty"`
	Update    *Module      `json:",omitempty"`
	Main      bool         `json:",omitempty"`
	Indirect  bool         `json:",omitempty"`
	Dir       string       `json:",omitempty"`
	GoMod     string       `json:",omitempty"`
	Error     *ModuleError `json:",omitempty"`
	GoVersion string       `json:",omitempty"`
}

// ModuleError represents an error returned by go list for a module.
type ModuleError struct {
	Err string
}

// Row is a rendered module update row.
type Row struct {
	Module          string
	Version         string
	NewVersion      string
	Direct          bool
	ValidTimestamps bool
}

// HasUpdate reports whether the row has an update.
func (r Row) HasUpdate() bool {
	return r.NewVersion != ""
}

func (m Module) effective() Module {
	if m.Replace != nil {
		return *m.Replace
	}
	return m
}

// CurrentVersion returns the module version, accounting for replace directives.
func (m Module) CurrentVersion() string {
	return m.effective().Version
}

// NewVersion returns the same-path update version, accounting for replace directives.
func (m Module) NewVersion() string {
	mod := m.effective()
	if mod.Update == nil {
		return ""
	}
	return mod.Update.Version
}

// HasUpdate reports whether the module has a same-path update.
func (m Module) HasUpdate() bool {
	return m.effective().Update != nil
}

// InvalidTimestamp reports whether go list returned an update older than the
// current module version timestamp.
func (m Module) InvalidTimestamp() bool {
	mod := m.effective()
	if mod.Time == nil || mod.Update == nil || mod.Update.Time == nil {
		return false
	}
	return mod.Time.After(*mod.Update.Time)
}

// ToolPathsFromGoMod returns tool package paths declared in a go.mod file.
func ToolPathsFromGoMod(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	file, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	tools := make([]string, 0, len(file.Tool))
	for _, tool := range file.Tool {
		tools = append(tools, tool.Path)
	}
	return tools, nil
}

func isToolModule(modulePath string, toolPaths []string) bool {
	for _, toolPath := range toolPaths {
		if toolPath == modulePath || strings.HasPrefix(toolPath, modulePath+"/") {
			return true
		}
	}
	return false
}

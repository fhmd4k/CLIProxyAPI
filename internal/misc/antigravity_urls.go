// Package misc - antigravity base URL helpers.
package misc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Default antigravity base URLs. Used when apiUrl.json is missing or a field is empty.
const (
	DefaultAntigravityBaseURLProd         = "https://cloudcode-pa.googleapis.com"
	DefaultAntigravityBaseURLDaily        = "https://daily-cloudcode-pa.googleapis.com"
	DefaultAntigravitySandboxBaseURLDaily = "https://daily-cloudcode-pa.sandbox.googleapis.com"
)

// AntigravityBaseURLConfigFileName is the filename probed alongside the executable
// (and in the working directory) to override the default base URLs at runtime.
const AntigravityBaseURLConfigFileName = "apiUrl.json"

type antigravityBaseURLConfig struct {
	Prod         string `json:"prod"`
	Daily        string `json:"daily"`
	SandboxDaily string `json:"sandbox_daily"`
}

var (
	antigravityURLOnce sync.Once
	antigravityURLs    []string
)

// AntigravityBaseURLs returns the ordered list of antigravity base URLs:
// prod, daily, sandbox_daily. Values from apiUrl.json (located next to the
// executable, with the working directory as a fallback) override the defaults;
// missing or blank fields fall back to the defaults. The lookup happens once
// per process.
func AntigravityBaseURLs() []string {
	antigravityURLOnce.Do(loadAntigravityBaseURLs)
	out := make([]string, len(antigravityURLs))
	copy(out, antigravityURLs)
	return out
}

func loadAntigravityBaseURLs() {
	cfg := antigravityBaseURLConfig{
		Prod:         DefaultAntigravityBaseURLProd,
		Daily:        DefaultAntigravityBaseURLDaily,
		SandboxDaily: DefaultAntigravitySandboxBaseURLDaily,
	}

	for _, dir := range antigravityBaseURLCandidateDirs() {
		if dir == "" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, AntigravityBaseURLConfigFileName))
		if err != nil {
			continue
		}
		var fileCfg antigravityBaseURLConfig
		if err := json.Unmarshal(data, &fileCfg); err != nil {
			continue
		}
		if v := strings.TrimSpace(fileCfg.Prod); v != "" {
			cfg.Prod = v
		}
		if v := strings.TrimSpace(fileCfg.Daily); v != "" {
			cfg.Daily = v
		}
		if v := strings.TrimSpace(fileCfg.SandboxDaily); v != "" {
			cfg.SandboxDaily = v
		}
		break
	}

	antigravityURLs = []string{cfg.Prod, cfg.Daily, cfg.SandboxDaily}
}

func antigravityBaseURLCandidateDirs() []string {
	dirs := make([]string, 0, 2)
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Dir(exe))
	}
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, wd)
	}
	return dirs
}

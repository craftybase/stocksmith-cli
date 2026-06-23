// Package selfupdate replaces the running binary with the latest GitHub release.
// All host/OS/network inputs are injected via Config so the logic is testable.
package selfupdate

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
	defaultAPIBaseURL      = "https://api.github.com"
	defaultDownloadBaseURL = "https://github.com"
)

// Config carries everything the updater needs.
type Config struct {
	BinaryName      string
	Repo            string
	CurrentVersion  string
	GOOS            string
	GOARCH          string
	ExecPath        string
	APIBaseURL      string // default https://api.github.com
	DownloadBaseURL string // default https://github.com
	HTTPClient      *http.Client
	Out             io.Writer
}

func (c *Config) apiBase() string {
	if c.APIBaseURL != "" {
		return c.APIBaseURL
	}
	return defaultAPIBaseURL
}

func (c *Config) downloadBase() string {
	if c.DownloadBaseURL != "" {
		return c.DownloadBaseURL
	}
	return defaultDownloadBaseURL
}

func (c *Config) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func (c *Config) out() io.Writer {
	if c.Out != nil {
		return c.Out
	}
	return os.Stdout
}

// assetNames returns the release archive and checksums filenames for version,
// matching .goreleaser.yml's name_template.
func (c *Config) assetNames(version string) (archive, checksums string) {
	vnum := strings.TrimPrefix(version, "v")
	archive = fmt.Sprintf("%s_%s_%s_%s.tar.gz", c.BinaryName, vnum, c.GOOS, c.GOARCH)
	checksums = fmt.Sprintf("%s_%s_checksums.txt", c.BinaryName, vnum)
	return archive, checksums
}

// canonicalVersion returns the semver-comparable form ("vX.Y.Z"), adding the
// leading v if absent.
func canonicalVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return v
}

// isDevVersion reports whether v cannot be compared against released tags.
func isDevVersion(v string) bool {
	return !semver.IsValid(canonicalVersion(v))
}

// updateAvailable reports whether latest is strictly newer than current.
func updateAvailable(current, latest string) bool {
	return semver.Compare(canonicalVersion(latest), canonicalVersion(current)) > 0
}

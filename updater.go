package update

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-version"
)

type Updater struct {
	AppName    string
	AppVersion string
	Author     string
}

func NewUpdater(appName string, appVersion string, author string) *Updater {
	return &Updater{
		AppName:    appName,
		AppVersion: appVersion,
		Author:     author,
	}
}

type Failure struct {
	Message string `json:"message"`
}

type Success struct {
	TagName string `json:"tag_name"`
}

func (u *Updater) fetchLatestVersion() (*version.Version, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.Author, u.AppName)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the response is a rate limit error
	if resp.StatusCode != http.StatusOK {
		var result Failure
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		return nil, errors.New(result.Message)
	}

	var result Success
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	latestVersion, err := version.NewVersion(result.TagName)
	if err != nil {
		return nil, err
	}

	return latestVersion, nil
}

func (u *Updater) CheckIfNewVersionIsAvailable() error {
	if u.AppVersion == "" {
		return nil
	}

	currentVersion, err := version.NewVersion(u.AppVersion)
	if err != nil {
		return err
	}

	latestVersion, err := u.fetchLatestVersion()
	if err != nil {
		return err
	}

	if latestVersion.GreaterThan(currentVersion) {
		fmt.Printf("A new version of %s is available (%s). Run `%s update` to update.\n\n", u.AppName, latestVersion, u.AppName)
	}

	return nil
}

// Determine the install path for the application
func (u *Updater) determineInstallPath() string {
	success, err := exec.Command("which", u.AppName).Output()
	if err == nil {
		return filepath.Dir(string(success))
	}

	if os.Getenv("GOBIN") != "" {
		return os.Getenv("GOBIN")
	}

	return "/usr/local/bin"
}

func (u *Updater) Update() error {
	version, err := u.fetchLatestVersion()
	if err != nil {
		return err
	}

	ldflags := fmt.Sprintf("-w -s -X main.AppVersion=%s -X main.AppName=%s", version, u.AppName)
	origin := fmt.Sprintf("github.com/%s/%s@%s", u.Author, u.AppName, version)
	args := []string{"install", "-ldflags", ldflags, origin}
	cmd := exec.Command("go", args...)

	// Use the GOBIN env var to tell go where to install the application
	path := u.determineInstallPath()
	cmd.Env = append(os.Environ(), "GOBIN="+path)

	// Capture the output of the command
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

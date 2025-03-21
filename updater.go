package updater

import (
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

// Returns the latest version if a new version is available
func (u *Updater) IsNewVersionAvailable() *version.Version {
	if u.AppVersion == "" {
		return nil
	}

	currentVersion, err := version.NewVersion(u.AppVersion)
	if err != nil {
		return nil
	}

	latestVersion, err := u.fetchLatestVersion()
	if err != nil {
		return nil
	}

	if latestVersion.GreaterThan(currentVersion) {
		return latestVersion
	}

	return nil
}

// Check if a new version is available and print a message
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
		fmt.Printf("A new version of %s is available (%s).\n\n", u.AppName, latestVersion)
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

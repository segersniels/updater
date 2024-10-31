# updater

A simple Go library for updating your Go applications.

## Installation

```bash
go get github.com/segersniels/updater
```

## Usage

```go
package main

import (
 "fmt"

 "github.com/segersniels/updater"
)

var (
 AppName    string
 AppVersion string
)

func main() {
  updater := update.NewUpdater(AppName, AppVersion, "segersniels")
  err := updater.CheckIfNewVersionIsAvailable()
  if err != nil {
    log.Debug("Failed to check for latest release")
  }

  err := updater.Update()
  if err != nil {
    log.Fatal("Failed to update application")
  }
}
```

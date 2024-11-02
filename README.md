# updater

A simple library for updating your Go applications.

## Installation

```bash
go get github.com/segersniels/updater
```

## Usage

```go
package main

import (
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
    println("Failed to check for latest release")
  }

  err := updater.Update()
  if err != nil {
    println("Failed to update application")
  }
}
```

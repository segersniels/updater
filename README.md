# updater

A simple library for updating your Go applications.

## Installation

```bash
go get -u github.com/segersniels/updater
```

## Usage

```go
package main

import (
 updater "github.com/segersniels/updater"
)

var (
 AppName    string
 AppVersion string
)

func main() {
  upd := updater.NewUpdater(AppName, AppVersion, "segersniels")
  err := upd.CheckIfNewVersionIsAvailable()
  if err != nil {
    println("Failed to check for latest release")
  }
}
```

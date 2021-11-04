# go-devcon

Go wrapper around the [Windows Device Console](https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon) (`devcon.exe`).

```
go install github.com/mikerourke/go-devcon
```

## Introduction

Here's a brief overview of DevCon taken from the Windows Hardware Developer documentation site:

> DevCon (Devcon.exe), the Device Console, is a command-line tool that displays detailed information about devices on computers running Windows. 
> You can use DevCon to enable, disable, install, configure, and remove devices.
> 
> DevCon runs on Microsoft Windows 2000 and later versions of Windows.

This package provides a handy mechanism for querying devices and performing various operations on said devices.
It provides methods that map directly to the [`devcon.exe` commands](https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-general-commands).

## Prerequisites

You'll need the `devcon.exe` executable. It's not included in the package because _Redistribution of Microsoft software without consent is usually a violation of Microsoft's End User License Agreements_.

For information regarding how to get your hands on `devcon.exe`, check out the [Microsoft documentation page](https://docs.microsoft.com/en-us/archive/blogs/deploymentguys/where-to-find-devcon-exe).

There are [other ways to get it](https://networchestration.wordpress.com/2016/07/11/how-to-obtain-device-console-utility-devcon-exe-without-downloading-and-installing-the-entire-windows-driver-kit-100-working-method/), but you didn't hear it from me.

## Usage

There is extensive documentation and examples available on the [docs site](https://pkg.go.dev/github.com/mikerourke/go-devcon), but here's a quick example of how to log all the devices on the computer:

```go
package main

import (
	"fmt"
	
	"github.com/mikerourke/go-devcon"
)

func main() {
	dc := devcon.New(`\path\to\devcon.exe`)
	
	devs, err := dc.FindAll("=net")
	if err != nil {
		panic(err)
    }
	
	for _, dev := range devs {
		fmt.Printf("ID: %s, Name: %s\n", dev.ID, dev.Name)
    }
}
```

## FAQ

#### Aren't you supposed to use the [`PnPUtil`](https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/pnputil) now?

Yes, but `PnPUtil` is only supported on Windows Vista and later. You can't use it on older operating systems.

#### Isn't this a little niche?

Yes, but one of my other projects requires me to check for the existence of a device and install a driver on Windows XP.
I wrote a bunch of parsing code to get the output, so I figured why not open source it?

#### Will this work on Windows XP?

Yes, I'm using only using Go APIs available in version `1.10.7` (or earlier?), since executables built with that version of Go still run on Windows XP.
It's possible that earlier versions _may_ work, but I used `1.10.7` during development, so I'd advise sticking with that version.
To target Windows XP 32-bit, build your project like so:

```
GOOS=windows GOARCH=386 go1.10.7 build -o myproject.exe myproject.go  
```
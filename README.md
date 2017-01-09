[![Build Status](https://travis-ci.org/hreese/resolveip.svg?branch=master)](https://travis-ci.org/hreese/resolveip)

# resolveip

```resolveip``` finds IP addresses (both IPv4 and IPv6 by default) in texts and resolves them using the system's local resolver.
It is primarily meant as an interactive tool but can also read from files and pipes.

![screencast](res/.screencast01.gif)

When reading files (this includes dropping text files on the resolveip icon) on Windows, resolveip adds the console input to the list of inputs. This keeps the console window open after reading all initial inputs. To use resolveip in pipes on Windows, add the ```-batch``` flag.

## Binary releases

### Linux
* [x86 (32bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_linux_386.tar.bz2)
* **[x86-64 (64bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_linux_amd64.tar.bz2)**
* [ARM](https://stuff.heiko-reese.de/resolveip/resolveip_linux_arm.tar.bz2)
* [ARM64](https://stuff.heiko-reese.de/resolveip/resolveip_linux_arm64.tar.bz2)

### Windows
* [x86 (32bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_windows_386.zip)
* **[x86-64 (64bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_windows_amd64.zip)**

### macOS

* [x86 (32bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_darwin_386.dmg)
* **[x86-64 (64bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_darwin_amd64.dmg)**

Open the DMG, copy resolveip into your ```$PATH``` and remove the quarantine attribute if necessary:
```
xattr -rd com.apple.quarantine resolveip
```

## Building from source

### Get Go

[Install](https://golang.org/dl) and [configure](https://golang.org/doc/install) the [Go](https://golang.org/) toolchain. Most Linux distributions already have it packaged. [Homebrew](http://brew.sh) also has a package.

### Get the source

```
go get github.com/hreese/resolveip
# install this if you want to compile for Windows
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo
```
### Build

#### Build for your local architecture
```
go install github.com/hreese/resolveip/cmd/resolveip
```

#### Build and package for all supported architectures

```
cd $GOPATH/src/github.com/hreese/resolveip
make
```

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
Compiling macOS binaries on Linux does not seem to yield working software at the moment. Please build it yourself or wait for me to get a usable build machine.

* ~~[x86 (32bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_darwin_386.tar.bz2)~~
* ~~**[x86-64 (64bit Intel/AMD)](https://stuff.heiko-reese.de/resolveip/resolveip_darwin_amd64.tar.bz2)**~~

## Building from source
1. [Install](https://golang.org/dl) and [configure](https://golang.org/doc/install) the [Go](https://golang.org/) toolchain. Most Linux distributions already have it packaged. [Homebrew](http://brew.sh) also has a package.
2. Get the source:
```
go get git.heiko-reese.de/hreese/resolveip
# install this if you want to compile for Windows
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo
```
3. Build:

Build for your local architecture:
```
go install git.heiko-reese.de/hreese/resolveip/cmd/resolveip
```

Build and package for all supported architectures:

```
cd $GOPATH/src/git.heiko-reese.de/hreese/resolveip
make
```

## ToDo

* better documentation
* add more commandline switches for output customization:
    * different ANSI codes for text, results, matches and non-matches
    * remove resolved ip addresses

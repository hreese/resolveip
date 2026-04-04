[![CI and Release](https://github.com/hreese/resolveip/actions/workflows/release.yml/badge.svg)](https://github.com/hreese/resolveip/actions/workflows/release.yml)

# resolveip

```resolveip``` finds IP addresses (both IPv4 and IPv6 by default) in texts and resolves them using the system's local resolver.
It is primarily meant as an interactive tool but can also read from files and pipes.

![screencast](res/.screencast01.gif)

When reading files (this includes dropping text files on the resolveip icon) on Windows, resolveip adds the console input to the list of inputs. This keeps the console window open after reading all initial inputs. To use resolveip in pipes on Windows, add the ```-batch``` flag.

## Binary releases

Binaries for Linux, macOS (Intel and Apple Silicon) and Windows are automatically built by GitHub Actions and available
on the [project's releases page](https://github.com/hreese/resolveip/releases).

Note to macOS users: Extract the archive, copy resolveip into your ```$PATH``` and remove the
quarantine attribute if necessary:
```
xattr -rd com.apple.quarantine resolveip
```

## Building from source

### Get Go

[Install](https://golang.org/dl) and [configure](https://golang.org/doc/install) the [Go](https://golang.org/) toolchain. Most Linux distributions already have it packaged. [Homebrew](http://brew.sh) also has a package.

### Get the source

```
git clone https://github.com/hreese/resolveip.git
cd resolveip
# install this if you want to compile for Windows
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
```
### Build

#### Build for your local architecture
```
go install github.com/hreese/resolveip/v2/cmd/resolveip@latest
```

#### Build and package for all supported architectures

Cross-compilation and release packaging is handled by [GoReleaser](https://goreleaser.com/).
To build a snapshot locally:

```
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser build --snapshot --clean
```

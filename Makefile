packagename := resolveip
source := github.com/hreese/resolveip/cmd/resolveip
builddir := build
ldflags := "-s -w"
PLATFORMS_UNIX := linux/386/tar/ linux/amd64/tar/ linux/arm/tar/ linux/arm64/tar/ darwin/386/dmg/ darwin/amd64/dmg/
PLATFORMS_CURRENTLY_UNSUPPORTED := solaris/amd64/tar/ freebsd/386/tar/ freebsd/amd64/tar/
PLATFORMS_WIN  := windows/386/zip/.exe windows/amd64/zip/.exe
PLATFORMS := $(PLATFORMS_UNIX) $(PLATFORMS_WIN)

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
packer = $(word 3, $(temp))
ext = $(word 4, $(temp))

tar = cd '$(builddir)/$(os)-$(arch)/' && tar cjf '../$(packagename)_$(os)_$(arch).tar.bz2' * && cd .. && rm -rf '$(os)-$(arch)/'
zip = cd '$(builddir)/$(os)-$(arch)/' && zip -9  '../$(packagename)_$(os)_$(arch).zip' * && cd .. && rm -rf '$(os)-$(arch)/'
dmg = cd '$(builddir)/$(os)-$(arch)/' && genisoimage -V '$(packagename)' -D -R -apple -no-pad -o ../$(packagename)_$(os)_$(arch).dmg * && cd .. && rm -rf '$(os)-$(arch)/'

release: $(PLATFORMS)

$(PLATFORMS): README.md
	GOOS=$(os) GOARCH=$(arch) go build -ldflags=$(ldflags) -o '$(builddir)/$(os)-$(arch)/$(packagename)$(ext)' $(source)
	cd '$(builddir)/$(os)-$(arch)/' && sha256sum -b * > sha256sum.txt
	sed -r -e '/.screencast01.gif/d' README.md > '$(builddir)/$(os)-$(arch)/README.md'
	$(call $(packer))

deploy: $(PLATFORMS)
	rsync -vaP build/* deploy_binary_reolveip:/srv/www/stuff.heiko-reese.de/resolveip/

clean:
	rm -rf $(builddir)

.PHONY: release clean $(PLATFORMS)

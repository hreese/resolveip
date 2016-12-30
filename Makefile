packagename := resolveip
source := cmd/resolveip.go
ldflags := "-s -w"
PLATFORMS := linux/386/tar/ linux/amd64/tar/ linux/arm/tar/ linux/arm64/tar/ freebsd/386/tar/ freebsd/amd64/tar/ darwin/386/tar/ darwin/amd64/tar/ windows/386/zip/.exe windows/amd64/zip/.exe solaris/amd64/tar/

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
packer = $(word 3, $(temp))
ext = $(word 4, $(temp))

tar = cd 'build/$(os)-$(arch)/' && tar cjf '../$(packagename)_$(os)_$(arch).tar.bz2' * && cd .. && rm -rf '$(os)-$(arch)/'
zip = cd 'build/$(os)-$(arch)/' && zip -9  '../$(packagename)_$(os)_$(arch).zip' * && cd .. && rm -rf '$(os)-$(arch)/'

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -ldflags=$(ldflags) -o 'build/$(os)-$(arch)/$(packagename)$(ext)' $(source)
	cd 'build/$(os)-$(arch)/' && sha256sum -b * > sha256sum.txt
	$(call $(packer))

deploy: $(PLATFORMS)
	rsync -vaP build/* deploy_binary_reolveip:/srv/www/stuff.heiko-reese.de/resolveip/

clean:
	rm -rf build

.PHONY: release clean $(PLATFORMS)

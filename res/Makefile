.PHONY: default clean

resolutions := 16 32 48 256
sourceicon := icon.svg
iconfiles := $(addsuffix .png,$(resolutions))

default: icon.ico

clean:
	rm -f  $(iconfiles) icon.ico

%.ico: $(iconfiles)
	convert $? $@

%.png: $(sourceicon)
	inkscape --without-gui --export-width=$* --export-height=$* --export-png=$@ $<

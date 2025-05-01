.PHONY: all dev build clean-i18n
YQ=yq
TAGS=-tags webkit2_41


all: all-i18n dev


dev:
	wails dev $(TAGS)

build:
	wails build -clean -u -upx --trimpath $(TAGS)

run: build
	build/bin/PolAIn

all-i18n: $(patsubst locales/%.yaml,i18n/%.json,$(wildcard locales/*.yaml))

i18n/%.json: locales/%.yaml
	$(YQ) $< -M -o json > $@

clean-i18n:
	rm -f i18n/*.json

YQ=yq


all: all-i18n dev


dev:
	wails dev -tags webkit2_41

all-i18n: $(patsubst locales/%.yaml,i18n/%.json,$(wildcard locales/*.yaml))

i18n/%.json: locales/%.yaml
	$(YQ) $< -M -o json > $@

clean-i18n:
	rm -f i18n/*.json

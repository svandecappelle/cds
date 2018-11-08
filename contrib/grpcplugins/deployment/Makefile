.PHONY: clean

VERSION := $(if ${CDS_SEMVER},${CDS_SEMVER},snapshot${CDS_VERSION})
GITHASH := $(if ${GIT_HASH},${GIT_HASH},`git log -1 --format="%H"`)
BUILDTIME := `date "+%m/%d/%y-%H:%M:%S"`

TARGET_DIR = dist
TARGET_NAME = $(filter-out $@ examples,$(MAKECMDGOALS))

define PLUGIN_MANIFEST_BINARY
os: %os%
arch: %arch%
cmd: ./%filename%
endef
export PLUGIN_MANIFEST_BINARY

TARGET_LDFLAGS = -ldflags "-X github.com/ovh/cds/sdk.VERSION=$(VERSION) -X github.com/ovh/cds/sdk.GOOS=$$GOOS -X github.com/ovh/cds/sdk.GOARCH=$$GOARCH -X github.com/ovh/cds/sdk.GITHASH=$(GITHASH) -X github.com/ovh/cds/sdk.BUILDTIME=$(BUILDTIME) -X github.com/ovh/cds/sdk.BINARY=$(TARGET_NAME)"
TARGET_OS = $(if ${OS},${OS},windows darwin linux freebsd openbsd)
TARGET_ARCH = $(if ${ARCH},${ARCH},amd64 arm 386)

GO_BUILD = go build -v
DIRS = `ls -d */ | cut -f1 -d'/'`

default: build-all

clean:
	@rm -rf */$(TARGET_DIR)

build-all:
	@for plugin in $(DIRS); do \
		pluginname=$(basename "$$plugin"); \
		mkdir -p $$plugin/$(TARGET_DIR); \
		cp $$plugin/plugin.yml $$plugin/$(TARGET_DIR)/plugin.yml; \
		for GOOS in $(TARGET_OS); do \
			for GOARCH in $(TARGET_ARCH); do \
				EXTENSION=""; \
				if test "$$GOOS" = "windows" ; then EXTENSION=".exe"; fi; \
				FILENAME=deployment-plugin-$$pluginname-$$GOOS-$$GOARCH$$EXTENSION; \
				echo Compiling $$plugin/$(TARGET_DIR)/$$FILENAME $(VERSION); \
				GOOS=$$GOOS GOARCH=$$GOARCH $(GO_BUILD) $(TARGET_LDFLAGS) -o $$plugin/$(TARGET_DIR)/$$FILENAME ./$$pluginname; \
				echo "$$PLUGIN_MANIFEST_BINARY" > $$plugin/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
				perl -pi -e s,%os%,$$GOOS,g $$plugin/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
				perl -pi -e s,%arch%,$$GOARCH,g $$plugin/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
				perl -pi -e s,%filename%,$$FILENAME,g $$plugin/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
			done; \
		done; \
	done

publish-all:
	@echo "Updating plugins..."
	@for plugin in $(DIRS); do \
		pluginname=$(basename "$$plugin"); \
		cdsctl admin plugins import $$plugin/$(TARGET_DIR)/plugin.yml; \
		for GOOS in $(TARGET_OS); do \
			for GOARCH in $(TARGET_ARCH); do \
			EXTENSION=""; \
			if test "$$GOOS" = "windows" ; then EXTENSION=".exe"; fi; \
			FILENAME=deployment-plugin-$$pluginname-$$GOOS-$$GOARCH$$EXTENSION; \
			PLUGIN_YML=$$plugin/$(TARGET_DIR)/plugin.yml; \
			NAME_PLUGIN=`grep 'name:' $$PLUGIN_YML|cut -d ':' -f2;`; \
			echo "Updating plugin $$NAME_PLUGIN with binary $$FILENAME"; \
			cdsctl admin plugins binary-add $$NAME_PLUGIN $$plugin/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml $$plugin/$(TARGET_DIR)/$$FILENAME; \
			done; \
		done; \
	done;

build:
	@mkdir -p $(TARGET_NAME)/$(TARGET_DIR)
	@cp $(TARGET_NAME)/plugin.yml $(TARGET_NAME)/$(TARGET_DIR)/plugin.yml
	@for GOOS in $(TARGET_OS); do \
		for GOARCH in $(TARGET_ARCH); do \
			EXTENSION=""; \
			if test "$$GOOS" = "windows" ; then EXTENSION=".exe"; fi; \
			FILENAME=deployment-$(TARGET_NAME)-$$GOOS-$$GOARCH$$EXTENSION; \
			echo Compiling $(TARGET_DIR)/$$FILENAME $(VERSION); \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GO_BUILD) $(TARGET_LDFLAGS) -o $(TARGET_NAME)/$(TARGET_DIR)/$$FILENAME ./$(TARGET_NAME); \
			echo "$$PLUGIN_MANIFEST_BINARY" > $(TARGET_NAME)/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
			perl -pi -e s,%os%,$$GOOS,g $(TARGET_NAME)/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
			perl -pi -e s,%arch%,$$GOARCH,g $(TARGET_NAME)/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
			perl -pi -e s,%filename%,$$FILENAME,g $(TARGET_NAME)/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml; \
		done; \
	done

publish:
	@echo "Updating plugin $(TARGET_NAME)..."; \
	cdsctl admin plugins import $(TARGET_NAME)/$(TARGET_DIR)/plugin.yml;
	@for GOOS in $(TARGET_OS); do \
		for GOARCH in $(TARGET_ARCH); do \
			EXTENSION=""; \
			if test "$$GOOS" = "windows" ; then EXTENSION=".exe"; fi; \
			FILENAME=deployment-$(TARGET_NAME)-$$GOOS-$$GOARCH$$EXTENSION; \
			PLUGIN_YML=$(TARGET_NAME)/$(TARGET_DIR)/plugin.yml; \
			NAME_PLUGIN=`grep 'name:' $$PLUGIN_YML|cut -d ':' -f2;`; \
			echo "Updating plugin $$NAME_PLUGIN with binary $$FILENAME"; \
			cdsctl admin plugins binary-add $$NAME_PLUGIN $(TARGET_NAME)/$(TARGET_DIR)/plugin-$$GOOS-$$GOARCH.yml $(TARGET_NAME)/$(TARGET_DIR)/$$FILENAME; \
		done; \
	done

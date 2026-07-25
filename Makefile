BIN=mail-alias-resolve-ldap
HASH   := $(shell git rev-parse --short HEAD)
REMOTE ?= gh
RELEASE_BRANCH ?= dev

.PHONY: all build strip clean release release-candidate rc _check-remote _check-up-to-date

all: build strip

build:
	go build -o $(BIN)

strip:
	strip $(BIN)

clean:
	rm -f $(BIN)

_check-remote:
	@git remote get-url $(REMOTE) > /dev/null 2>&1 || \
	    { echo "Error: no remote '$(REMOTE)' — add one with: git remote add $(REMOTE) <url>"; exit 1; }

_check-up-to-date: _check-remote
	@git fetch $(REMOTE) $(RELEASE_BRANCH) > /dev/null 2>&1
	@git merge-base --is-ancestor $(REMOTE)/$(RELEASE_BRANCH) HEAD || \
	    { echo "Error: $(RELEASE_BRANCH) has commits you don't have — pull/rebase before tagging a release."; exit 1; }

release: _check-up-to-date
	$(eval TAG := release-$(HASH))
	git tag -f $(TAG)
	@printf 'Tagged %s as %s\n' "$(HASH)" "$(TAG)"
	@printf 'Push tag to trigger a release? [y/N] ' && read ans && \
	    case "$$ans" in [yY]) git push $(REMOTE) $(TAG) ;; \
	    *) git tag -d $(TAG); echo 'Aborted — tag removed.' ;; esac

release-candidate rc: _check-up-to-date
	$(eval TAG := rc-$(HASH))
	git tag -f $(TAG)
	@printf 'Tagged %s as %s\n' "$(HASH)" "$(TAG)"
	@printf 'Push tag to trigger a release candidate? [y/N] ' && read ans && \
	    case "$$ans" in [yY]) git push $(REMOTE) $(TAG) ;; \
	    *) git tag -d $(TAG); echo 'Aborted — tag removed.' ;; esac

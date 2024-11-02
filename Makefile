VERSION := $(shell git log --grep="^chore(release):" -1 --format="%s" | sed 's/chore(release): //')

version:
	@echo $(VERSION)

.PHONY: get-version


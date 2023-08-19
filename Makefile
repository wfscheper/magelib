MAGE = bin/mage

.PHONY: default
default: all

.PHONY: help
help: | $(MAGE)
	@$(MAGE) -l

.PHONY: clean
clean:
	$(RM) -r bin/

$(MAGE): mage/go.mod mage/go.sum mage/mage.go mage/magefile.go ; $(info ▶ building $@)
	cd mage/ && go run mage.go -compile ../$@

mage/go.mod mage/go.sum mage/mage.go mage/magefile.go:
	@:

%: | $(MAGE)
	@$(MAGE) $(if $(V),-v) $*

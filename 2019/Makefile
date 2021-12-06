DAYS=$(shell find . -mindepth 1 -type d  | sed -e 's/\.\///g')
.PHONY: ${DAYS}

all: ${DAYS}

${DAYS}:
	cat ${@}/input.txt | go run ${@}/main.go

fmt:
	for f in $$(find -name '*.go'); do go fmt $$f; done

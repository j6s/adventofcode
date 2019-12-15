DAYS=$(shell ls 2019)

all: ${DAYS}
${DAYS}:
	cat 2019/${@}/input.txt | go run 2019/${@}/main.go
	@echo -e "\n"

fmt:
	for f in $$(find -name '*.go'); do go fmt $$f; done
all: run

run:
	go build -o webrdd cmd/webrdd/* && ./webrdd

reformat:
	find . \( -name '*.m' -o -name '*.h' \) -exec clang-format -i --style="{BasedOnStyle: llvm, IndentWidth: 4}" '{}' \;

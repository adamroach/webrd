all: run

run:
	go build main.go && ./main

reformat:
	find . \( -name '*.m' -o -name '*.h' \) -exec clang-format -i --style="{BasedOnStyle: llvm, IndentWidth: 4}" '{}' \;
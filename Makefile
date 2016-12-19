.PHONY: all
all: ga

.PHONY: clean
clean:
	rm -f ga

ga: ga.go helper.go datastructures.go
	go build $^

.PHONY: test
test: ga
	./ga

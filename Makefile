.PHONY: all

all: ssmt-ssu

ssmt-ssu: main.go
	go build -o $@ $^

clean: 
	rm -f ssmt-ssu
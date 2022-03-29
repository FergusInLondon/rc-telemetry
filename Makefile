
clean:
	rm -f parser

example: clean
	go build -o parser ./examples/ltm/main.go

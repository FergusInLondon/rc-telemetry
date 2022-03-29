
clean:
	rm -f parser

example: clean
	go build -o parser ./example/main.go

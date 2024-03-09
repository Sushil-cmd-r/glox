run: build
	./glox

build: clean
	@go build -o glox

test:
	@go test ./...

clean:
	@rm -f glox

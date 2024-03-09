run: build
	./glox

build: clean
	@go build -o glox

clean:
	@rm -f glox

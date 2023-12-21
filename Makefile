build: 
	@go build -o bin/rssagg
run: build
	@./bin/rssagg
test: 
	@go test -v ./...

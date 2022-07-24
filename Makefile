test:
	go clean -testcache
	go test -v --coverprofile coverage.out ./...
	go tool cover -html=coverage.out

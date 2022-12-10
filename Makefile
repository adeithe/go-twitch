test:
	go clean -testcache
	go test -timeout 300s -v --coverprofile coverage.out ./...
	go tool cover -html=coverage.out

bench:
	go clean -testcache
	go test -run=Benchmark -bench=. ./...

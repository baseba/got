run:
	@templ generate
	@go run cmd/main.go
run_nodocker:
	@~/go/bin/templ generate
	@go run cmd/main.go
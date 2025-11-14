build-ui:
	cd ui && bun run build && cd ..

dev: build-ui
	go run dgstuff.go
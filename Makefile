run:
	go run -race ./cmd/switch/*.go

test:
	go test -v --cover ./...

release:
	@echo "Enter the release version (format vx.x.x).."; \
	read VERSION; \
	git tag -a $$VERSION -m "Releasing "$$VERSION; \
	git push origin $$VERSION
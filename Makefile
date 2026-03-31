.PHONY: run-gin run-fiber run-stdlib run-echo seed

run-gin:
	go run gin/cmd/main.go

run-fiber:
	go run fiber/cmd/main.go

run-stdlib:
	go run stdlib/cmd/main.go

run-echo:
	go run echo/cmd/main.go

seed:
	go run scripts/seed.go

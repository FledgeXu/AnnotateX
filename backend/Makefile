.DEFAULT_GOAL := build

build: gen_mock gen_wire static_check
	go build -o build/app

run: static_check gen_wire
	go run main.go

gen_mock:
	mockery

static_check:
	staticcheck ./...

gen_wire:
	wire wire/wire.go

clear:
	rm wire/wire_gen.go
	rm mocks/*

.PHONY: run gen_mock staticcheck wire clear

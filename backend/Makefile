.DEFAULT_GOAL := run

run: gen_wire
	go run main.go

gen_wire:
	wire wire/wire.go

clear:
	rm wire/wire_gen.go
	rm mocks/*

:xa
.PHONY: run wire clear

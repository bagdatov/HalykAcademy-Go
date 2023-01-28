.PHONY: build

build:
	go build

.PHONY: test

test:
	go test -v api/*

# only for linux:
.PHONY: bench-profile

bench-profile:
	go build
	@sleep 1
	gnome-terminal --geometry=80x40+250+200 -e "./web ;bash"
	@sleep 1
	gnome-terminal --geometry=80x40+1650+200 -e "go tool pprof web http://localhost:8080/debug/pprof/profile?seconds=10;bash"
	gnome-terminal --geometry=80x40+650+200 -e "bash request.sh ;bash"

# only for linux:
.PHONY: bench-heap

bench-heap:
	go build
	@sleep 1
	gnome-terminal --geometry=80x40+250+200 -e "./web ;bash"
	@sleep 1
	gnome-terminal --geometry=80x40+1650+200 -e "go tool pprof -alloc_space web http://localhost:8080/debug/pprof/heap?seconds=10;bash"
	gnome-terminal --geometry=80x40+650+200 -e "bash request.sh ;bash"

.DEFAULT_GOAL := build
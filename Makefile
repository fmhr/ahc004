build: main.go
	go build


example:
	go run main.go < example.in > example.out

vis:
	cd tools &&\
	cargo run --release --bin vis ../example.in ../example.out


runscript:
	make build
	go run script/main.go


pprof:
	.solver -cpuprofile cpu.prof < tools/in/0000.txt
	pprof -http=localhost:8080 cpu.prof

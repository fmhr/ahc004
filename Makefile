build: main.go
	go build


example:
	go run main.go < example.in > example.out

vis:
	cd tools &&\
	cargo run --release --bin vis ../example.in ../example.out


runscript:
	go build
	go run script/main.go


pprof:
	go build
	./solver -cpuprofile cpu.prof < tools/in/0000.txt
	pprof -http=localhost:8080 cpu.prof

setuptools:
	curl -O https://img.atcoder.jp/ahc004/222362f13a30b1342bf79d0041bd4d39.zip
	unzip 222362f13a30b1342bf79d0041bd4d39.zip
	rm 222362f13a30b1342bf79d0041bd4d39.zip
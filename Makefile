build_all: main.go
	go build -o Bin/out main.go
	cp -r Datasets/ Bin/ 

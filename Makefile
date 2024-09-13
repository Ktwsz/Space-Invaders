linux:
	GOOS=linux go build -C src -o build/main

windows:
	GOOS=windows go build -C src -o build/main

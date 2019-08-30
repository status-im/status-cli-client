all:
	go build -o status-cli-client main.go

clean:
	rm -f ./status-cli-client

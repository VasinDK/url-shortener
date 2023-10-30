CGO_ENABLED=1
CONFIG_PATH=.\config\local.yaml

run: .\cmd\url-shortener\main.go
	go run .\cmd\url-shortener\main.go
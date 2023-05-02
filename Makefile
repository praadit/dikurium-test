gqlinit:
	cd pkg && gqlgen init

gqlgen:
	cd pkg && gqlgen generate

server:
	cd cmd && go run main.go

test:
	cd pkg/test && go test
FROM golang:1.18

WORKDIR /auth0-demo-api-go
COPY ./ ./

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["go", "run", "main.go"]

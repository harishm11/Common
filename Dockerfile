FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

COPY config/config.yaml /app/config.yaml

RUN go build -o /Common ./cmd/main.go

EXPOSE 3000

CMD [ "/Common" ]

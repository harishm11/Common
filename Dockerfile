FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /API-Gateway ./cmd/main.go

EXPOSE 3000

CMD [ "/API-Gateway" ]

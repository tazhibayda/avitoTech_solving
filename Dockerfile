FROM golang:1.18


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY pkg/handler/*.go pkg/database/db.go cmd/main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8808

CMD [ "/docker-gs-ping" ]
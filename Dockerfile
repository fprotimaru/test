FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /app/bin cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=build /app/bin /app/bin

EXPOSE 8080
EXPOSE 8081

CMD ["./bin"]
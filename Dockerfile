FROM golang:latest AS builder
WORKDIR /builder
# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


FROM alpine:latest
ENV GIN_MODE=release\
    PORT=8080
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /builder/app .
EXPOSE 8080
ENTRYPOINT ["./app"]

FROM golang:1.26 AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o redactlog .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/redactlog .
EXPOSE 8080
CMD ["./redactlog", "-serve"]
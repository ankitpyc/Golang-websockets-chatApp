FROM golang:1.21.6
WORKDIR /app
RUN mkdir server
RUN mkdir redis-cache
WORKDIR /app
COPY server/* /app/server/
COPY redis-cache/* /app/redis-cache/
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /webserver
CMD ["/webserver"]
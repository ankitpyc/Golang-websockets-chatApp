FROM golang:1.21.6
WORKDIR /app
COPY . /app
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /webserver
CMD ["/webserver"]
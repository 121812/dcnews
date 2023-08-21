FROM golang:1.20

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod download
RUN go build
RUN mv docker-config.json config.json

CMD ["./dcnews"]

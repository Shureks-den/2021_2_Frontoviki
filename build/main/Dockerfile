FROM golang:alpine AS run_step
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN apk add --update vips-dev
RUN apk add build-base
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o compressor cmd/compressor/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main/main.go
RUN cp /app/configs/start.sh /app/start.sh
ENTRYPOINT [ "./start.sh" ]
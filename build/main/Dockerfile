FROM golang:alpine AS build_step
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN apk add --update vips-dev
RUN apk add build-base
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o compressor cmd/compressor/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main/main.go

FROM alpine AS run_step
WORKDIR /app
RUN apk add --update vips-dev
COPY --from=build_step app/main .
COPY --from=build_step app/compressor .
ENTRYPOINT [ "./configs/start.sh" ]
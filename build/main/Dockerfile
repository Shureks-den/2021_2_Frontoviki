FROM golang:alpine AS build_step
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN apk add --update vips-dev
RUN apk add build-base
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o compressor cmd/compressor/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main/main.go
RUN mv .env_for_docker .env
RUN mv config_for_docker.yml config_dev.yml

FROM alpine AS run_step
WORKDIR /app
RUN apk add --update vips-dev
COPY --from=build_step app/static static
COPY --from=build_step app/main .
COPY --from=build_step app/compressor .
COPY --from=build_step app/.env .
COPY --from=build_step app/config_dev.yml .
COPY --from=build_step app/start.sh .
RUN cp /app/config_dev.yml /root/config_dev.yml
RUN chmod 777 ./start.sh
ENTRYPOINT [ "./start.sh" ]
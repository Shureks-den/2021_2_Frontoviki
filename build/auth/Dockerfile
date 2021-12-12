FROM golang:alpine AS build_step
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth cmd/auth/main.go
RUN mv .env_for_docker .env
RUN mv config_for_docker.yml config_dev.yml

FROM alpine AS run_step
WORKDIR /app
COPY --from=build_step app/auth .
COPY --from=build_step app/.env .
COPY --from=build_step app/config_dev.yml .
ENTRYPOINT [ "./auth" ]
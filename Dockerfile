FROM golang:1.16-alpine as builder
COPY . /app
WORKDIR /app
RUN go build -o out/gitlab-webhook-create

FROM alpine:3.14
COPY --from=builder /app/out/gitlab-webhook-create /app/gitlab-webhook-create

WORKDIR /app

CMD ["./gitlab-webhook-create"]
FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/telenotion-bot ./

FROM alpine:3.20

# 安装时区数据
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=build /out/telenotion-bot /app/telenotion-bot

ENTRYPOINT ["/app/telenotion-bot"]

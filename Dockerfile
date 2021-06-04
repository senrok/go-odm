
FROM golang:1.15-alpine as builder

# Add the keys
ARG token
ENV token=$token
ARG user
ENV user=$user

RUN apk update && apk add git

RUN git config \
    --global \
    url."https://${user}:${token}@github.com/".insteadOf \
    "https://github.com/"

WORKDIR /root
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN export GO111MODULE=on && CGO_ENABLED=0 GOOS=linux go build  -ldflags "-s -w" -o build/main cmd/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /root
COPY --from=builder /root/build/main ./

ENTRYPOINT ["/root/main"]

# Builder Image
FROM golang:1.20 as builder
ARG version

# create and set working directory
RUN mkdir -p /app
WORKDIR /app

# install dependencies
ADD go.mod go.mod
ADD go.sum go.sum

RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
RUN go mod download
# add code
ADD . .
# build the source
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=$version'" -o main ./cmd/main.go

# Final Image
FROM alpine:3.16
# add ca-certificates
RUN apk update && apk --no-cache  add ca-certificates
# set working directory
WORKDIR /app
# copy the binary from builder
COPY --from=builder /app/main ./main
# copy the timezone information
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# run the binary
EXPOSE 8080
EXPOSE 2112
CMD ["./main"]

FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN  go build -v -buildvcs=false github.com/lordralex/github-ip-to-opnsense

FROM alpine:latest

COPY --from=builder /build/github-ip-to-opnsense .

EXPOSE 8080
ENTRYPOINT ["./github-ip-to-opnsense"]
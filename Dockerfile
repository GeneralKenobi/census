FROM golang:1.18.0 as builder
WORKDIR /go/github.com/GeneralKenobi/census

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY internal internal
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o census ./cmd/census


FROM alpine:3.15
WORKDIR /opt/census
COPY --from=builder /go/github.com/GeneralKenobi/census/census census
RUN chmod 755 census

EXPOSE 8080
ENTRYPOINT ["/opt/census/census"]

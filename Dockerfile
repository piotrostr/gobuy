FROM golang:alpine AS builder

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main ./main.go


FROM scratch

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]

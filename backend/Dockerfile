FROM golang:1.25 AS builder
WORKDIR /usr/app/diplom
COPY ./ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .
RUN go build -o main .
CMD ["./main"]



# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /app

COPY ./services/log/go.mod ./
COPY ./services/log/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o log ./services/log/cmd

FROM alpine:3.18

COPY --from=stage1 /app/log .

EXPOSE 8080

CMD ["./log"]

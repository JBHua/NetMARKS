# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /app

COPY ./services/echo/go.mod ./
COPY ./services/echo/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o echo ./services/echo/cmd

FROM alpine:3.18

COPY --from=stage1 /app/echo .

EXPOSE 8080

CMD ["./echo"]

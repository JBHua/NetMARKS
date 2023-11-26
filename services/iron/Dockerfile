# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /app

COPY ./services/iron/go.mod ./
COPY ./services/iron/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o iron ./services/iron/cmd

FROM alpine:3.18

COPY --from=stage1 /app/iron .

EXPOSE 8080

CMD ["./iron"]

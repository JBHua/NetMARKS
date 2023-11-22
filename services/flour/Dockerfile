# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /app

COPY ./services/flour/go.mod ./
COPY ./services/flour/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o flour ./services/flour/cmd

FROM alpine:3.18

COPY --from=stage1 /app/flour .

EXPOSE 8080

CMD ["./flour"]

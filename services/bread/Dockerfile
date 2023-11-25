# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /app

COPY ./services/bread/go.mod ./
COPY ./services/bread/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o bread ./services/bread/cmd

FROM alpine:3.18

COPY --from=stage1 /app/bread .

EXPOSE 8080

CMD ["./bread"]

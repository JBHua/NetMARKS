# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /app

COPY ./services/grain/go.mod ./
COPY ./services/grain/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o grain ./services/grain/cmd

FROM alpine:3.18

COPY --from=stage1 /app/grain .

EXPOSE 8080

CMD ["./grain"]

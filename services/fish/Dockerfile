# Stage 1
FROM golang:1.21.1-alpine3.18 as stage1

# Set the working directory inside the container
WORKDIR /

COPY ./services/fish/go.mod ./
COPY ./services/fish/go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o fish ./services/fish/cmd

FROM alpine:3.18

COPY --from=stage1 /fish .

EXPOSE 8080

CMD ["./fish"]

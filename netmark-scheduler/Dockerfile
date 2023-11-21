#FROM golang:alpine AS builder
#WORKDIR /
#COPY . .
#
#RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o scheduler main.go
#
## Final Stage
#FROM alpine AS runner
#WORKDIR /
#
#COPY --from=builder /scheduler /
#
#CMD ["/scheduler"]


FROM golang:alpine AS builder

WORKDIR /scheduler
COPY ./ /scheduler

RUN \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o scheduler main.go

FROM alpine AS runner
WORKDIR /go/scheduler

COPY --from=builder /scheduler/scheduler .
COPY --from=builder /scheduler/infra/scheduler.yaml /etc/

VOLUME ["./scheduler"]

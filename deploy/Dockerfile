FROM golang:1.18-alpine as build

WORKDIR /us-build

COPY ./config ./config
COPY ./controller ./controller
COPY ./model ./model
COPY ./service ./service
COPY main.go go.mod go.sum ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -o usmain main.go

## main container
FROM alpine

COPY --from=build /us-build/usmain .
ENTRYPOINT [ "./usmain" ]
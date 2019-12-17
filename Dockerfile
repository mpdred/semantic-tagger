FROM golang:alpine as build

WORKDIR /app/src
COPY go.mod .
RUN go mod download

COPY main.go .
COPY pkg/ ./pkg/

RUN mkdir /app/out \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o /app/out/semtag

###

FROM scratch
COPY --from=build /app/out/semtag /semtag
CMD "foo"

## We specify image we need for our go application build
FROM golang:1.15.3-alpine3.12 AS build

## "go get" command requires git
RUN apk add git

RUN mkdir /app

# copy source only
COPY main.go /app

WORKDIR /app

## download dendencies
RUN go get github.com/gorilla/mux

## we run go build to compile the binary
## executable of our Go program
RUN go build -o main .

# runtime image
FROM alpine:3.12
# bring exe from build image step
COPY --from=build /app/main /bin/threatcloudapi
# main service
ENTRYPOINT ["/bin/threatcloudapi"]
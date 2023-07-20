FROM golang:1.20-alpine3.17
# Note that if you use this Dockerfile alone, you'll need to add the environmental variables to it

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/pass-server .

CMD ["pass-server"]
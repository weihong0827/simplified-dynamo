FROM golang:latest

# FROM alpine:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./bin/server ./server
RUN go build -o ./bin/webclient ./webclient

# COPY ./bin/server /app/bin/server
# COPY ./bin/webclient /app/bin/webclient
#
# # Make the binaries executable
# RUN chmod +x /app/bin/server
# RUN chmod +x /app/bin/webclient

EXPOSE 8080

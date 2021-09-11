FROM golang:1.16-alpine as build

WORKDIR /app
COPY ./* /app/
RUN go build -o go-cf-auth

FROM bitnami/minideb:buster as runtime 
RUN apt-get update && apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good -y
COPY --from=build /app/go-cf-auth /
CMD ./go-cf-auth

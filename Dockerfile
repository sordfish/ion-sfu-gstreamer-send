FROM sordfish/build-tools:v1.17 as build

WORKDIR /app
COPY ./* /app/
RUN apt-get update && apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good -y
RUN go build -o ion-sfu-gstreamer-send


FROM sordfish/minideb-gstreamer:latest as runtime

COPY --from=build /app/ion-sfu-gstreamer-send /
COPY startup.sh /
CMD ["/bin/sh" "startup.sh"]
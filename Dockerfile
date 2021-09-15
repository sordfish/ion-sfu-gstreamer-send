FROM sordfish/build-tools:v1.17 as build

WORKDIR /app
COPY ./* /app/
RUN go build -o ion-sfu-gstreamer-send


FROM sordfish/ubuntu-gstreamer:latest as runtime

COPY --from=build /app/ion-sfu-gstreamer-send /
COPY startup.sh /
CMD ["/bin/sh", "startup.sh"]
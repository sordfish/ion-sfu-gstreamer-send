FROM bitnami/minideb:buster

RUN apt-get update && apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good -y


#!/bin/sh
echo "Starting Listener on port 9000 for session 'test session'"
/ion-sfu-gstreamer-send -addr "localhost:50051" -session "test session" -video-src "udpsrc port=9000 ! h264parse ! rtph264pay config-interval=10 pt=96 ! rtph264depay ! avdec_h264 ! videoconvert"
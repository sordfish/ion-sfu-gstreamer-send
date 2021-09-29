#!/bin/sh
echo "Starting Listener on port $ION_CLIENT_SFU_ADDR for session $ION_CLIENT_SFU_SESSION"
echo "./ion-sfu-gstreamer-send -addr $ION_CLIENT_SFU_ADDR -session $ION_CLIENT_SFU_SESSION -video-src \"$ION_CLIENT_VIDEO_SRC\" "
./ion-sfu-gstreamer-send -addr "$ION_CLIENT_SFU_ADDR" -session "$ION_CLIENT_SFU_SESSION" -video-codec "vp8" -video-src \""$ION_CLIENT_VIDEO_SRC"\"
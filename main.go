package main

import (
	"os"

	"fmt"
	"net/http"

	ilog "github.com/pion/ion-log"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/webrtc/v3"
	gst "github.com/sordfish/go-gstreamer"
)

var (
	log = ilog.NewLoggerWithFields(ilog.DebugLevel, "", nil)
)

func healthz(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func main() {

	env_addr := os.Getenv("ISGS_ADDR")
	env_session := os.Getenv("ISGS_SESSION")
	env_videoSrc := os.Getenv("ISGS_VIDEO_SRC")
	env_audioSrc := os.Getenv("ISGS_AUDIO_SRC")
	env_turnAddr := os.Getenv("ISGS_TURN_ADDR")
	env_turnUser := os.Getenv("ISGS_TURN_USER")
	env_turnPass := os.Getenv("ISGS_TURN_PASS")

	log.Infof("This is the env addr: %s", env_addr)
	log.Infof("This is the env session: %s", env_session)
	log.Infof("This is the env videosrc: %s", env_videoSrc)
	log.Infof("This is the env audiosrc: %s", env_audioSrc)
	log.Infof("This is the env turnaddr: %s", env_turnAddr)
	log.Infof("This is the env turnaddr: %s", env_turnUser)
	log.Infof("This is the env turnaddr: %s", env_turnPass)

	servicename, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	var webrtcCfg webrtc.Configuration

	if len(env_turnAddr) > 0 {

		webrtcCfg = webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs:       []string{"turn:" + env_turnAddr},
					Username:   env_turnUser,
					Credential: env_turnPass,
				},
			},
		}

	} else {

		webrtcCfg = webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				webrtc.ICEServer{},
			},
		}

	}

	config := sdk.RTCConfig{
		WebRTC: sdk.WebRTCTransportConfig{
			VideoMime:     "video/h264",
			Configuration: webrtcCfg,
		},
	}

	connector := sdk.NewConnector(env_addr)
	rtc := sdk.NewRTC(connector, config)

	//videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264", ClockRate: 90000, Channels: 0, SDPFmtpLine: "packetization-mode=1;profile-level-id=42e01f", RTCPFeedback: nil}, "video", servicename)
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264", ClockRate: 90000, Channels: 0, RTCPFeedback: nil}, "video", servicename)
	if err != nil {
		panic(err)
	}

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", servicename)
	if err != nil {
		panic(err)
	}

	// client join a session
	err = rtc.Join(env_session, sdk.RandomKey(4))

	if err != nil {
		log.Errorf("join err=%v", err)
		panic(err)
	}
	_, _ = rtc.Publish(videoTrack, audioTrack)

	// Start pushing buffers on these tracks

	gst.CreatePipeline("opus", []*webrtc.TrackLocalStaticSample{audioTrack}, env_audioSrc).Start()
	gst.CreatePipeline("bare", []*webrtc.TrackLocalStaticSample{videoTrack}, env_videoSrc).Start()

	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":8090", nil)

	select {}

}

package main

import (
	"os"

	"fmt"
	"net/http"

	gst "github.com/pion/ion-sdk-go/pkg/gstreamer-src"

	ilog "github.com/pion/ion-log"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/webrtc/v3"
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

	// add stun servers
	webrtcCfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:       []string{"turn:" + env_turnAddr},
				Username:   env_turnUser,
				Credential: env_turnPass,
			},
		},
	}

	config := sdk.Config{
		WebRTC: sdk.WebRTCTransportConfig{
			Configuration: webrtcCfg,
			VideoMime:     "video/h264",
		},
	}
	// new sdk engine
	e := sdk.NewEngine(config)

	// get a client from engine
	c, err := sdk.NewClient(e, env_addr, "client id")

	var peerConnection *webrtc.PeerConnection = c.GetPubTransport().GetPeerConnection()

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Infof("Connection state changed: %s", state)
	})

	if err != nil {
		log.Errorf("client err=%v", err)
		panic(err)
	}

	err = e.AddClient(c)
	if err != nil {
		return
	}

	var videoTrack *webrtc.TrackLocalStaticSample
	var audioTrack *webrtc.TrackLocalStaticSample

	videoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264", ClockRate: 90000, Channels: 0, SDPFmtpLine: "packetization-mode=1;profile-level-id=42e01f", RTCPFeedback: nil}, "video", servicename)
	if err != nil {
		panic(err)
	}

	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	if env_audioSrc != "" {
		audioTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", servicename)
		if err != nil {
			panic(err)
		}
		_, err = peerConnection.AddTrack(audioTrack)
		if err != nil {
			panic(err)
		}
	}

	// client join a session
	err = c.Join(env_session, nil)

	if err != nil {
		log.Errorf("join err=%v", err)
		panic(err)
	}

	// Start pushing buffers on these tracks
	if env_audioSrc != "" {
		gst.CreatePipeline("opus", []*webrtc.TrackLocalStaticSample{audioTrack}, env_audioSrc).Start()
	}

	gst.CreatePipeline("h264", []*webrtc.TrackLocalStaticSample{videoTrack}, env_videoSrc).Start()

	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":8090", nil)

	select {}

}

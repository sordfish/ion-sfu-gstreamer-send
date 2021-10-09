package main

import (
	"os"

	gst "github.com/pion/ion-sdk-go/pkg/gstreamer-src"

	ilog "github.com/pion/ion-log"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/webrtc/v3"
)

var (
	log = ilog.NewLoggerWithFields(ilog.DebugLevel, "", nil)
)

// func varctl(envvar string, flag string) string {

// 	if flag != "" {
// 		return flag
// 	} else {
// 		return envvar
// 	}

// }

func main() {

	env_addr := os.Getenv("ISGS_ADDR")
	env_session := os.Getenv("ISGS_SESSION")
	env_videocodec := os.Getenv("ISGS_VIDEO_CODEC")
	env_videoSrc := os.Getenv("ISGS_VIDEO_SRC")
	env_audioSrc := os.Getenv("ISGS_AUDIO_SRC")

	log.Infof("This is the env addr: %s", env_addr)
	log.Infof("This is the env session: %s", env_session)
	log.Infof("This is the env videocodec: %s", env_videocodec)
	log.Infof("This is the env videosrc: %s", env_videoSrc)
	log.Infof("This is the env audiosrc: %s", env_audioSrc)

	servicename, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// add stun servers
	webrtcCfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			// webrtc.ICEServer{
			// 	URLs: []string{"stun:stun.stunprotocol.org:3478", "stun:stun.l.google.com:19302"},
			// },
		},
	}

	config := sdk.Config{
		WebRTC: sdk.WebRTCTransportConfig{
			Configuration: webrtcCfg,
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

	switch env_videocodec {
	case "vp8":
		videoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", servicename+"-video")
		if err != nil {
			panic(err)
		}
	case "h264":
		videoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264"}, "video", servicename+"-video")
		if err != nil {
			panic(err)
		}
	default:
		videoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", servicename+"-video")
		if err != nil {
			panic(err)
		}
	}

	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	if env_audioSrc != "" {
		audioTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", servicename+"-audio")
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

	switch env_videocodec {
	case "vp8":
		gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{videoTrack}, env_videoSrc).Start()
	case "h264":
		gst.CreatePipeline("h264", []*webrtc.TrackLocalStaticSample{videoTrack}, env_videoSrc).Start()
	default:
		gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{videoTrack}, env_videoSrc).Start()

	}

	select {}
}

package main

import (
	"flag"
	"os"

	gst "github.com/pion/ion-sdk-go/pkg/gstreamer-src"

	ilog "github.com/pion/ion-log"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/webrtc/v3"
)

var (
	log = ilog.NewLoggerWithFields(ilog.DebugLevel, "", nil)
)

func main() {
	// parse flag
	var session, addr, videocodec string
	flag.StringVar(&addr, "addr", "localhost:50051", "Ion-sfu grpc addr")
	flag.StringVar(&session, "session", "test room", "join session name")
	flag.StringVar(&videocodec, "video-codec", "vp8", "set video codec vp8 or h264")
	audioSrc := flag.String("audio-src", "audiotestsrc", "GStreamer audio src")
	videoSrc := flag.String("video-src", "videotestsrc", "GStreamer video src")
	flag.Parse()

	log.Infof("Running video source: %s", *videoSrc)
	log.Infof("Running audio source: %s", *audioSrc)

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
	c, err := sdk.NewClient(e, addr, "client id")

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

	switch videocodec {
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

	if audioSrc != nil {
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
	err = c.Join(session, nil)

	if err != nil {
		log.Errorf("join err=%v", err)
		panic(err)
	}

	// Start pushing buffers on these tracks
	if audioSrc != nil {
		gst.CreatePipeline("opus", []*webrtc.TrackLocalStaticSample{audioTrack}, *audioSrc).Start()
	}

	switch videocodec {
	case "vp8":
		gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{videoTrack}, *videoSrc).Start()
	case "h264":
		gst.CreatePipeline("h264", []*webrtc.TrackLocalStaticSample{videoTrack}, *videoSrc).Start()
	default:
		gst.CreatePipeline("vp8", []*webrtc.TrackLocalStaticSample{videoTrack}, *videoSrc).Start()

	}

	select {}
}

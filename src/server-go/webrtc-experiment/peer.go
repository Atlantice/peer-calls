package main

import (
	"fmt"
	"time"

	"github.com/pion/webrtc/v2"
)

func main() {

	webrtcICEServers := []webrtc.ICEServer{{
		URLs: []string{"stun:stun.l.google.com:19302"},
	}, {
		URLs: []string{"stun:global.stun.twilio.com:3478?transport=udp"},
	}}

	webrtcConfig := webrtc.Configuration{
		ICEServers: webrtcICEServers,
	}

	mediaEngine := webrtc.MediaEngine{}
	mediaEngine.RegisterDefaultCodecs()

	settingEngine := webrtc.SettingEngine{}
	// settingEngine.SetTrickle(true)
	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(mediaEngine),
		webrtc.WithSettingEngine(settingEngine),
	)
	p1, err := api.NewPeerConnection(webrtcConfig)
	chkErr(err)

	p2, err := api.NewPeerConnection(webrtcConfig)
	chkErr(err)

	p1.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
		fmt.Println("p1.OnTrack", track.Codec(), track.PayloadType(), track.Kind(), track.SSRC(), track.ID(), track.Label())
	})
	p2.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
		fmt.Println("p2.OnTrack", track.Codec(), track.PayloadType(), track.Kind(), track.SSRC(), track.ID(), track.Label())
	})

	addStateChangeListener(1, p1)
	addStateChangeListener(2, p2)

	negotiate := func() {
		fmt.Println("Negotiate start")
		offer, err := p1.CreateOffer(nil)
		// fmt.Println(offer.Type, offer.SDP)

		chkErr(err)
		err = p1.SetLocalDescription(offer)
		chkErr(err)
		err = p2.SetRemoteDescription(offer)
		chkErr(err)

		answer, err := p2.CreateAnswer(nil)
		// fmt.Println(answer.Type, answer.SDP)
		chkErr(err)
		err = p2.SetLocalDescription(answer)
		chkErr(err)
		err = p1.SetRemoteDescription(answer)
		chkErr(err)
		fmt.Println("Negotiate end")
	}

	negotiate()

	fmt.Println("p1.NewTrack")
	track, err := p1.NewTrack(webrtc.DefaultPayloadTypeVP9, 1234, "first", "track-one")
	chkErr(err)

	// go func() {
	// 	// Open a IVF file and start reading using our IVFReader
	// 	file, ivfErr := os.Open("output.ivf")
	// 	if ivfErr != nil {
	// 		panic(ivfErr)
	// 	}

	// 	ivf, header, ivfErr := ivfreader.NewWith(file)
	// 	if ivfErr != nil {
	// 		panic(ivfErr)
	// 	}

	// 	// Send our video file frame at a time. Pace our sending so we send it at the same speed it should be played back as.
	// 	// This isn't required since the video is timestamped, but we will such much higher loss if we send all at once.
	// 	sleepTime := time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000)
	// 	for {
	// 		frame, _, ivfErr := ivf.ParseNextFrame()
	// 		if ivfErr != nil {
	// 			panic(ivfErr)
	// 		}

	// 		time.Sleep(sleepTime)
	// 		if ivfErr = track.WriteSample(media.Sample{Data: frame, Samples: 90000}); ivfErr != nil {
	// 			panic(ivfErr)
	// 		}
	// 	}
	// }()

	// fmt.Println("p2.AddTransceiver")
	// _, err = p2.AddTransceiver(track.Kind())
	// chkErr(err)

	// _, err = track.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	// chkErr(err)

	fmt.Println("p1.AddTrack")
	_, err = p1.AddTrack(track)
	chkErr(err)

	negotiate()

	time.Sleep(10 * time.Second)
}

func addStateChangeListener(id int, peer *webrtc.PeerConnection) {
	peer.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		fmt.Println("Peer connection state change:", id, state)
	})
	peer.OnSignalingStateChange(func(state webrtc.SignalingState) {
		fmt.Println("Signaling state change:", id, state)
	})
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}

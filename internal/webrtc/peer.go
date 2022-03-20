package webrtc

import (
	"sync"

	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog"

	"gitlab.com/demodesk/neko/server/internal/types"
)

type WebRTCPeerCtx struct {
	mu          sync.Mutex
	logger      zerolog.Logger
	connection  *webrtc.PeerConnection
	dataChannel *webrtc.DataChannel
	changeVideo func(videoID string) error
	iceTrickle  bool
}

func (peer *WebRTCPeerCtx) CreateOffer(ICERestart bool) (*webrtc.SessionDescription, error) {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		return nil, types.ErrWebRTCConnectionNotFound
	}

	offer, err := peer.connection.CreateOffer(&webrtc.OfferOptions{
		ICERestart: ICERestart,
	})
	if err != nil {
		return nil, err
	}

	return peer.setLocalDescription(offer)
}

func (peer *WebRTCPeerCtx) CreateAnswer() (*webrtc.SessionDescription, error) {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		return nil, types.ErrWebRTCConnectionNotFound
	}

	answer, err := peer.connection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	return peer.setLocalDescription(answer)
}

func (peer *WebRTCPeerCtx) setLocalDescription(description webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	if !peer.iceTrickle {
		// Create channel that is blocked until ICE Gathering is complete
		gatherComplete := webrtc.GatheringCompletePromise(peer.connection)

		if err := peer.connection.SetLocalDescription(description); err != nil {
			return nil, err
		}

		<-gatherComplete
	} else {
		if err := peer.connection.SetLocalDescription(description); err != nil {
			return nil, err
		}
	}

	return peer.connection.LocalDescription(), nil
}

func (peer *WebRTCPeerCtx) SetOffer(sdp string) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		return types.ErrWebRTCConnectionNotFound
	}

	return peer.connection.SetRemoteDescription(webrtc.SessionDescription{
		SDP:  sdp,
		Type: webrtc.SDPTypeOffer,
	})
}

func (peer *WebRTCPeerCtx) SetAnswer(sdp string) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		return types.ErrWebRTCConnectionNotFound
	}

	return peer.connection.SetRemoteDescription(webrtc.SessionDescription{
		SDP:  sdp,
		Type: webrtc.SDPTypeAnswer,
	})
}

func (peer *WebRTCPeerCtx) SetCandidate(candidate webrtc.ICECandidateInit) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		return types.ErrWebRTCConnectionNotFound
	}

	return peer.connection.AddICECandidate(candidate)
}

func (peer *WebRTCPeerCtx) SetVideoID(videoID string) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection == nil {
		return types.ErrWebRTCConnectionNotFound
	}

	peer.logger.Info().Str("video_id", videoID).Msg("change video id")
	return peer.changeVideo(videoID)
}

func (peer *WebRTCPeerCtx) Destroy() {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	if peer.connection != nil {
		err := peer.connection.Close()
		peer.logger.Err(err).Msg("peer connection destroyed")
		peer.connection = nil
	}
}

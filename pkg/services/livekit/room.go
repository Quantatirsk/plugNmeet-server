package livekitservice

import (
	"github.com/livekit/protocol/livekit"
)

// EndRoom will send API request to livekit
func (s *LivekitService) EndRoom(roomId string) (string, error) {
	data := livekit.DeleteRoomRequest{
		Room: roomId,
	}

	res, err := s.lkc.DeleteRoom(s.ctx, &data)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}
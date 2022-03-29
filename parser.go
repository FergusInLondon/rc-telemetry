// @see https://github.com/iNavFlight/inav/wiki/Lightweight-Telemetry-(LTM)
package telemetry

import (
	"io"
	"log"
)

type DecodableFrame interface {
	FromReader(io.Reader) (byte, error)
}

type frameGetter func() DecodableFrame

var frames = map[byte]frameGetter{
	FRAME_TYPE_GPS_FRAME:          func() DecodableFrame { return &GPSFrame{} },
	FRAME_TYPE_ALTITUDE_FRAME:     func() DecodableFrame { return &AltitudeFrame{} },
	FRAME_TYPE_STATUS_FRAME:       func() DecodableFrame { return &StatusFrame{} },
	FRAME_TYPE_ORIGIN_FRAME:       func() DecodableFrame { return &OriginFrame{} },
	FRAME_TYPE_NAVIGATION_FRAME:   func() DecodableFrame { return &NavigationFrame{} },
	FRAME_TYPE_GPS_EXTENDED_FRAME: func() DecodableFrame { return &GPSExtendedFrame{} },
	FRAME_TYPE_TUNING_FRAME:       func() DecodableFrame { return &TuningFrame{} },
}

func frameFromFunctionByte(fn byte) DecodableFrame {
	if getFrame, haveFrame := frames[fn]; haveFrame {
		return getFrame()
	}

	return nil
}

const (
	_PARSER_AWAITING_MESSAGE = iota
	_PARSER_AWAITING_HEADER
	_PARSER_AWAITING_FRAME
	_PARSER_AWAITING_CHECKSUM
)

func Parse(data io.Reader) ([]DecodableFrame, error) {
	var (
		err               error
		crcByte           byte
		currentFrame      DecodableFrame
		input             = make([]byte, 1)
		state             = _PARSER_AWAITING_MESSAGE
		frames            = make([]DecodableFrame, 0)
		crcFailures       = 0
		malformedMessages = 0
		unknownFrameTypes = 0
	)

	defer log.Printf(
		"Unknown Frame Types: %d\n Malformed Frames: %d\n CRC Failures: %d\n",
		unknownFrameTypes, malformedMessages, crcFailures,
	)

	for {
		if _, err := data.Read(input); err != nil {
			return frames, err
		}

		switch state {
		case _PARSER_AWAITING_MESSAGE:
			if input[0] == '$' {
				state = _PARSER_AWAITING_HEADER
			}
		case _PARSER_AWAITING_HEADER:
			if input[0] != 'T' {
				state = _PARSER_AWAITING_MESSAGE
				continue
			}
			state = _PARSER_AWAITING_FRAME
		case _PARSER_AWAITING_FRAME:
			currentFrame = frameFromFunctionByte(input[0])
			if currentFrame == nil {
				unknownFrameTypes++
				state = _PARSER_AWAITING_MESSAGE
				continue
			}

			if crcByte, err = currentFrame.FromReader(data); err != nil {
				malformedMessages++
				state = _PARSER_AWAITING_MESSAGE
			}

			state = _PARSER_AWAITING_CHECKSUM
		case _PARSER_AWAITING_CHECKSUM:
			if crcByte != input[0] {
				crcFailures++
			} else {
				frames = append(frames, currentFrame)
			}

			state = _PARSER_AWAITING_MESSAGE
		}
	}
}

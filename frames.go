package telemetry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	FRAME_TYPE_GPS_FRAME          = 'G'
	FRAME_TYPE_ALTITUDE_FRAME     = 'A'
	FRAME_TYPE_STATUS_FRAME       = 'S'
	FRAME_TYPE_ORIGIN_FRAME       = 'O'
	FRAME_TYPE_NAVIGATION_FRAME   = 'N'
	FRAME_TYPE_GPS_EXTENDED_FRAME = 'X'
	FRAME_TYPE_TUNING_FRAME       = 'T'
)

type position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Altitude  float64 `json:"alt"`
}

func (pos *position) parsePosition(lat, long, alt int32) {
	pos.Latitude = float64(lat) / 1e7
	pos.Longitude = float64(long) / 1e7
	pos.Altitude = float64(alt) / 100
}

type GPSFrame struct {
	position
	GroundSpeed int8 `json:"spd"`
	Fix         int8 `json:"fix"`
	Sats        int8 `json:"numsat"`
}

func (gpsFrame *GPSFrame) String() string {
	return fmt.Sprintf("[GPS] lat: %.6f lon: %.6f alt: %.2f gspd: %dm/s fix: %d sats %d",
		gpsFrame.Latitude, gpsFrame.Longitude, gpsFrame.Altitude,
		gpsFrame.GroundSpeed, gpsFrame.Fix, gpsFrame.Sats)
}

func (gpsFrame *GPSFrame) FromReader(inbound io.Reader) (byte, error) {
	payload := make([]byte, 14)
	if err := readBytes(inbound, payload); err != nil {
		return byte(0), err
	}

	gpsFrame.parsePosition(
		readInt32(payload[0:4]),
		readInt32(payload[4:8]),
		readInt32(payload[9:13]),
	)

	gpsFrame.GroundSpeed = int8(payload[8])
	gpsFrame.Fix = int8(payload[13] & 3)
	gpsFrame.Sats = int8(payload[13] >> 2)
	return crcByte(payload), nil
}

type AltitudeFrame struct {
	Pitch   int16 `json:"angx"`
	Roll    int16 `json:"angy"`
	Heading int16 `json:"heading"`
}

func (altitudeFrame *AltitudeFrame) String() string {
	return fmt.Sprintf("[ALT] pitch: %d roll: %d heading: %d",
		altitudeFrame.Pitch, altitudeFrame.Roll, altitudeFrame.Heading)
}

func (altitudeFrame *AltitudeFrame) FromReader(inbound io.Reader) (byte, error) {
	payload := make([]byte, 6)
	if err := readBytes(inbound, payload); err != nil {
		return byte(0), err
	}

	altitudeFrame.Pitch = readInt16(payload[0:2])
	altitudeFrame.Roll = readInt16(payload[2:4])
	altitudeFrame.Heading = readInt16(payload[4:])
	return crcByte(payload), nil
}

type StatusFrame struct {
	BatteryVoltage     float64 `json:"vbat"`
	BatteryConsumption float64 `json:"vcurr"`
	RSSI               byte    `json:"rssi"`
	Airspeed           byte    `json:"airspeed"`
	Status             Status  `json:"status"`
	IsArmed            bool    `json:"armed"`
	IsFailsafe         bool    `json:"failsafe"`
}

func (statusFrame *StatusFrame) String() string {
	isArmed := "N"
	isFailSafe := "N"

	if statusFrame.IsArmed {
		isArmed = "Y"
	}

	if statusFrame.IsFailsafe {
		isFailSafe = "Y"
	}

	return fmt.Sprintf("[STA] vbat: %.2fV cons: %.3fAh rssi: %d aspd: %dm/s arm: %s fail: %s status: %s",
		statusFrame.BatteryVoltage, statusFrame.BatteryConsumption, statusFrame.RSSI,
		statusFrame.Airspeed, isArmed, isFailSafe, statusFrame.Status)
}

func (statusFrame *StatusFrame) parseStatus(status byte) {
	statusFrame.IsArmed = (status & 0x01) == 0x01
	statusFrame.IsArmed = (status & 0x02) == 0x02
	statusFrame.Status = Status(status >> 2)
}

func (statusFrame *StatusFrame) FromReader(inbound io.Reader) (byte, error) {
	payload := make([]byte, 7)
	if err := readBytes(inbound, payload); err != nil {
		return byte(0), err
	}

	statusFrame.BatteryVoltage = float64(readUint16(payload[0:2])) / 1000
	statusFrame.BatteryConsumption = float64(readUint16(payload[2:4])) / 1000
	statusFrame.RSSI = payload[4]
	statusFrame.Airspeed = payload[5]
	statusFrame.parseStatus(payload[6])
	return crcByte(payload), nil
}

type OriginFrame struct {
	position
	OSD bool `json:"osd"`
	Fix byte `json:"fix"`
}

func (originFrame *OriginFrame) String() string {
	isOSDOn := "Y"
	if !originFrame.OSD {
		isOSDOn = "N"
	}

	return fmt.Sprintf("[ORI] lat: %.6f lon: %.6f alt: %.2fm fix: %d osd %s",
		originFrame.Latitude, originFrame.Longitude, originFrame.Altitude,
		originFrame.Fix, isOSDOn)
}

func (originFrame *OriginFrame) FromReader(inbound io.Reader) (byte, error) {
	payload := make([]byte, 14)
	if err := readBytes(inbound, payload); err != nil {
		return byte(0), err
	}

	originFrame.parsePosition(
		readInt32(payload[0:4]),
		readInt32(payload[4:8]),
		readInt32(payload[8:12]),
	)

	originFrame.OSD = (payload[12] & 0x01) == 1
	originFrame.Fix = payload[13]
	return crcByte(payload), nil
}

type NavigationFrame struct {
	NavMode        NavMode   `json:"nav_mode"`
	GPSMode        GPSMode   `json:"gps_mode"`
	NavAction      NavAction `json:"action"`
	NavError       NavError  `json:"nav_error"`
	WaypointNumber int8      `json:"wp_number"`
	Flags          byte      `json:"flags"`
}

func (navigationFrame *NavigationFrame) String() string {
	return fmt.Sprintf("[NAV] nav: %s gps: %s act: %s err %s wpt: %d",
		navigationFrame.NavMode, navigationFrame.GPSMode, navigationFrame.NavAction,
		navigationFrame.NavError, navigationFrame.WaypointNumber)
}

func (navigationFrame *NavigationFrame) FromReader(inbound io.Reader) (byte, error) {
	payload := make([]byte, 6)
	if err := readBytes(inbound, payload); err != nil {
		return byte(0), err
	}

	navigationFrame.GPSMode = GPSMode(payload[0])
	navigationFrame.NavMode = NavMode(payload[1])
	navigationFrame.NavAction = NavAction(payload[2])
	navigationFrame.NavError = NavError(payload[4])
	navigationFrame.WaypointNumber = int8(payload[3])
	navigationFrame.Flags = payload[5]
	return crcByte(payload), nil
}

type GPSExtendedFrame struct {
	HDOP           float64 `json:"hdop"`
	HardwareStatus byte    `json:"hw_status"`
	LTMXCounter    byte    `json:"ltm_x_count"`
	DisarmReason   byte    `json:"disarm_reason"` // DisarmReason should be another enum type?
	Unused         byte    `json:"-"`
}

func (gpsExtendedFrame *GPSExtendedFrame) String() string {
	return fmt.Sprintf("[GPX] hdop: %.2f hw: 0x%x cnt: %d disarm: %d",
		gpsExtendedFrame.HDOP, gpsExtendedFrame.HardwareStatus,
		gpsExtendedFrame.LTMXCounter, gpsExtendedFrame.DisarmReason)
}

func (gpsExtendedFrame *GPSExtendedFrame) FromReader(inbound io.Reader) (byte, error) {
	payload := make([]byte, 6)
	if err := readBytes(inbound, payload); err != nil {
		return byte(0), err
	}

	gpsExtendedFrame.HDOP = float64(readUint16(payload[0:2])) / 100
	gpsExtendedFrame.HardwareStatus = payload[2]
	gpsExtendedFrame.LTMXCounter = payload[3]
	gpsExtendedFrame.DisarmReason = payload[4]
	gpsExtendedFrame.Unused = payload[5]
	return crcByte(payload), nil
}

type TuningFrame struct {
	// not txed
	PRoll, IRoll, DRoll          byte
	PPitch, PPitch2, IPitch      byte
	DYaw, IYaw                   byte
	RollRate, PitchRate, YawRate byte
}

func (tuningFrame *TuningFrame) FromReader(inbound io.Reader) (byte, error) {
	// This is just here for documentation reasons; I want to look in to the inav
	// source code and update their docs for this.
	return byte(0), nil
}

func readInt16(in []byte) int16 {
	return int16(readUint16(in))
}

func readUint16(in []byte) uint16 {
	return binary.LittleEndian.Uint16(in)
}

func readInt32(in []byte) int32 {
	return int32(readUint32(in))
}

func readUint32(in []byte) uint32 {
	return binary.LittleEndian.Uint32(in)
}

func readBytes(inbound io.Reader, buf []byte) error {
	nRead, err := inbound.Read(buf)
	if err != nil {
		return err
	}

	if nRead != len(buf) {
		return errors.New("invalid number of bytes for frame type")
	}

	return nil
}

func crcByte(buf []byte) byte {
	crc := byte(0)
	for _, b := range buf {
		crc = crc ^ b
	}

	return crc
}

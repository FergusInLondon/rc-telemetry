# LTM Parsing Example

## Usage

```
➜  telemetry ./parser example.data
2022/03/29 21:21:12 Unknown Frame Types: 0
 Malformed Frames: 0
 CRC Failures: 0
File 'example.data' contains 5978 frames.

➜  telemetry ./parser example.data --verbose
  ...
[GPS] lat: 50.910427 lon: -1.535210 alt: 0.05 gspd: 0m/s fix: 3 sats 9
[ALT] pitch: 4 roll: 0 heading: 111
[STA] vbat: 10.70V cons: 0.000Ah rssi: 0 aspd: 0m/s arm: N fail: N status: STATUS_HORIZON
[ORI] lat: 50.910484 lon: -1.535072 alt: 0.00m fix: 1 osd Y
[ALT] pitch: 4 roll: 0 heading: 111
[GPS] lat: 50.910429 lon: -1.535209 alt: 0.06 gspd: 0m/s fix: 3 sats 9
[ALT] pitch: 4 roll: 0 heading: 111
[STA] vbat: 10.70V cons: 0.000Ah rssi: 0 aspd: 0m/s arm: N fail: N status: STATUS_HORIZON
[ALT] pitch: 4 roll: 0 heading: 111
[GPS] lat: 50.910430 lon: -1.535208 alt: 0.10 gspd: 0m/s fix: 3 sats 9
➜  telemetry 
```

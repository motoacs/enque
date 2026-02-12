# enque

Enque is a Windows desktop batch encoder frontend for rigaya's NVEncC.

## Scope (v1)
- Encoder execution: NVEncC 8.x+
- Queue, profile, output resolution, logging, metadata file-time restore
- QSVEnc/ffmpeg: stubbed adapter contract (`E_ENCODER_NOT_IMPLEMENTED`)

## External tools
- NVEncC: [rigaya NVEnc releases](https://github.com/rigaya/NVEnc/releases)
- QSVEnc: [rigaya QSVEnc releases](https://github.com/rigaya/QSVEnc/releases)
- ffmpeg: [official builds](https://ffmpeg.org/download.html)

## Built-in presets
- `HEVC Quality`: balanced HEVC quality/speed
- `AV1 Fast`: throughput-oriented AV1
- `Camera Archive`: metadata + file-time preservation
- `H.264 Compatible`: compatibility-first output

## Build and test
```bash
go test ./...
cd frontend && npm install && npm run build && npm run test
wails dev
wails build -platform windows/amd64
```

## Notes
- Target OS is Windows 11 x64.
- Progress parse failures do not stop encoding; raw stderr remains available in logs.
- Job artifacts are written under `%APPDATA%/Enque/` (`config.json`, `profiles.json`, `logs/`, `runtime/temp_index.json`).

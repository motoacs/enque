# Enque

A professional batch video encoder frontend for [rigaya's NVEncC](https://github.com/rigaya/NVEnc). Built for NVIDIA GPU users who need efficient, concurrent batch encoding with advanced profile management and metadata preservation.

**[Japanese / 日本語版はこちら](README.ja.md)**

## Overview

Enque is a Windows desktop application designed for users who understand encoder CLI options and want to streamline batch encoding workflows. It provides GUI widgets for commonly used NVEncC options while allowing arbitrary CLI options via text input, giving you full control over the encoding process.

### Key Features

- **Batch Encoding** -- Drag & drop multiple files or folders to build a job queue, then encode them all at once
- **Concurrent Execution** -- Run up to 8 NVEncC processes in parallel to maximize GPU utilization (e.g. RTX 5090's 3 NVENC engines)
- **Profile Management** -- Save, duplicate, and switch between encoding presets. 4 built-in presets included
- **Full NVEncC Option Coverage** -- GUI widgets for major settings + free-form text input for any NVEncC option
- **Real-time Progress** -- Per-job progress bars with fps, bitrate, and ETA. Live stderr log viewer
- **Command Preview** -- Always see the exact NVEncC command line that will be executed, with clipboard copy
- **File Timestamp Preservation** -- Restore original file creation/modification times after encoding (Win32 API)
- **Metadata Pass-through** -- Copy container metadata, chapters, subtitles, data tracks, and attachments
- **Job Control** -- Skip individual jobs, graceful stop (finish current then stop), or abort all
- **Post-Encode Actions** -- Automatically shutdown, sleep, or run a custom command after all jobs complete
- **Bilingual UI** -- Japanese and English interface

## Requirements

- **OS**: Windows 11 (x64)
- **GPU**: NVIDIA GPU with NVENC support
- **NVEncC**: [NVEncC64.exe](https://github.com/rigaya/NVEnc) 8.x or later (required)

### Optional

- QSVEncC64.exe -- Detected and stored for future support
- ffmpeg.exe / ffprobe.exe -- Detected and stored for future support

## Installation

1. Download the latest release ZIP
2. Extract to any folder
3. Place `NVEncC64.exe` in the same folder as `Enque.exe`, or ensure it's in your system PATH
4. Run `Enque.exe`

On first launch, Enque will auto-detect NVEncC in the application folder and PATH. If not found, you'll be prompted to set the path manually in Settings.

No installer required. Configuration is stored in `%APPDATA%\Enque\`.

## Quick Start

1. **Add files** -- Drag & drop video files onto the window, or use the file/folder dialog
2. **Select a profile** -- Choose from built-in presets or create your own
3. **Configure output** -- Set the output folder, filename template, and container format
4. **Review the command** -- Check the command preview at the bottom of the profile editor
5. **Start encoding** -- Click the start button and monitor progress in real-time

## Built-in Presets

| Preset | Codec | Quality | Preset | Notes |
|--------|-------|---------|--------|-------|
| **HEVC Quality** | HEVC | QVBR 28 | P4 | 10-bit, balanced speed/quality, split-enc auto |
| **AV1 Fast** | AV1 | QVBR 32 | P1 | 10-bit, maximum throughput, split-enc auto |
| **Camera Archive** | HEVC | QVBR 24 | P7 | 10-bit, max quality, audio copy, all metadata + file timestamps preserved |
| **H.264 Compatible** | H.264 | QVBR 26 | P4 | 8-bit, AAC 256kbps, maximum device compatibility |

Presets cannot be edited directly. Duplicate a preset to create your own customized version.

## Profile Settings

### GUI Options (Layer 1)

Enque provides GUI widgets for the following NVEncC options:

- **Video**: Codec (H.264/HEVC/AV1), rate control (QVBR/CQP/CBR/VBR), quality/bitrate, preset (P1-P7), output depth (8/10-bit), multipass, output resolution
- **Video Detail**: B-frames, reference frames, lookahead, GOP length, spatial/temporal AQ
- **Speed**: Split encoding, parallel encoding, decoder (avhw/avsw), device selection
- **Audio**: Copy / AAC / Opus with bitrate control
- **Color/HDR**: Color matrix, transfer, primaries, range, HDR10+ pass-through
- **Metadata**: Container/video/audio metadata copy, chapters, subtitles, data tracks, attachments, file timestamp restoration
- **Advanced**: Interlace, input/output CSP, tune, max bitrate, VBR quality, weighted P frames, MV precision, level/profile/tier, SSIM/PSNR metrics, trim/seek, and more

### Custom Options (Layer 2)

Any NVEncC option not covered by the GUI can be entered as free-form text. These are appended after GUI options with "later wins" precedence.

### Priority Order

```
GUI Standard Options -> GUI Advanced Options -> Custom Options (highest priority)
```

## Parallel Encoding

Enque supports two layers of parallelization:

- **Process-level** (app setting): Run 1-8 NVEncC processes simultaneously
- **Job-level** (profile setting): NVEncC's `--split-enc` and `--parallel` options

Default: 1 concurrent process + split-enc auto. For many short clips, try 2-3 concurrent processes with split-enc off.

## Output Settings

- **Output folder**: Same as input file, or a specified folder
- **Filename template**: Use `{name}` (original name without extension) and `{ext}` (output container extension). Default: `{name}_encoded.{ext}`
- **Container**: mp4, mkv, mov, webm, etc.
- **Overwrite handling**: Prompt for confirmation, or auto-rename with sequential numbering

Encoding outputs to a temporary file (`{name}.{id}.tmp.{ext}`) first, then renames to the final path on success. This prevents incomplete files from appearing as final output.

## Job Control During Encoding

- **Skip**: Mark individual pending jobs to be skipped
- **Cancel**: Force-terminate a specific running job
- **Stop After Current**: Let running jobs finish, then skip the rest
- **Abort All**: Force-terminate all running jobs immediately

## Configuration

Settings are stored in `%APPDATA%\Enque\`:

```
%APPDATA%\Enque\
  config.json         # App settings
  profiles.json       # Encoding profiles
  logs/
    {job_id}.json     # Job execution record (command line, exit code, timing)
    {job_id}.stderr.log  # Full encoder stderr output
```

### App Settings

| Setting | Default | Description |
|---------|---------|-------------|
| Concurrent jobs | 1 | Number of simultaneous encoder processes (1-8) |
| On error | Skip | Skip failed job and continue, or stop the queue |
| Decoder fallback | Off | Auto-retry with --avsw if --avhw fails |
| Output timeout | 600s | Kill process if no stderr output for this long |
| Progress timeout | 300s | Kill process if progress stalls for this long |
| Post-complete action | None | Shutdown, sleep, or custom command after all jobs |

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.23 |
| Desktop Framework | [Wails](https://wails.io/) v2 |
| Frontend | React 18 + TypeScript + Vite |
| State Management | Zustand |
| UI Components | Radix UI + Tailwind CSS |
| Persistence | JSON files |
| Process Management | Windows Job Objects |

## Building from Source

### Prerequisites

- Go 1.23+
- Node.js 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

### Build

```bash
# Install dependencies
cd frontend && npm install && cd ..

# Development mode (hot reload)
wails dev

# Production build
wails build
```

The built executable will be at `build/bin/Enque.exe`.

### Run Tests

```bash
# Backend tests
go test ./...

# Frontend type check
cd frontend && npm run build
```

## Architecture

Enque follows a clear separation between frontend and backend with a well-defined API boundary:

- **Before encoding**: Frontend (Zustand stores) is the source of truth for queue and profile editing
- **During encoding**: Backend (Go) is the source of truth. Progress is pushed to frontend via Wails events
- **Encoder abstraction**: Registry + adapter pattern. v1 implements NVEncC; QSVEncC and ffmpeg adapters can be added without changing core logic

```
React + Zustand (UI & State)
    |
    +-- Wails Bindings (function calls)
    +-- Wails Events (Go -> JS, real-time updates)
    |
Go Backend
    +-- queue/     Session management, worker pool
    +-- encoder/   Adapter registry, process execution, timeout guard
    |   +-- nvencc/  Command builder, progress parser
    +-- profile/   CRUD, migration, presets
    +-- config/    App settings, migration
    +-- detector/  Tool detection, GPU capabilities
    +-- metadata/  Win32 file timestamp restoration
    +-- logging/   Job records, stderr capture
```

## License

[MIT License](LICENSE) - Copyright (c) 2026 motoacs

## Acknowledgments

- [rigaya/NVEnc](https://github.com/rigaya/NVEnc) -- The NVEncC encoder that Enque wraps
- [Wails](https://wails.io/) -- Go + Web desktop application framework
- [Radix UI](https://www.radix-ui.com/) -- Accessible UI component primitives

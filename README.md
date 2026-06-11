# Komga SSR
Simple Komga web client built in go. Built since my iPad mini 3 (running iOS 12.5.7) doesn't support ES6 to run recent Komga web pages.

## Features
* Browse libraries, series and books
* Read comics page-by-page with LTR and RTL support
* Syncs read progress back to your Komga server
* Preloads the next page for faster reading

## Requirements
* Go 1.26+ (maybe you can run this with an earlier version, only uses standard libraries)
* A running Komga instance accessible to your network

## Configuration
Config is done via environment variables, personally set it in systemd unit.
* `KOMGA_URL`: Base URL of your Komga instance
* `KOMGA_USERNAME`: Komga account username
* `KOMGA_PASSWORD`: Komga account password
* `KOMGA_SSR_PORT`: Port to host this page, defaults to 25601

## Building
```go
go build -o komga-ssr .
```

## Usage
Navigate to `http://your-host:25601` (or the port you overwrote with) in your browser.
* **Home**: list all Komga libraries
* **Library**: shows series in that library with cover thumbnails
* **Series**: shows volumes/books with read progress status
* **Reader**: page by page reader with prev/next navigation (can also click image to go to next), direction toggle, and progress sync

## Notes
* Images are served as JPEG (for maximum compatability) and cached for 1 hour in your device
* Read progress is synced async (lol) in the background so page turns aren't delayed
* Minimal UI on purpose, some inline styles and small CSS block, no major JS usage

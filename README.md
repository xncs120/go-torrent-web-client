# GoTorrentWebClient
GoTorrentWebClient is a simple localhost torrent client similar to uTorrent web without any ads or bloatware. This is an open source project, feel free to reuse the code and extend more function as you like.

## Features
- [x] Simple Webpage UI (Home page and MediaPlayer page)
- [x] Start download torrent content with progress detail
- [x] Resume back to download the torrent content that stopped
- [x] Streaming the music/video that is downloading

P.S.:
- To resume download stopped torrent, the partially downloaded file need to be in the original path.
- To stream the media the file need to be in the original path, also not all video type support partial downloaded stream.\
-- .mp4 depend on its encoding\
-- .webm work best\
-- .mkv have picture but dont have audio

## Getting started
### Installation
As dev:
1. Install golang version >= 1.23.4
2. Fork this project code into your local
```sh
cd project-folder
go mod tidy

// linux build
go build -o GoTorrentWebClient
// windows build
GOOS=windows GOARCH=amd64 go build -o GoTorrentWebClient.exe
```

As for windows user:
1. Just download the builded .exe from my github [HERE](https://github.com/xncs120/go-torrent-web-client/releases/tag/v0.2.0-alpha)

### How to use
1. Create a folder called GoTorrentWebClient (or any name)
2. Move the builded GoTorrentWebClient.exe in to the folder created
3. Then run the builded GoTorrentWebClient.exe
4. Open your browser and go to http://localhost:60000/

## Reference and external source
- [Anacrolix torrent package](https://github.com/anacrolix/torrent)

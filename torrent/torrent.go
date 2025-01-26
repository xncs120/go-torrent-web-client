package torrent

import (
	"fmt"
	"sync"

	"github.com/anacrolix/torrent"

	"go-torrent-web-client/database"
)

type TorrentManager struct {
	Client     *torrent.Client
	torrents   map[string]*torrent.Torrent
	torrentMux sync.Mutex
	jd         *database.JsonData
}

func NewTorrentManager(jd *database.JsonData, dataDir string) (*TorrentManager, error) {
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = dataDir

	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return &TorrentManager{
		Client:   client,
		torrents: make(map[string]*torrent.Torrent),
		jd:       jd,
	}, nil
}

func prioritizeStartAndEnd(t *torrent.Torrent) {
	<-t.GotInfo()

	numPieces := t.NumPieces()
	piecesToPrioritize := 10

	for i := 0; i < piecesToPrioritize && i < numPieces; i++ {
		t.Piece(i).SetPriority(torrent.PiecePriorityHigh)
	}

	for i := numPieces - piecesToPrioritize; i < numPieces; i++ {
		if i >= 0 {
			t.Piece(i).SetPriority(torrent.PiecePriorityHigh)
		}
	}
}

func (tm *TorrentManager) AddDownload(source string) (string, error) {
	tm.torrentMux.Lock()
	defer tm.torrentMux.Unlock()

	var t *torrent.Torrent
	var err error

	t, err = tm.Client.AddMagnet(source)
	if err != nil {
		return "", fmt.Errorf("error adding magnet: %v", err)
	}

	go func() {
		<-t.GotInfo()
		prioritizeStartAndEnd(t)
		t.DownloadAll()
	}()

	tm.torrents[t.Name()] = t
	tm.jd.Insert(database.TorrentData{
		Name:   t.Name(),
		Source: source,
	})
	return t.Name(), nil
}

func (tm *TorrentManager) GetProgresses() map[string]float64 {
	tm.torrentMux.Lock()
	defer tm.torrentMux.Unlock()

	progress := make(map[string]float64)
	for name, t := range tm.torrents {
		if t.Info() != nil {
			progress[name] = float64(t.BytesCompleted()) / float64(t.Info().TotalLength()) * 100
		}
	}
	return progress
}

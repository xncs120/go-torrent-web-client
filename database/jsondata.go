package database

import (
	"encoding/json"
	"os"
)

type TorrentData struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type JsonData struct {
	FilePath string
	Data     map[string]TorrentData
}

func NewJsonData(filePath string) (*JsonData, error) {
	jd := &JsonData{
		FilePath: filePath,
		Data:     make(map[string]TorrentData),
	}

	if err := jd.load(); err != nil {
		return nil, err
	}
	return jd, nil
}

func (jd *JsonData) Insert(torrentData TorrentData) error {
	jd.Data[torrentData.Name] = torrentData
	return jd.save()
}

func (jd *JsonData) SelectAll() ([]TorrentData, error) {
	torrentData := make([]TorrentData, 0, len(jd.Data))
	for _, torrent := range jd.Data {
		torrentData = append(torrentData, torrent)
	}
	return torrentData, nil
}

func (jd *JsonData) SelectOne(name string) TorrentData {
	return jd.Data[name]
}

func (jd *JsonData) save() error {
	file, err := os.OpenFile(jd.FilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(jd.Data)
}

func (jd *JsonData) load() error {
	file, err := os.Open(jd.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&jd.Data)
}

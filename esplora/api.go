package esplora

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type EsploraAPI struct {
	URL string
}

func NewEsploraAPI(url string) *EsploraAPI {
	return &EsploraAPI{URL: url}
}

func (e *EsploraAPI) get(path string) ([]byte, error) {
	p, err := url.JoinPath(e.URL, path)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (e *EsploraAPI) GetUTXOs(address string) ([]*UTXO, error) {
	path := "/address/" + address + "/utxo"
	body, err := e.get(path)
	if err != nil {
		return nil, err
	}

	var utxos []*UTXO
	err = json.Unmarshal(body, &utxos)
	if err != nil {
		log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}

	return utxos, nil
}

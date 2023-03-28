package blob

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RemoteOracle struct {
	baseUrl string
}

func NewRemoteOracle(baseUrl string) *RemoteOracle {
	return &RemoteOracle{
		baseUrl,
	}
}

func (self *RemoteOracle) GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]Blob, []KZGCommitment, error) {
	//todo: support mulity version
	if len(versions) != 1 {
		return nil, nil, errors.New("only support 1 version query")
	}

	resp, err := http.Get(fmt.Sprintf("%s?versionHash=0x%x", self.baseUrl, versions[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("get remote version failed: %w", err)
	}
	defer resp.Body.Close()
	var ret BlobWithCommitment
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, nil, err
	}
	return []Blob{ret.Blob}, []KZGCommitment{ret.Commitment}, nil
}

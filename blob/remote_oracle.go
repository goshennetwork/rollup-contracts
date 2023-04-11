package blob

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//todo: If EIP4844 is enabled, just need to modify logic of this Oracle
type RemoteOracle struct {
	baseUrl string
}

func NewRemoteOracle(baseUrl string) *RemoteOracle {
	if len(baseUrl) == 0 {
		panic(1)
	}
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	return &RemoteOracle{
		baseUrl,
	}
}

func (self *RemoteOracle) GetBlobsWithCommitmentVersions(versions ...[32]byte) ([]Blob, []KZGCommitment, error) {
	//todo: support mulity version
	if len(versions) != 1 {
		return nil, nil, errors.New("only support 1 version query")
	}

	url := fmt.Sprintf("%sblobOracle?versionHash=0x%x", self.baseUrl, versions[0])
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("get remote version failed: %w, url: %s", err, url)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var ret BlobWithCommitment
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, nil, err
	}
	return []Blob{ret.Blob}, []KZGCommitment{ret.Commitment}, nil
}

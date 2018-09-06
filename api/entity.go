package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Entity interface for json http response
type Entity interface {
	UnmarshalHTTP(*http.Response) error
}

func GetEntity(r *http.Response, v Entity) error {
	return v.UnmarshalHTTP(r)
}

func readEntityFromResp(resp *http.Response, dat Entity) error {
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, &dat)
}

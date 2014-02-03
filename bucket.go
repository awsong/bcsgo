package bcsgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Bucket struct {
	bcs  *BCS
	Name string `json:"bucket_name"`
}

func (this *Bucket) getUrl() string {
	return this.bcs.simpleSign(GET, this.Name, "/")
}
func (this *Bucket) putUrl() string {
	return this.bcs.simpleSign(PUT, this.Name, "/")
}
func (this *Bucket) Create() error {
	link := this.putUrl()
	resp, _, err := this.bcs.httpClient.Put(link, nil, 0)
	if resp.StatusCode != http.StatusOK {
		err = errors.New("request not ok, status: " + string(resp.StatusCode))
	}
	return err
}
func (this *Bucket) Object(absolutePath string) *Object {
	o := Object{}
	o.bucket = this
	o.AbsolutePath = absolutePath
	return &o
}

func (this *Bucket) ListObjects(prefix string, start, limit int) (*ObjectCollection, error) {
	params := url.Values{}
	params.Set("start", string(start))
	params.Set("limit", string(limit))
	link := this.getUrl() + "&" + params.Encode()
	_, data, err := this.bcs.httpClient.Get(link)
	fmt.Println(string(data))
	if err != nil {
		return nil, err
	} else {
		var objectsInfo ObjectCollection
		err := json.Unmarshal(data, &objectsInfo)
		fmt.Println(objectsInfo)
		if err != nil {
			return nil, err
		} else {
			for i, _ := range objectsInfo.Objects {
				objectsInfo.Objects[i].bucket = this
			}
			return &objectsInfo, nil
		}
	}
}

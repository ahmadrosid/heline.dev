package solr

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Map map[string]interface{}

func Insert(payload io.Reader) error {

	url := "http://localhost:8984/solr/heline/update?&commitWithin=1000&overwrite=true&wt=json"
	// url := "http://heline.dev:8984/solr/heline/update?&commitWithin=1000&overwrite=true&wt=json"

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		println(err.Error())
		return err
	}

	req.Header.Add("Content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err.Error())
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	return nil
}

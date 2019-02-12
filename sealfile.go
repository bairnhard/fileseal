package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2" //YAML Parser for external configuration
)

//Conf stores config parameters
type Conf struct {
	APIKey  string `yaml:"apikey"` //CWSeal API Key
	SealURL string `yaml:"baseurl"`
}

type Cryptowerk struct {
	MaxSupportedAPIVersion int `json:"maxSupportedAPIVersion"`
	Documents              []struct {
		RetrievalID string `json:"retrievalId"`
	} `json:"documents"`
	MinSupportedAPIVersion int `json:"minSupportedAPIVersion"`
}

var (
	Cfg Conf // Config File Content
)

func main() {

	readconfig() //get config file content

	file := filename()
	cryptowerkapi := Cfg.SealURL

	//get the hash
	hashresult := filehasher(file)
	fmt.Println("...Registering to Blockchain")
	retrievalId, err := registerToBlockchain(hashresult, cryptowerkapi, Cfg.APIKey)
	errlog(err)

	fmt.Println("RetrievalId: ", retrievalId)

	//log retrieval ID to CSV
	seallog(retrievalId, file)

}

func readconfig() {

	fmt.Println("Reading config")

	confFile, err := ioutil.ReadFile("sealfile.cfg")
	errlog(err)

	err = yaml.Unmarshal(confFile, &Cfg)
	errlog(err)

}

func errlog(err error) {
	if err != nil {
		log.Println(time.Now(), err)
	}
}

func filename() string { // reads filename from command line with flag -f

	wordPtr := flag.String("f", "demo.txt", "filename")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	returnvalue := *wordPtr
	return returnvalue
}

func filehasher(filename string) string { //here we hash the files with SHA256
	f, err := os.Open(filename)
	if err != nil {
		errlog(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		errlog(err)
	}
	hstring := hex.EncodeToString(h.Sum(nil))
	return hstring
}

func registerToBlockchain(bufhash string, cryptowerkapi string, cryptowerkkey string) (retrievalId string, err error) {

	headers := map[string][]string{
		"Accept":    []string{"application/json"},
		"X-API-Key": []string{cryptowerkkey},
	}

	data := url.Values{}
	data.Set("version", "6")
	data.Add("hashes", bufhash)

	req, err := http.NewRequest("POST", cryptowerkapi, bytes.NewBufferString(data.Encode()))
	errlog(err)

	req.Header = headers

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cryptowerk responded with error %d", err)
	}

	defer resp.Body.Close()

	// Successful responses will always return status code 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cryptowerk responded with unexepcted status code %d", resp.StatusCode)
	}

	seals := &Cryptowerk{}
	if err := json.NewDecoder(resp.Body).Decode(seals); err != nil {
		return "", fmt.Errorf("unable to decode Cryptowerk response: %s", err)
	}

	registerID := seals.Documents[0].RetrievalID

	return registerID, nil
}

func seallog(retrievalid string, filename string) { //write and append retrieval id to seal log

	t := time.Now()

	f, err := os.OpenFile("seallog.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	errlog(err)
	defer f.Close()
	var data = []string{t.String(), retrievalid, filename}

	w := csv.NewWriter(f)
	w.Write(data)
	w.Flush()
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2" //YAML Parser for external configuration
)

//Conf stores config parameters
type Conf struct {
	APIKey        string `yaml:"apikey"` //CWSeal API Key
	apiCredential string `yaml:"apicredential"`
	apiVersion    string `yaml:"version"`
}

var (
	Cfg Conf // Config File Content
)

func main() {

	readconfig()
	dateiname := filename()

	fmt.Println("Der Dateiname ist: %s", dateiname)
	fmt.Println("APIKey: %s", Cfg.APIKey)

	//get the hash
	hashresult := filehasher(dateiname)
	fmt.Println("Hash in Main: %s", hashresult)

	//	headers := map[string][]string{
	//		"Accept":    []string{"application/json"},
	//		"X-API-Key": []string{Cfg.APIKey, " ", Cfg.apiCredential},
	//	}
	/*
		data := bytes.NewBuffer([]byte{jsonReq})
		req, err := http.NewRequest("POST", "https://developers.cryptowerk.com/platform/API/v6/register", data)
		errlog(err)
		req.Header = headers

		client := &http.Client{}
		resp, err := client.Do(req)
		errlog(err)
	*/
}

func readconfig() {

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

	fmt.Println("f:", *wordPtr)
	fmt.Println("tail:", flag.Args())
	returnvalue := *wordPtr
	return returnvalue
}

func filehasher(filename string) string {
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
	fmt.Printf(" Hash: %x", h.Sum(nil))
	return hstring
}

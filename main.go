package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

// env or default if envvar does not exist
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// authenticate to backend: IN apiKey, OUT token
func makeAuthRequest(apiKey string) string {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://rep.checkpoint.com/rep-auth/service/v1.0/request", nil)
	req.Header.Set("Client-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// log.Println(string(body))
	return string(body)
}

func CheckResource(res string, apiKey string, token string) string {
	type RequestDetail struct {
		Resource string `json:"resource"`
	}
	type RequestDetails []RequestDetail
	type RequestPayload struct {
		Request RequestDetails `json:"request"`
	}
	data := RequestPayload{RequestDetails{RequestDetail{res}}}

	payloadBytes, err := json.Marshal(data)

	if err != nil {
		log.Fatalln(err)
	}
	body := bytes.NewReader(payloadBytes)

	var url string
	switch {
	case IsIpv4Net(res):
		url = fmt.Sprintf("https://rep.checkpoint.com/ip-rep/service/v2.0/query?resource=%s", res)
	case IsHash(res):
		url = fmt.Sprintf("https://rep.checkpoint.com/file-rep/service/v2.0/query?resource=%s", res)

	default:
		url = fmt.Sprintf("https://rep.checkpoint.com/url-rep/service/v2.0/query?resource=%s", res)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Client-Key", apiKey)
	req.Header.Set("Token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	return bodyString
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Usage: /query/resource/(ip|hash|domain)")
}

func QueryResource(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	res := vars["resource"]

	w.Header().Add("Content-Type", "application/json")

	jsonResp := CheckResource(res, _apiKey, _token)
	fmt.Fprintln(w, jsonResp)

}

func serve() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/query/resource/{resource}", QueryResource)
	port := getEnv("PORT", "60080")
	fmt.Printf("Listening on port %v\n usage http://127.0.0.1:%v/query/resource/(ip|hash|domain)\n\n", port, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

var _apiKey = ""
var _token = ""

func IsIpv4Net(host string) bool {
	return net.ParseIP(host) != nil
}

var (
	md5Regex, _    = regexp.Compile(`^([0-9a-fA-F]){32}$`)
	sha1Regex, _   = regexp.Compile(`^([0-9a-fA-F]){40}$`)
	sha256Regex, _ = regexp.Compile(`^([0-9a-fA-F]){64}$`)
)

func IsTrimmedMatching(s string, r *regexp.Regexp) bool {
	s = strings.Trim(s, " ")
	return r.MatchString(s)
}
func IsMd5Hash(hash string) bool {
	return IsTrimmedMatching(hash, md5Regex)
}

func IsSha1Hash(hash string) bool {
	return IsTrimmedMatching(hash, sha1Regex)
}

func IsSha256Hash(hash string) bool {
	return IsTrimmedMatching(hash, sha256Regex)
}

func IsHash(hash string) bool {
	return IsMd5Hash(hash) || IsSha1Hash(hash) || IsSha256Hash(hash)
}

func showHashes() {
	s := "Foo"

	md5 := md5.Sum([]byte(s))
	sha1 := sha1.Sum([]byte(s))
	sha256 := sha256.Sum256([]byte(s))

	fmt.Printf("%x %d\n", md5, 2*len((md5)))
	fmt.Printf("%x %d\n", sha1, 2*len(sha1))
	fmt.Printf("%x %d\n", sha256, 2*len(sha256))
}

func main() {

	var domain string
	flag.StringVar(&domain, "domain", "kancelareroku.eu", "domain to query in TheatCloud")

	var apiKey string
	flag.StringVar(&apiKey, "apikey", getEnv("TC_API_KEY", "bring-your-own-api-key"), "service API key")
	flag.Parse()

	token := makeAuthRequest(apiKey)

	_apiKey = apiKey
	_token = token

	// CheckResource(domain, apiKey, token)

	serve()
}

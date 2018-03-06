package msfrpc

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/vmihailenco/msgpack"
)

// MSFRPC is the entrypoint for MSFRPC connections
type MSFRPC struct {
	// Configuration
	Host     string
	Port     string
	URI      string
	Username string
	Password string
	Ssl      bool
	// Runtime
	isConnected bool
	authToken   string
}

////////////////////////////////////////////////////////////
// Constructor
////////////////////////////////////////////////////////////

// NewMsfrpc create a new MSFRPC object with specified parameters
func NewMsfrpc(host string, port string, uri string, username string, password string, ssl bool) *MSFRPC {
	msfrpc := MSFRPC{
		Host:     host,
		Port:     port,
		URI:      uri,
		Username: username,
		Password: password,
		Ssl:      ssl,
	}
	return &msfrpc
}

// Login login to rpc server
func (msfrpc *MSFRPC) Login() error {
	var result struct {
		Result string `msgpack:"result"`
		Token  string `msgpack:"token"`
	}
	err := msfrpc.CallAndUnmarshall("auth.login", []interface{}{msfrpc.Username, msfrpc.Password}, &result)
	msfrpc.authToken = result.Token
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////
// Remote calls
////////////////////////////////////////////////////////////

// CallAndUnmarshall call rpc method and unmarshal in data interface
func (msfrpc *MSFRPC) CallAndUnmarshall(method string, options []interface{}, data interface{}) error {
	stringBody, err := msfrpc.Call(method, options)
	if err != nil {
		return err
	}
	decodeMsgpack([]byte(stringBody), data)
	return nil
}

// Call call rpc method and return output
func (msfrpc *MSFRPC) Call(method string, options []interface{}) (string, error) {
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	requestData := []interface{}{method}
	if method != "auth.login" {
		requestData = append(requestData, msfrpc.authToken)
	}
	requestData = append(requestData, options...)
	requestBody, err := encodeMsgpack(requestData)
	if err != nil {
		return "", err
	}

	scheme := "http"
	if msfrpc.Ssl {
		scheme += "s"
	}
	request, err := http.NewRequest("POST", scheme+"://"+msfrpc.Host+":"+msfrpc.Port+msfrpc.URI, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "binary/message-pack")
	request.Header.Set("Accept", "binary/message-pack")
	request.Header.Set("Accept-Charset", "UTF-8")

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	stringBody := fmt.Sprintf("%s", rawBody)

	return stringBody, nil
}

func encodeMsgpack(data interface{}) ([]byte, error) {
	return msgpack.Marshal(data)
}

func decodeMsgpack(bytes []byte, destination interface{}) {
	msgpack.Unmarshal(bytes, destination)
}

// safeString return a safe string for both meterpreter and classic command line
func safeString(input string) string {
	if regexp.MustCompile(`^[^:]+:\\$`).MatchString(input) {
		return input
	}
	return "\"" + input + "\""
}

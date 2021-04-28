package policyservice

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func (ps *PolicyServer) httpServe(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/invoke", ps.InvokePolicyHttp)
	httpServer := &http.Server{Handler: mux}
	return httpServer.ServeTLS(l, ps.Config.CertFile, ps.Config.KeyFile)
}

func (ps *PolicyServer) InvokePolicyHttp(w http.ResponseWriter, req *http.Request) {
	log.Info("CALLED")

	bodyReader := bufio.NewReader(req.Body)
	headers, err := ps.parseHeaders(bodyReader)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse header content"))
		return
	}

	body, err := ps.parseBody(bodyReader)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse body content"))
		return
	}

	// blah blah
	headers["newone"] = "some header policy"
	body = reverseBody(body)

	// Write the response.
	w.Header().Set("content-type", "application/octet-stream")
	for hdr, value := range headers {
		w.Write([]byte(fmt.Sprintf("%s: %s\r\n", hdr, value)))
	}
	w.Write([]byte("\r\n"))
	w.Write(body)
}

func (ps *PolicyServer) parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := map[string]string{}
	for {
		hdr, err := reader.ReadString('\n')
		if err == io.EOF {
			return headers, nil
		}
		if err != nil {
			return nil, err
		}
		hdr = strings.TrimSuffix(hdr, "\r\n")
		if hdr == "" {
			return headers, nil
		}
		headerParts := strings.SplitN(hdr, ": ", 2)
		if len(headerParts) != 2 {
			return nil, fmt.Errorf("unexpected header format: %s", hdr)
		}
		headers[headerParts[0]] = headerParts[1]
	}
}

func (ps *PolicyServer) parseBody(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

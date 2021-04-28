package policyservice

import (
	"fmt"
	"io"
)

func (ps *PolicyServer) InvokePolicy(stream PolicyServer_InvokePolicyServer) error {
	resp := &InvokeReply{}
	first := true

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Info("Premature EOF")
			return fmt.Errorf("premature EOF")
		}
		if err != nil {
			log.Error(err.Error())
			return err
		}

		log.Infof("Got request: %v", req)

		// The headers are on the first request only.
		if first {
			if ps.Config.HeaderPolicy {
				log.Info("Setting headers")
				policyHeaders := map[string]string{
					"x-added": "new key",
				}
				for key, value := range req.Headers {
					policyHeaders[key] = value
				}
				resp.Headers = policyHeaders
				resp.Id = req.Id
			}
		}

		if ps.Config.BodyPolicy {
			if len(req.Body) > 0 {
				// SHOULD DO BUFFERING, CHUNKING AND WHAT NOT
				log.Info("Setting body")
				resp.Body = reverseBody(req.GetBody())
			} else {
				log.Info("No body")
			}
		}

		if req.EndOfStream {
			log.Info("Request finished")
			break
		}
		first = false
	}

	// Send the response
	resp.EndOfStream = true
	log.Infof("Responding %v", resp)
	return stream.Send(resp)
}

func reverseBody(input []byte) []byte {
	// -1 to preseve tailing null
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}

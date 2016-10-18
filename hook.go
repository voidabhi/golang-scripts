package handler

// https://developer.github.com/webhooks/
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Replace with your hook's secret
const secret = "shhhhh!!"

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

type HookContext struct {
	Signature string
	Event     string
	Id        string
	Payload   []byte
}

func ParseHook(secret []byte, req *http.Request) (*HookContext, error) {
	hc := HookContext{}

	if hc.Signature = req.Header.Get("x-hub-signature"); len(hc.Signature) == 0 {
		return nil, errors.New("No signature!")
	}

	if hc.Event = req.Header.Get("x-github-event"); len(hc.Event) == 0 {
		return nil, errors.New("No event!")
	}

	if hc.Id = req.Header.Get("x-github-delivery"); len(hc.Id) == 0 {
		return nil, errors.New("No event Id!")
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	if !verifySignature(secret, hc.Signature, body) {
		return nil, errors.New("Invalid signature")
	}

	hc.Payload = body

	return &hc, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {

	hc, err := ParseHook([]byte(secret), r)

	w.Header().Set("Content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed processing hook! ('%s')", err)
		io.WriteString(w, "{}")
		return
	}

	log.Printf("Received %s", hc.Event)

  // parse `hc.Payload` or do additional processing here

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
	return
}

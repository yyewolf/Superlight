package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// buildHeader : Builds the header for the request
func buildHeader(req *http.Request, body []byte) {
	req.Header.Set("client_id", config.ClientID)
	req.Header.Set("sign_method", "HMAC-SHA256")

	ts := fmt.Sprint(time.Now().UnixNano() / 1e6)
	req.Header.Set("t", ts)

	if token.Result.AccessToken != "" {
		req.Header.Set("access_token", token.Result.AccessToken)
	}

	sign := buildSign(req, body, ts)
	req.Header.Set("sign", sign)
}

// buildSign : Builds the sign for the request
func buildSign(req *http.Request, body []byte, t string) string {
	headers := getHeaderStr(req)
	urlStr := getUrlStr(req)
	contentSha256 := Sha256(body)
	stringToSign := req.Method + "\n" + contentSha256 + "\n" + headers + "\n" + urlStr
	signStr := config.ClientID + token.Result.AccessToken + t + stringToSign
	sign := strings.ToUpper(HmacSha256(signStr, config.Secret))
	return sign
}

// Sha256 : Returns the sha256 of the given string
func Sha256(data []byte) string {
	sha256Contain := sha256.New()
	sha256Contain.Write(data)
	return hex.EncodeToString(sha256Contain.Sum(nil))
}

// getUrlStr : Returns the url string
func getUrlStr(req *http.Request) string {
	url := req.URL.Path
	keys := make([]string, 0, 10)

	query := req.URL.Query()
	for key := range query {
		keys = append(keys, key)
	}
	if len(keys) > 0 {
		url += "?"
		sort.Strings(keys)
		for _, keyName := range keys {
			value := query.Get(keyName)
			url += keyName + "=" + value + "&"
		}
	}

	if url[len(url)-1] == '&' {
		url = url[:len(url)-1]
	}
	return url
}

// getHeaderStr : Returns the header string
func getHeaderStr(req *http.Request) string {
	signHeaderKeys := req.Header.Get("Signature-Headers")
	if signHeaderKeys == "" {
		return ""
	}
	keys := strings.Split(signHeaderKeys, ":")
	headers := ""
	for _, key := range keys {
		headers += key + ":" + req.Header.Get(key) + "\n"
	}
	return headers
}

// HmacSha256 : Returns the hmac sha256 of the given string
func HmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

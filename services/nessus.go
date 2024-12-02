package services

import (
	"NessusEssAutomation/config"
	"NessusEssAutomation/models"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Authenticate performs authentication and returns Nessus client
func Authenticate() (models.Nessus, error) {
	context, err := config.LoadNessusConfig()
	if err != nil {
		return models.Nessus{}, err
	}

	customTransport := (http.DefaultTransport.(*http.Transport))
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	context.HttpClient = &http.Client{Transport: customTransport}

	context, err = GetToken(context)
	if err != nil {
		return models.Nessus{}, err
	}

	keys, err := GetSessionKeys(context)
	if err != nil {
		return models.Nessus{}, err
	}

	return keys, nil
}

// GetToken retrieves the authentication token
func GetToken(context models.Nessus) (models.Nessus, error) {
	body := []byte(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, context.Username, context.Password))
	posturl := context.Url + "/session"
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		return models.Nessus{}, err
	}
	r.Header.Add("Content-Type", "application/json")
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return models.Nessus{}, err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseobj models.Session
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return models.Nessus{}, err
	}
	context.Token = responseobj.Token
	return context, nil
}

// GetSessionKeys retrieves the session keys for Nessus
func GetSessionKeys(context models.Nessus) (models.Nessus, error) {
	puturl := context.Url + "/session/keys"
	r, err := http.NewRequest("PUT", puturl, nil)
	if err != nil {
		return models.Nessus{}, err
	}
	token := fmt.Sprintf("token=%s", context.Token)
	r.Header.Add("X-Cookie", token)
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return models.Nessus{}, err
	}
	if res.StatusCode != 200 {
		return models.Nessus{}, errors.New("not able to authenticate with server")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Nessus{}, err
	}
	var responseobj models.SessionKeys
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return models.Nessus{}, err
	}
	context.AccessKey = responseobj.AccessKey
	context.SecretKey = responseobj.SecretKey
	return context, nil
}
func GetScan(_context models.Nessus) ([]models.Scan, error) {
	geturl := _context.Url + "/scans?"
	r, err := http.NewRequest("GET", geturl, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", _context.AccessKey, _context.SecretKey))

	res, err := _context.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {

		return nil, errors.New("not able to authenticate with server")
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseobj models.ScanResponse
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return nil, err
	}
	return responseobj.Scans, nil
}

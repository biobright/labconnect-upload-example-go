package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

var apiKey = ""
var basePath = "https://dbxdemo.biobright.xyz/api/"
var client = &http.Client{}

type Source struct {
	Id string
}

type SourcesResponse struct {
	Data []Source
}

// gets the ID of an agent or virtual folder
func getLLCAgentId(name string) (string, error) {
	respBytes, err := doLLCRequest("GET", "sources", [][]string{{"instrument_name", name}, {"include_instrument_name", "true"}}, nil, "")
	if err != nil {
		return "", err
	}
	var sourcesResonse SourcesResponse
	err = json.Unmarshal(respBytes, &sourcesResonse)
	if err != nil {
		return "", err
	}
	if len(sourcesResonse.Data) == 0 {
		return "", errors.New("No sources found for instrument name: " + name)
	}
	return sourcesResonse.Data[0].Id, nil
}

// create a folder in a virtual folder given the agent ID
func createVirtualFolder(folder_path string, agentId string) error {
	var err error = nil
	_, err = doLLCRequest("POST", "virtual-folder/folders", [][]string{{"folder_path", folder_path}, {"client_instance_id", agentId}}, nil, "")
	return err
}

// performs a multipart form data file upload
func uploadFileToVirtualFolder(fileHash string, agentId string, filePathIncludingFileName string, fileName string, fileReader io.Reader) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	err = w.WriteField("encrypted", "false")
	if err != nil {
		return err
	}
	err = w.WriteField("etag", fileHash)
	if err != nil {
		return err
	}
	err = w.WriteField("id", agentId)
	if err != nil {
		return err
	}
	err = w.WriteField("clientPath", filePathIncludingFileName)
	if err != nil {
		return err
	}
	err = w.WriteField("partSize", strconv.Itoa(10*1024*1024))
	if err != nil {
		return err
	}
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, fileReader)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	_, err = doLLCRequest("POST", "virtual-folder/upload", nil, &b, w.FormDataContentType())
	return err
}

// helper class for making requests to LLC
func doLLCRequest(method string, path string, queryParams [][]string, body io.Reader, formDataContentType string) ([]byte, error) {
	req, err := http.NewRequest(method, basePath+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	if formDataContentType != "" {
		req.Header.Add("Content-Type", formDataContentType)
	}
	finalQueryParams := req.URL.Query()
	if queryParams != nil {
		for _, param := range queryParams {
			finalQueryParams.Add(param[0], param[1])
		}
	}
	req.URL.RawQuery = finalQueryParams.Encode()
	//println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := io.ReadAll(resp.Body)
	respString := string(respBytes)
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Request failed with status code: %d, message %s", resp.StatusCode, respString)
	}
	//println(respString)
	return respBytes, nil
}

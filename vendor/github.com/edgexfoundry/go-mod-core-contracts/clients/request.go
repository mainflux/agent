/*******************************************************************************
 * Copyright 2019 Dell Inc.
 * Copyright 2019 Joan Duran
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/types"
)

// Helper method to get the body from the response after making the request
func getBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Helper method to make the request and return the response
func makeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

// GetRequest will make a GET request to the specified URL with the root URL retrieved by the URLClient prepended.
// It returns the body as a byte array if successful and an error otherwise.
func GetRequest(ctx context.Context, urlSuffix string, urlClient interfaces.URLClient) ([]byte, error) {
	urlPrefix, err := urlClient.Prefix()
	if err != nil {
		return nil, err
	}

	return GetRequestWithURL(ctx, urlPrefix+urlSuffix)
}

// GetRequestWithURL will make a GET request to the specified URL.
// It returns the body as a byte array if successful and an error otherwise.
func GetRequestWithURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c := NewCorrelatedRequest(ctx, req)
	resp, err := makeRequest(c.Request)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, types.ErrResponseNil{}
	}
	defer resp.Body.Close()

	bodyBytes, err := getBody(resp)
	if err != nil {
		return nil, err
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusAccepted) {
		return nil, types.NewErrServiceClient(resp.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}

// Helper method to make the count request
func CountRequest(ctx context.Context, urlSuffix string, urlClient interfaces.URLClient) (int, error) {
	// do not get URLPrefix here since GetRequest does it
	data, err := GetRequest(ctx, urlSuffix, urlClient)
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Helper method to make the post JSON request and return the body
func PostJSONRequest(
	ctx context.Context,
	urlSuffix string,
	data interface{},
	urlClient interfaces.URLClient) (string, error) {

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	ctx = context.WithValue(ctx, ContentType, ContentTypeJSON)

	// do not get URLPrefix here since PostRequest does it
	return PostRequest(ctx, urlSuffix, jsonStr, urlClient)
}

// PostJSONRequestWithURL will make a POST request to the specified URL with the object passed in
// marshaled into a JSON formatted byte array.
// It returns the body on success and an error otherwise.
func PostJSONRequestWithURL(ctx context.Context, url string, data interface{}) (string, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	ctx = context.WithValue(ctx, ContentType, ContentTypeJSON)

	return PostRequestWithURL(ctx, url, jsonStr)
}

// Helper method to make the post request and return the body
func PostRequest(ctx context.Context, urlSuffix string, data []byte, urlClient interfaces.URLClient) (string, error) {
	urlPrefix, err := urlClient.Prefix()
	if err != nil {
		return "", err
	}

	return PostRequestWithURL(ctx, urlPrefix+urlSuffix, data)
}

// PostRequestWithURL will make a POST request to the specified URL.
// It returns the body as a byte array if successful and an error otherwise.
func PostRequestWithURL(ctx context.Context, url string, data []byte) (string, error) {
	content := FromContext(ctx, ContentType)
	if content == "" {
		content = ContentTypeJSON
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set(ContentType, content)

	c := NewCorrelatedRequest(ctx, req)
	resp, err := makeRequest(c.Request)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", types.ErrResponseNil{}
	}
	defer resp.Body.Close()

	bodyBytes, err := getBody(resp)
	if err != nil {
		return "", err
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusAccepted) {
		return "", types.NewErrServiceClient(resp.StatusCode, bodyBytes)
	}

	bodyString := string(bodyBytes)
	return bodyString, nil
}

// Helper method to make a post request in order to upload a file and return the request body
func UploadFileRequest(
	ctx context.Context,
	urlSuffix string,
	filePath string,
	urlClient interfaces.URLClient) (string, error) {

	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Create multipart/form-data request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(formFileWriter, bytes.NewReader(fileContents))
	if err != nil {
		return "", err
	}
	writer.Close()

	urlPrefix, err := urlClient.Prefix()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, urlPrefix+urlSuffix, body)
	if err != nil {
		return "", err
	}
	req.Header.Add(ContentType, writer.FormDataContentType())

	c := NewCorrelatedRequest(ctx, req)
	resp, err := makeRequest(c.Request)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", types.ErrResponseNil{}
	}
	defer resp.Body.Close()

	bodyBytes, err := getBody(resp)
	if err != nil {
		return "", err
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusAccepted) {
		return "", types.NewErrServiceClient(resp.StatusCode, bodyBytes)
	}

	bodyString := string(bodyBytes)
	return bodyString, nil
}

// Helper method to make the update request
func UpdateRequest(ctx context.Context, urlSuffix string, data interface{}, urlClient interfaces.URLClient) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// do not get URLPrefix here since PutRequest does it
	_, err = PutRequest(ctx, urlSuffix, jsonStr, urlClient)
	return err
}

// Helper method to make the put request
func PutRequest(ctx context.Context, urlSuffix string, body []byte, urlClient interfaces.URLClient) (string, error) {
	var err error
	var req *http.Request

	urlPrefix, err := urlClient.Prefix()
	if err != nil {
		return "", err
	}
	if body != nil {
		req, err = http.NewRequest(http.MethodPut, urlPrefix+urlSuffix, bytes.NewReader(body))
		if err != nil {
			return "", err
		}

		content := FromContext(ctx, ContentType)
		if content == "" {
			content = ContentTypeJSON
		}
		req.Header.Set(ContentType, content)
	} else {
		req, err = http.NewRequest(http.MethodPut, urlPrefix+urlSuffix, nil)
	}
	if err != nil {
		return "", err
	}

	c := NewCorrelatedRequest(ctx, req)
	resp, err := makeRequest(c.Request)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", types.ErrResponseNil{}
	}
	defer resp.Body.Close()

	bodyBytes, err := getBody(resp)
	if err != nil {
		return "", err
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusAccepted) {
		return "", types.NewErrServiceClient(resp.StatusCode, bodyBytes)
	}

	bodyString := string(bodyBytes)
	return bodyString, nil
}

// Helper method to make the delete request
func DeleteRequest(ctx context.Context, urlSuffix string, urlClient interfaces.URLClient) error {
	urlPrefix, err := urlClient.Prefix()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, urlPrefix+urlSuffix, nil)
	if err != nil {
		return err
	}

	c := NewCorrelatedRequest(ctx, req)
	resp, err := makeRequest(c.Request)
	if err != nil {
		return err
	}
	if resp == nil {
		return types.ErrResponseNil{}
	}
	defer resp.Body.Close()

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusAccepted) {
		bodyBytes, err := getBody(resp)
		if err != nil {
			return err
		}

		return types.NewErrServiceClient(resp.StatusCode, bodyBytes)
	}

	return nil
}

// CorrelatedRequest is a wrapper type for use in managing Correlation IDs during service to service API calls.
type CorrelatedRequest struct {
	*http.Request
}

// NewCorrelatedRequest will add the Correlation ID header to the supplied request. If no Correlation ID header is
// present in the supplied context, one will be created along with a value.
func NewCorrelatedRequest(ctx context.Context, req *http.Request) CorrelatedRequest {
	c := CorrelatedRequest{Request: req}
	correlation := FromContext(ctx, CorrelationHeader)
	if len(correlation) == 0 {
		correlation = uuid.New().String()
	}
	c.Header.Set(CorrelationHeader, correlation)
	return c
}

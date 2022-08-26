package docatl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Docat struct {
	Host   string
	ApiKey string
}

type ProjectClaim struct {
	Token string
}

func (docat *Docat) Post(project string, version string, docsPath string) error {
	file, err := os.Open(docsPath)
	if err != nil {
		return fmt.Errorf("unable to upload documentation because it isn't accessible locally at '%s'", docsPath)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("unable to upload documentation because cannot create form file for file '%s'", docsPath)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("unable to upload documentation because cannot copy file content of file '%s' to request", docsPath)
	}
	writer.Close()

	apiUrl := fmt.Sprintf("%s/api/%s/%s", docat.Host, project, version)
	request, err := http.NewRequest(http.MethodPost, apiUrl, body)
	if err != nil {
		return fmt.Errorf("unable to upload documentation because cannot create POST request: %s", err)
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	if docat.ApiKey != "" {
		request.Header.Add("Docat-Api-Key", docat.ApiKey)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("unable to upload documentation: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to upload documentation and read it's response (status code: %d", response.StatusCode)
		}
		return fmt.Errorf("unable to upload documentation: (status code: %d) %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

func (docat *Docat) Delete(project string, version string) error {
	apiUrl := fmt.Sprintf("%s/api/%s/%s", docat.Host, project, version)
	request, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to delete documentation because cannot create DELETE request: %s", err)
	}
	request.Header.Add("Docat-Api-Key", docat.ApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("unable to delete documentation because request failed: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to delete documentation and read it's response (status code: %d", response.StatusCode)
		}
		return fmt.Errorf("unable to delete documentation: (status code: %d) %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

func (docat *Docat) Claim(project string) (ProjectClaim, error) {
	apiUrl := fmt.Sprintf("%s/api/%s/claim", docat.Host, project)
	response, err := http.Get(apiUrl)
	if err != nil {
		return ProjectClaim{}, fmt.Errorf("unable to claim project because request failed: %s", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ProjectClaim{}, fmt.Errorf("unable to claim project and read it's response (status code: %d", response.StatusCode)
	}

	if response.StatusCode != http.StatusCreated {
		return ProjectClaim{}, fmt.Errorf("unable to claim project: (status code: %d) %s", response.StatusCode, string(bodyBytes))
	}

	var claim ProjectClaim
	if err = json.Unmarshal(bodyBytes, &claim); err != nil {
		return ProjectClaim{}, fmt.Errorf("unable to claim project, because cannot unmarshal response from server: %s", string(bodyBytes))
	}
	return claim, nil
}

func (docat *Docat) Tag(project string, version string, tag string) error {
	apiUrl := fmt.Sprintf("%s/api/%s/%s/tags/%s", docat.Host, project, version, tag)
	request, err := http.NewRequest(http.MethodPut, apiUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to tag documentation because cannot create PUT request: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("unable to tag documentation because request failed: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to tag documentation and read it's response (status code: %d", response.StatusCode)
		}
		return fmt.Errorf("unable to tag documentation: (status code: %d) %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

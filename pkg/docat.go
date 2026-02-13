package docatl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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
	defer func() { _ = file.Close() }()

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
	if err = writer.Close(); err != nil {
		return fmt.Errorf("unable to upload documentation because cannot close multipart writer: %s", err)
	}

	apiUrl, err := url.JoinPath(docat.Host, "api", project, version)
	if err != nil {
		return fmt.Errorf("unable to upload documentation because cannot create an url for host: %s error: %s", docat.Host, err)
	}

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
	defer func() { _ = response.Body.Close() }()

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
	apiUrl, err := url.JoinPath(docat.Host, "api", project, version)
	if err != nil {
		return fmt.Errorf("unable to delete documentation because cannot create an url for host: %s error: %s", docat.Host, err)
	}

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
	defer func() { _ = response.Body.Close() }()

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
	apiUrl, err := url.JoinPath(docat.Host, "api", project, "claim")
	if err != nil {
		return ProjectClaim{}, fmt.Errorf("unable to claim project because cannot create an url for host: %s error: %s", docat.Host, err)
	}

	response, err := http.Get(apiUrl)
	if err != nil {
		return ProjectClaim{}, fmt.Errorf("unable to claim project because request failed: %s", err)
	}
	defer func() { _ = response.Body.Close() }()

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
	apiUrl, err := url.JoinPath(docat.Host, "api", project, version, "tags", tag)
	if err != nil {
		return fmt.Errorf("unable to tag documentation because cannot create an url for host: %s error: %s", docat.Host, err)
	}

	request, err := http.NewRequest(http.MethodPut, apiUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to tag documentation because cannot create PUT request: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("unable to tag documentation because request failed: %s", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to tag documentation and read it's response (status code: %d", response.StatusCode)
		}
		return fmt.Errorf("unable to tag documentation: (status code: %d) %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

func (docat *Docat) PushIcon(project string, iconPath string) error {
	file, err := os.Open(iconPath)
	if err != nil {
		return fmt.Errorf("unable to upload icon because the path '%s' does not exist", iconPath)
	}
	defer func() { _ = file.Close() }()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("unable to upload icon because create form file for file '%s' failed", iconPath)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("unable to upload icon because copy file content from file '%s' to request failed", iconPath)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("unable to upload icon because cannot close multipart writer: %s", err)
	}

	apiUrl, err := url.JoinPath(docat.Host, "api", project, "icon")
	if err != nil {
		return fmt.Errorf("unable to upload icon because creating an url for host: %s failed with error: %s", docat.Host, err)
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		return fmt.Errorf("unable to upload icon because creating the POST request failed: %s", err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	if docat.ApiKey != "" {
		request.Header.Add("Docat-Api-Key", docat.ApiKey)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return fmt.Errorf("unable to upload icon because the request failed: %s", err)
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("unable to upload icon and read it's response (status code: %d", response.StatusCode)
	}

	return fmt.Errorf("unable to upload icon: (status code: %d) %s", response.StatusCode, string(bodyBytes))
}

func (docat *Docat) Rename(project string, newName string) error {
	apiUrl, err := url.JoinPath(docat.Host, "api", project, "rename", newName)

	if err != nil {
		return fmt.Errorf("unable to rename project because creating an url failed for host: %s error: %s", docat.Host, err)
	}

	request, err := http.NewRequest(http.MethodPut, apiUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to rename project because creating PUT request failed: %s", err)

	}

	if docat.ApiKey != "" {
		request.Header.Add("Docat-Api-Key", docat.ApiKey)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return fmt.Errorf("unable to rename project because the request failed: %s", err)
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to rename project and read the response (status code: %d", response.StatusCode)
	}
	return fmt.Errorf("unable to rename project: (status code: %d) %s", response.StatusCode, string(bodyBytes))
}

func (docat *Docat) HideOrShowVersion(project string, version string, hide bool) error {
	var hideOrShow string
	if hide {
		hideOrShow = "hide"
	} else {
		hideOrShow = "show"
	}

	apiUrl, err := url.JoinPath(docat.Host, "api", project, version, hideOrShow)
	if err != nil {
		return fmt.Errorf("unable to %s version because creating an url failed for host: %s error: %s", hideOrShow, docat.Host, err)
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to %s version because creating POST request failed: %s", hideOrShow, err)
	}
	if docat.ApiKey != "" {
		request.Header.Add("Docat-Api-Key", docat.ApiKey)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("unable to %s version because request failed: %s", hideOrShow, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to %s version and read it's response (status code: %d", hideOrShow, response.StatusCode)
		}
		return fmt.Errorf("unable to %s version: (status code: %d) %s", hideOrShow, response.StatusCode, string(bodyBytes))
	}

	return nil
}

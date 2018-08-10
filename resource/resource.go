package resource

import (
	"encoding/json"
	"github.com/Piszmog/cloudconfigclient/client"
	"github.com/pkg/errors"
)

const (
	defaultApplicationName    = "default"
	defaultApplicationProfile = "default"
)

type Resource interface {
	GetFile(directory string, file string, interfaceType *interface{}) error
	GetFileFromBranch(branch string, directory string, file string, interfaceType *interface{}) error
}

type Client struct {
	configClient *client.ConfigClient
}

func CreateResourceClient(configClient *client.ConfigClient) *Client {
	return &Client{configClient: configClient}
}

func (client *Client) GetFile(directory string, file string, interfaceType interface{}) error {
	fileFound := false
	for _, configClient := range client.configClient.Clients {
		resp, err := configClient.Get(defaultApplicationName, defaultApplicationProfile, directory, file+"?useDefaultLabel=true")
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve file from %s",
				configClient.GetFullUrl(defaultApplicationName, defaultApplicationProfile, directory, file+"?useDefaultLabel=true"))
		}
		if resp.StatusCode != 200 {
			return errors.Errorf("server responded with status code %d from url %s",
				resp.StatusCode,
				configClient.GetFullUrl(defaultApplicationName, defaultApplicationProfile, directory, file+"?useDefaultLabel=true"))
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(interfaceType)
		resp.Body.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to decode response from url %s",
				configClient.GetFullUrl(defaultApplicationName, defaultApplicationProfile, directory, file+"?useDefaultLabel=true"))
		}
		fileFound = true
	}
	if !fileFound {
		return errors.Errorf("failed to find file %s in the Config Server", file)
	}
	return nil
}

func (client *Client) GetFileFromBranch(branch string, directory string, file string, interfaceType interface{}) error {
	fileFound := false
	for _, configClient := range client.configClient.Clients {
		resp, err := configClient.Get(defaultApplicationName, defaultApplicationProfile, branch, directory, file)
		if resp != nil && resp.StatusCode == 404 {
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve file from %s",
				configClient.GetFullUrl(defaultApplicationName, defaultApplicationProfile, branch, directory, file))
		}
		if resp.StatusCode != 200 {
			return errors.Errorf("server responded with status code %d from url %s", resp.StatusCode,
				configClient.GetFullUrl(defaultApplicationName, defaultApplicationProfile, branch, directory, file))
		}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(interfaceType)
		resp.Body.Close()
		if err != nil {
			return errors.Wrapf(err, "failed to decode response from url %s",
				configClient.GetFullUrl(defaultApplicationName, defaultApplicationProfile, branch, directory, file))
		}
		fileFound = true
	}
	if !fileFound {
		return errors.Errorf("failed to find file %s in the Config Server", file)
	}
	return nil
}

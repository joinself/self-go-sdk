// Copyright 2020 Self Group Ltd. All Rights Reserved.

package request

import (
	"encoding/json"
)

type RestTransport interface {
	Get(path string) ([]byte, error)
	Post(path string, ctype string, data []byte) ([]byte, error)
	BuildURL(path string) string
}

type Client struct {
	selfID   string
	deviceID string
	api      RestTransport
}

type Config struct {
	SelfID   string
	DeviceID string
	API      RestTransport
}

func New(config Config) *Client {
	return &Client{
		selfID:   config.SelfID,
		deviceID: config.DeviceID,
		api:      config.API,
	}
}

func (c *Client) SetAPI(api RestTransport) {
	c.api = api
}

// builds a list of all devices associated with a list of identities
func (c *Client) FormatRecipients(recipients []string) ([]string, error) {
	devices := make([]string, 0)
	for _, selfID := range recipients {
		dds, err := c.getDevices(selfID)
		if err != nil {
			return nil, err
		}

		for i := range dds {
			// if is not the current device
			if selfID != c.selfID && dds[i] != c.deviceID {
				devices = append(devices, selfID+":"+dds[i])
			}
		}
	}

	return devices, nil
}

func (c *Client) getDevices(selfID string) ([]string, error) {
	var resp []byte
	var err error

	if len(selfID) > 11 {
		resp, err = c.api.Get("/v1/apps/" + selfID + "/devices")
	} else {
		resp, err = c.api.Get("/v1/identities/" + selfID + "/devices")
	}
	if err != nil {
		return nil, err
	}

	var devices []string
	err = json.Unmarshal(resp, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

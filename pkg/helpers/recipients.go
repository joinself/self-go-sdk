package helpers

import "encoding/json"

type RestTransport interface {
	Get(path string) ([]byte, error)
}

func PrepareRecipients(recipients, excludedDevices []string, api RestTransport) ([]string, error) {
	devices := make([]string, 0)
	for _, sID := range recipients {
		dds, err := getDevices(api, sID)
		if err != nil {
			return nil, err
		}

		for i := range dds {
			dd := sID + ":" + dds[i]
			if !stringInSlice(dd, excludedDevices) {
				devices = append(devices, dd)
			}
		}
	}

	return devices, nil
}

func getDevices(api RestTransport, selfID string) ([]string, error) {
	var resp []byte
	var err error

	if len(selfID) > 11 {
		resp, err = api.Get("/v1/apps/" + selfID + "/devices")
	} else {
		resp, err = api.Get("/v1/identities/" + selfID + "/devices")
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

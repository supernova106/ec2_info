package models

type AWSPrice struct {
	Config struct {
		Currencies []string `json:"currencies"`
		Rate       string   `json:"rate"`
		Regions    []struct {
			InstanceTypes []struct {
				Sizes []struct {
					ECU          string `json:"ECU"`
					MemoryGiB    string `json:"memoryGiB"`
					Size         string `json:"size"`
					StorageGB    string `json:"storageGB"`
					VCPU         string `json:"vCPU"`
					ValueColumns []struct {
						Name   string `json:"name"`
						Prices struct {
							USD string `json:"USD"`
						} `json:"prices"`
					} `json:"valueColumns"`
				} `json:"sizes"`
				Type string `json:"type"`
			} `json:"instanceTypes"`
			Region string `json:"region"`
		} `json:"regions"`
		ValueColumns []string `json:"valueColumns"`
	} `json:"config"`
	Vers float64 `json:"vers"`
}

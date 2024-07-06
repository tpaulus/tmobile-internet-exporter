package main

type GatewayStatus struct {
	Device Device `json:"device,omitempty"`
	Time   Time   `json:"time,omitempty"`
	Signal Signal `json:"signal,omitempty"`
}
type Device struct {
	Type            string `json:"type,omitempty"`
	Manufacturer    string `json:"manufacturer,omitempty"`
	ManufacturerOUI string `json:"manufacturerOUI,omitempty"`
	Name            string `json:"name,omitempty"`
	FriendlyName    string `json:"friendlyName,omitempty"`
	IsEnabled       bool   `json:"isEnabled,omitempty"`
	Index           int    `json:"index,omitempty"`
	IsMeshSupported bool   `json:"isMeshSupported,omitempty"`
	HardwareVersion string `json:"hardwareVersion,omitempty"`
	SoftwareVersion string `json:"softwareVersion,omitempty"`
	Model           string `json:"model,omitempty"`
	Serial          string `json:"serial,omitempty"`
	MacID           string `json:"macId,omitempty"`
	UpdateState     string `json:"updateState,omitempty"`
	Role            string `json:"role,omitempty"`
}
type DaylightSavings struct {
	IsUsed bool `json:"isUsed,omitempty"`
}
type Time struct {
	UpTime          int             `json:"upTime,omitempty"`
	LocalTime       int             `json:"localTime,omitempty"`
	LocalTimeZone   string          `json:"localTimeZone,omitempty"`
	DaylightSavings DaylightSavings `json:"daylightSavings,omitempty"`
}
type Generic struct {
	Apn          string `json:"apn,omitempty"`
	Roaming      bool   `json:"roaming,omitempty"`
	Registration string `json:"registration,omitempty"`
	HasIPv6      bool   `json:"hasIPv6,omitempty"`
}
type FourG struct {
	ENBID int      `json:"eNBID,omitempty"`
	Cid   int      `json:"cid,omitempty"`
	Sinr  int      `json:"sinr,omitempty"`
	Rsrp  int      `json:"rsrp,omitempty"`
	Rsrq  int      `json:"rsrq,omitempty"`
	Rssi  int      `json:"rssi,omitempty"`
	Bands []string `json:"bands,omitempty"`
	Bars  int      `json:"bars,omitempty"`
}
type FiveG struct {
	GNBID int      `json:"gNBID,omitempty"`
	Cid   int      `json:"cid,omitempty"`
	Sinr  int      `json:"sinr,omitempty"`
	Rsrp  int      `json:"rsrp,omitempty"`
	Rsrq  int      `json:"rsrq,omitempty"`
	Rssi  int      `json:"rssi,omitempty"`
	Bands []string `json:"bands,omitempty"`
	Bars  int      `json:"bars,omitempty"`
}
type Signal struct {
	Generic Generic `json:"generic,omitempty"`
	FourG   FourG   `json:"4g,omitempty"`
	FiveG   FiveG   `json:"5g,omitempty"`
}

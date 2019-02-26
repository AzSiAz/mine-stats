package models

type MinecraftStatus struct {
	Hostname    string `json:"host"`
	Port        uint16 `json:"port"`
	Description struct {
		Motd string `json:"text"`
	} `json:"description"`
	ServerVersion struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	PlayerInfo struct {
		Max     int64 `json:"max"`
		Current int64 `json:"online"`
		Players []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	ModInfo struct {
		Type    string
		ModList []struct {
			ModID   string `json:"modid"`
			Version string `json:"version"`
		}
	} `json:"modinfo"`
}

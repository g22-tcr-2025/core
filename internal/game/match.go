package game

type MatchData struct {
	PUsername string  `json:"pusername"`
	PLevel    float64 `json:"plevel"`
	PEXP      float64 `json:"pexp"`
	PMana     float64 `json:"pmana"`
	PTroops   []Troop `json:"ptroops"`
	PTowers   []Tower `json:"ptowers"`

	OUsername string  `json:"ousername"`
	OLevel    float64 `json:"olevel"`
	OEXP      float64 `json:"oexp"`
	OMana     float64 `json:"omana"`
	OTroops   []Troop `json:"otroops"`
	OTowers   []Tower `json:"otowers"`
}

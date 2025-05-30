package game

type MatchData struct {
	PUsername      string  `json:"pusername"`
	PLevelMetadata float64 `json:"plevel_metadata"`
	PEXPMetadata   float64 `json:"pexp_metadata"`
	PLevel         float64 `json:"plevel"`
	PEXP           float64 `json:"pexp"`
	PMana          float64 `json:"pmana"`
	PTroops        []Troop `json:"ptroops"`
	PTowers        []Tower `json:"ptowers"`

	OUsername      string  `json:"ousername"`
	OLevelMetadata float64 `json:"olevel_metadata"`
	OEXPMetadata   float64 `json:"oexp_metadata"`
	OLevel         float64 `json:"olevel"`
	OEXP           float64 `json:"oexp"`
	OMana          float64 `json:"omana"`
	OTroops        []Troop `json:"otroops"`
	OTowers        []Tower `json:"otowers"`
}

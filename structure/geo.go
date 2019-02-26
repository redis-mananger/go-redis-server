package structure

type Geo struct {
    Values []map[string]*GeoValue `json:"values"`
}

type GeoValue struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

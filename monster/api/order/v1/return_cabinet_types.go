package v1

// ReturnCabinetInfo represents a cabinet available for returning a power bank.
// Non-proto type used for the HTTP JSON API.
type ReturnCabinetInfo struct {
	CabinetID      int64   `json:"cabinetId"`
	CabinetNo      string  `json:"cabinetNo"`
	StationID      int64   `json:"stationId"`
	StationName    string  `json:"stationName"`
	StationAddress string  `json:"stationAddress"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Distance       float64 `json:"distance"`
	EmptySlotCount int32   `json:"emptySlotCount"`
	TotalSlots     int32   `json:"totalSlots"`
}

// ListReturnCabinetsRequest is the request for listing cabinets with empty slots.
// Query params: latitude, longitude, radius_meters (optional).
type ListReturnCabinetsRequest struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters int32   `json:"radiusMeters"`
}

// ListReturnCabinetsResponse is the response for the return cabinet list.
type ListReturnCabinetsResponse struct {
	List []*ReturnCabinetInfo `json:"list"`
}

// SearchCabinetRequest is the request for searching a cabinet by number.
type SearchCabinetRequest struct {
	CabinetNo string `json:"cabinetNo"`
}

// SearchCabinetResponse is the response for cabinet search.
type SearchCabinetResponse struct {
	Cabinet *ReturnCabinetInfo `json:"cabinet"`
}

package mrbilit

type Response struct {
	Trains        []*Trains    `json:"trains"`
	Filters       Filters      `json:"filters"`
	ToLocation    ToLocation   `json:"toLocation"`
	FromLocation  FromLocation `json:"fromLocation"`
	ContentPath   string       `json:"contentPath"`
	Meta          Meta         `json:"meta"`
	TrainPackages []any        `json:"trainPackages"`
}
type CancellationTermEntries struct {
	Amount        any    `json:"amount"`
	PercentAmount int    `json:"percentAmount"`
	Title         string `json:"title"`
}
type WagonServices struct {
	HasService bool   `json:"hasService"`
	Title      string `json:"title"`
}
type Classes struct {
	ID                      int                       `json:"id"`
	Price                   int                       `json:"price"`
	Capacity                int                       `json:"capacity"`
	WagonName               string                    `json:"wagonName"`
	AirConditioning         bool                      `json:"airConditioning"`
	Media                   bool                      `json:"media"`
	IsCompartment           bool                      `json:"isCompartment"`
	CompartmentCapacity     int                       `json:"compartmentCapacity"`
	Discount                int                       `json:"discount"`
	WagonID                 int                       `json:"wagonID"`
	CapacityString          string                    `json:"capacityString"`
	HasHotel                bool                      `json:"hasHotel"`
	IsAvailable             bool                      `json:"isAvailable"`
	CancellationTerms       string                    `json:"cancellationTerms"`
	ReservationAvailable    bool                      `json:"reservationAvailable"`
	PngLogoPath             string                    `json:"pngLogoPath"`
	SvgLogoPath             string                    `json:"svgLogoPath"`
	MinPersons              int                       `json:"minPersons"`
	Owner                   int                       `json:"owner"`
	OwnerName               string                    `json:"ownerName"`
	HasCompartmentCapacity  bool                      `json:"hasCompartmentCapacity"`
	HotelID                 any                       `json:"hotelId"`
	HotelName               any                       `json:"hotelName"`
	Class                   any                       `json:"class"`
	RoundTrip               bool                      `json:"roundTrip"`
	HasSpecialDiscount      bool                      `json:"hasSpecialDiscount"`
	Score                   any                       `json:"score"`
	CancellationTermEntries []CancellationTermEntries `json:"cancellationTermEntries"`
	SuggestedServices       []any                     `json:"suggestedServices"`
	WagonServices           []WagonServices           `json:"wagonServices"`
}
type Prices struct {
	SellType int       `json:"sellType"`
	Classes  []Classes `json:"classes"`
}
type Trains struct {
	ID                int       `json:"id"`
	From              int       `json:"from"`
	To                int       `json:"to"`
	FromName          string    `json:"fromName"`
	ToName            string    `json:"toName"`
	TrainNumber       int       `json:"trainNumber"`
	DepartureTime     string    `json:"departureTime"`
	ArrivalTime       string    `json:"arrivalTime"`
	Provider          int       `json:"provider"`
	ProviderName      string    `json:"providerName"`
	CorporationID     int       `json:"corporationID"`
	CorporationIds    []int     `json:"corporationIds"`
	CorporationName   string    `json:"corporationName"`
	Weekday           string    `json:"weekday"`
	DateString        string    `json:"dateString"`
	ArrivalDateString string    `json:"arrivalDateString"`
	FromCache         bool      `json:"fromCache"`
	Cancellable       bool      `json:"cancellable"`
	IsForeign         bool      `json:"isForeign"`
	Prices            []*Prices `json:"prices"`
	Score             any       `json:"score"`
}
type Corporations struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Ids  []int  `json:"ids"`
}
type WagonTypes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Filters struct {
	Corporations []Corporations `json:"corporations"`
	MinPrice     int            `json:"minPrice"`
	MaxPrice     int            `json:"maxPrice"`
	WagonTypes   []WagonTypes   `json:"wagonTypes"`
}
type ToLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Title     string  `json:"title"`
}
type FromLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Title     string  `json:"title"`
}
type Meta struct {
	AlwaysHasTrain bool `json:"alwaysHasTrain"`
	RouteNotFound  bool `json:"routeNotFound"`
}

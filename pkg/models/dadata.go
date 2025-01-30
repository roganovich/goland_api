package models


type AddressRequest struct {
	Query string `json:"query"`
}

type AddressSuggestion struct {
	Value string `json:"value"`
}

type AddressResponse struct {
	Suggestions []AddressSuggestion `json:"suggestions"`
}

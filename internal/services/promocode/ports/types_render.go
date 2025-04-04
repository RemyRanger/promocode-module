package ports

import "net/http"

// Render : Pre-processing before a response is marshalled and sent across the wire
func (a Promocode) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Render : Pre-processing before a response is marshalled and sent across the wire
func (a PromocodeValidationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind : just a post-process after a decode..
func (na *PromocodeIn) Bind(r *http.Request) error {
	return nil
}

// Bind : just a post-process after a decode..
func (na *PromocodeValidation) Bind(r *http.Request) error {
	return nil
}

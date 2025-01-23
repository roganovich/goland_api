package models

type ErrorResponse struct {
	StatusCode 	int    						`json:"statusCode"`
	Message    	string 						`json:"message"`
	Errors	   	[]ValidationErrorResponse	`json:"errors"`

}

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
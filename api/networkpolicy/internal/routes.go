package internal

import (
	"net/http"

	"github.com/tedsuo/rata"
)

const (
	PostAddPolicies = "AddPolicies"
)

// Routes is a list of routes used by the rata library to construct request
// URLs.
var Routes = rata.Routes{
	{Path: "/policies", Method: http.MethodPost, Name: PostAddPolicies},
}

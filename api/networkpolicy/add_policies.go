package networkpolicy

import (
	"bytes"
	"encoding/json"

	"code.cloudfoundry.org/cli/api/networkpolicy/internal"
)

func (client Client) AddPolicies(policies []Policy) error {
	// _ := c.JsonClient.Do("POST", "/policies", reqPolicies, nil, token)

	body := struct {
		Policies []Policy `json:"policies"`
	}{Policies: policies}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := client.newRequest(requestOptions{
		RequestName: internal.PostAddPolicies,
		Body:        bytes.NewBuffer(bodyBytes),
	})
	if err != nil {
		return err
	}

	response := Response{
		Result: nil,
	}

	err = client.connection.Make(request, &response)
	if err != nil {
		return err
	}

	return nil
}

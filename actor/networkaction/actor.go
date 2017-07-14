package networkaction

type Actor struct {
	PolicyClient NetworkPolicyClient
}

func NewActor(networkPolicyClient NetworkPolicyClient) *Actor {
	return &Actor{
		PolicyClient: networkPolicyClient,
	}
}

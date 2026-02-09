package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		PolicyMap: []Policy{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	policyIndexMap := make(map[string]struct{})

	for _, elem := range gs.PolicyMap {
		index := fmt.Sprint(elem.PolicyId)
		if _, ok := policyIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for policy")
		}
		policyIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}

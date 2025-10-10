package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		RateMap: []Rate{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	rateIndexMap := make(map[string]struct{})

	for _, elem := range gs.RateMap {
		index := fmt.Sprint(elem.Symbol)
		if _, ok := rateIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for rate")
		}
		rateIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}

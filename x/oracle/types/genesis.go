package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		PriceMap: []Price{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	priceIndexMap := make(map[string]struct{})

	for _, elem := range gs.PriceMap {
		index := fmt.Sprint(elem.Symbol)
		if _, ok := priceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for price")
		}
		priceIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}

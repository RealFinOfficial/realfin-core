package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		AssetMap: []Asset{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	assetIndexMap := make(map[string]struct{})

	for _, elem := range gs.AssetMap {
		index := fmt.Sprint(elem.Symbol)
		if _, ok := assetIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for asset")
		}
		assetIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}

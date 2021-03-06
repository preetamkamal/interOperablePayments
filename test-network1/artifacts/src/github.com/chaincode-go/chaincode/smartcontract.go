package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID          string `json:"ID"`
	Issuerid    string `json:"issuerid"`
	IssuerName  string `json:"issuerName"`
	Owner       string `json:"owner"`
	Value       int    `json:"value"`
	State       string `json:"state"`
	Category    string `json:"category"`
	AssetName   string `json:"assetName"`
	Account     string `json:"account"`
	Amount      string `json:"amount"`
	Destination string `json:"destination"`
	Hash        string `json:"hash"`
}

// InitLedger adds a base set of assets to the ledger
// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	assets := []Asset{
// 		{ID: "asset1", Issuerid: "Axis7051", IssuerName: "Axis", Owner: "Axis", Value: 500, State: "Issued", Category: "Machines"},
// 		{ID: "asset2", Issuerid: "Hdfc8051", IssuerName: "Hdfc", Owner: "Hdfc", Value: 3300, State: "Issued", Category: "Machines"},
// 		{ID: "asset3", Issuerid: "Axis7051", IssuerName: "Axis", Owner: "Axis", Value: 300, State: "Issued", Category: "Machines"},
// 		{ID: "asset4", Issuerid: "Bofs9051", IssuerName: "Bofs", Owner: "Bofs", Value: 800, State: "Issued", Category: "Furniture"},
// 		{ID: "asset5", Issuerid: "Hdfc8051", IssuerName: "Hdfc", Owner: "Hdfc", Value: 900, State: "Issued", Category: "Furniture"},
// 	}

// 	for _, asset := range assets {
// 		assetJSON, err := json.Marshal(asset)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState(asset.ID, assetJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %v", err)
// 		}
// 	}

// 	return nil
// }

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, issueID string, issueName string, owner string, appraisedValue int, cat string, assetName string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:          id,
		Issuerid:    issueID,
		IssuerName:  issueName,
		Owner:       owner,
		Value:       appraisedValue,
		State:       "Issued",
		Category:    cat,
		AssetName:   assetName,
		Account:     "",
		Amount:      "",
		Destination: "",
		Hash:        "",
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset *Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string, amount string, account string, destination string, hash string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	asset.Amount = amount
	asset.Account = account
	asset.Destination = destination
	asset.Hash = hash

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAsstByCat retrives all the assets according to category.
func (s *SmartContract) GetAssetByCat(ctx contractapi.TransactionContextInterface, cat string) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset *Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		if strings.ToLower(asset.Category) == strings.ToLower(cat) {
			assets = append(assets, asset)
		}
	}

	return assets, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset *Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what descriptions up a car
type Car struct {
	Description   string `json:"description"`
	ProductID  string `json:"productID"`
	Colour string `json:"colour"`
	Status  string `json:"status"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{Description: "Toyota", ProductID: "Prius", Colour: "blue", Status: "Tomoko"},
		Car{Description: "Ford", ProductID: "Mustang", Colour: "red", Status: "Brad"},
		Car{Description: "Hyundai", ProductID: "Tucson", Colour: "green", Status: "Jin Soo"},
		Car{Description: "Volkswagen", ProductID: "Passat", Colour: "yellow", Status: "Max"},
		Car{Description: "Tesla", ProductID: "S", Colour: "black", Status: "Adriana"},
		Car{Description: "Peugeot", ProductID: "205", Colour: "purple", Status: "Michel"},
		Car{Description: "Chery", ProductID: "S22L", Colour: "white", Status: "Aarav"},
		Car{Description: "Fiat", ProductID: "Punto", Colour: "violet", Status: "Pari"},
		Car{Description: "Tata", ProductID: "Nano", Colour: "indigo", Status: "Valeria"},
		Car{Description: "Holden", ProductID: "Barina", Colour: "brown", Status: "Shotaro"},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, description string, productID string, colour string, status string) error {
	car := Car{
		Description:   description,
		ProductID:  productID,
		Colour: colour,
		Status:  status,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarStatus updates the status field of car with given id in world state
func (s *SmartContract) ChangeCarStatus(ctx contractapi.TransactionContextInterface, carNumber string, newStatus string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	car.Status = newStatus

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}

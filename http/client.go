package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// UserStatus represents user data, balance and debts
type UserStatus struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	BalanceAmount string              `json:"balance_amount"`
	Debts         []map[string]string `json:"debts"`
}

// GetUserStatusSync hit necessary endpoints and join user's data sync.
func GetUserStatusSync(serverURL, userID string) (UserStatus, error) {
	userResponse, _ := http.Get(fmt.Sprintf("%s/users/%s", serverURL, userID))
	balanceResponse, _ := http.Get(fmt.Sprintf("%s/balance/%s", serverURL, userID))
	debtsResponse, _ := http.Get(fmt.Sprintf("%s/user-debts/%s", serverURL, userID))
	var userInfo map[string]string
	unmarshalResponse(userResponse, &userInfo)
	var userBalance map[string]string
	unmarshalResponse(balanceResponse, &userBalance)
	var userDebts []map[string]string
	unmarshalResponse(debtsResponse, &userDebts)
	return UserStatus{
		ID:            userInfo["id"],
		Name:          userInfo["name"],
		BalanceAmount: userBalance["amount"],
		Debts:         userDebts,
	}, nil
}

// GetUserStatusAsyncWaitGroup hit necessary endpoints and join user's data async with waitgroups.
func GetUserStatusAsyncWaitGroup(serverURL, userID string) (UserStatus, error) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(3)
	var userResponse, balanceResponse, debtsResponse *http.Response
	go func() {
		userResponse, _ = http.Get(fmt.Sprintf("%s/users/%s", serverURL, userID))
		waitgroup.Done()
	}()
	go func() {
		balanceResponse, _ = http.Get(fmt.Sprintf("%s/balance/%s", serverURL, userID))
		waitgroup.Done()
	}()
	go func() {
		debtsResponse, _ = http.Get(fmt.Sprintf("%s/user-debts/%s", serverURL, userID))
		waitgroup.Done()
	}()
	waitgroup.Wait()
	var userInfo map[string]string
	unmarshalResponse(userResponse, &userInfo)
	var userBalance map[string]string
	unmarshalResponse(balanceResponse, &userBalance)
	var userDebts []map[string]string
	unmarshalResponse(debtsResponse, &userDebts)
	return UserStatus{
		ID:            userInfo["id"],
		Name:          userInfo["name"],
		BalanceAmount: userBalance["amount"],
		Debts:         userDebts,
	}, nil
}

// GetUserStatusAsyncChannels hit necessary endpoints and join user's data async with waitgroups.
func GetUserStatusAsyncChannels(serverURL, userID string) (UserStatus, error) {

	userResponse := make(chan *http.Response)
	balanceResponse := make(chan *http.Response)
	debtsResponse := make(chan *http.Response)
	defer close(userResponse)
	defer close(balanceResponse)
	defer close(debtsResponse)

	go func() {
		result, _ := http.Get(fmt.Sprintf("%s/users/%s", serverURL, userID))
		userResponse <- result
	}()
	go func() {
		result, _ := http.Get(fmt.Sprintf("%s/balance/%s", serverURL, userID))
		balanceResponse <- result
	}()
	go func() {
		result, _ := http.Get(fmt.Sprintf("%s/user-debts/%s", serverURL, userID))
		debtsResponse <- result
	}()

	var userInfo map[string]string
	unmarshalResponse(<-userResponse, &userInfo)
	var userBalance map[string]string
	unmarshalResponse(<-balanceResponse, &userBalance)
	var userDebts []map[string]string
	unmarshalResponse(<-debtsResponse, &userDebts)
	return UserStatus{
		ID:            userInfo["id"],
		Name:          userInfo["name"],
		BalanceAmount: userBalance["amount"],
		Debts:         userDebts,
	}, nil
}

func unmarshalResponse(r *http.Response, b interface{}) {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(bytes, &b)
}

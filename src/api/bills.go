// Get bills from the API
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Bill struct {
  BillId int `json:"bill_id"`
  Amount float64 `json:"amount"`
}

type GetBillRes struct {
  Result []Bill `json:"result"`
}

type PayBillsRes struct {
	Result []int `json:"result"`
}

type BillPayReq struct {
	Tx string `json:"tx"`
	BillIds []int `json:"bill_ids"`
}

/**
 * Get an array of Bill objects from the API. This is an authenticated request,
 * so a valid JSON web token must be included
 *
 * @param  api      Base URI for the hub API
 * @param  token    JSON web token for the agent
 * @return          (array of bills, error)
 */
func GetBills(api string, token string) (*[]Bill, error) {
  var result = new(GetBillRes)
  var bills = new([]Bill)
  client := &http.Client{}
  req, _ := http.NewRequest("GET", api+"/Bills", nil)
  req.Header.Set("x-access-token", token)
  res, _ := client.Do(req)
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return bills, fmt.Errorf("Could not read response body (%s)", err)
  } else {
    err2 := json.Unmarshal(body, &result)
    if err2 != nil {
      return bills, fmt.Errorf("Could not unmarshal body (%s)", err)
    }
  }
  return &result.Result, nil

}

/**
 * Get an array of Bill objects from the API. This is an authenticated request,
 * so a valid JSON web token must be included
 *
 * @param  ids           Array of bill_ids (SQL ids from the API)
 * @param  tx            Raw transaction string signed by agent
 * @param  api           Base URI for the hub API
 * @param  auth_token    JSON web token for the agent
 * @return               (array of bill ids, error)
 */
func PayBills(ids []int, tx string, api string, auth_token string) ([]int, error) {
	payload := BillPayReq{tx, ids}
	b, _ := json.Marshal(payload)
	var result = new(PayBillsRes)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", api+"/PayBills", bytes.NewBuffer(b))
  req.Header.Set("x-access-token", auth_token)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, fmt.Errorf("Could not read response body (%s)", err)
  } else {
    err2 := json.Unmarshal(body, &result)
    if err2 != nil {
      return nil, fmt.Errorf("Could not unmarshal body (%s)", err)
    }
  }
  return result.Result, nil

}

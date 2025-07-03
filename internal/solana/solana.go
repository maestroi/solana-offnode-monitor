package solana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

func (c *Client) rpcRequest(method string, params interface{}, result interface{}) error {
	body, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	resp, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  interface{}     `json:"error"`
	}
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("rpc error: %v", rpcResp.Error)
	}
	return json.Unmarshal(rpcResp.Result, result)
}

// Example: GetVoteAccounts returns the full vote accounts response
func (c *Client) GetVoteAccounts() (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := c.rpcRequest("getVoteAccounts", []interface{}{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Example: GetBalance returns the balance for a pubkey
func (c *Client) GetBalance(pubkey string) (uint64, error) {
	var result struct {
		Value uint64 `json:"value"`
	}
	if err := c.rpcRequest("getBalance", []interface{}{pubkey}, &result); err != nil {
		return 0, err
	}
	return result.Value, nil
}

// GetEpochInfo returns the current epoch info
func (c *Client) GetEpochInfo() (uint64, error) {
	var result struct {
		Epoch uint64 `json:"epoch"`
	}
	if err := c.rpcRequest("getEpochInfo", []interface{}{}, &result); err != nil {
		return 0, err
	}
	return result.Epoch, nil
}

type ValidatorInfo struct {
	Identity string
	Name     string
	Website  string
	Details  string
}

// GetValidatorInfo returns a map of identity pubkey to ValidatorInfo
func (c *Client) GetValidatorInfo() (map[string]ValidatorInfo, error) {
	var rpcResult struct {
		Result []struct {
			Info struct {
				Identity string `json:"identity"`
				Name     string `json:"name"`
				Website  string `json:"website"`
				Details  string `json:"details"`
			} `json:"info"`
		} `json:"result"`
	}
	if err := c.rpcRequest("getValidatorInfo", []interface{}{}, &rpcResult); err != nil {
		return nil, err
	}
	infoMap := make(map[string]ValidatorInfo)
	for _, v := range rpcResult.Result {
		infoMap[v.Info.Identity] = ValidatorInfo{
			Identity: v.Info.Identity,
			Name:     v.Info.Name,
			Website:  v.Info.Website,
			Details:  v.Info.Details,
		}
	}
	return infoMap, nil
}

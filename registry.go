package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// getchains() gets the list of chains from the chain registry
func GetChains() Chains {
	resp, err := http.Get("https://cosmos-chain.directory/chains")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var chains Chains

	err = json.Unmarshal([]byte(body), &chains)
	if err != nil {
		fmt.Println(err)
	}
	return chains
}

// getchain() gets one chain's records from the chain registry
func GetChain(chainid string) Chain {
	resp, err := http.Get("https://cosmos-chain.directory/chains/" + chainid)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var chain Chain

	err = json.Unmarshal([]byte(body), &chain)
	if err != nil {
		fmt.Println(err)
	}
	return chain
}

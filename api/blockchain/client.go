package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainClient struct {
	Client  *ethclient.Client
	ChainID *big.Int
}

func NewClient(networkURL string) (*BlockchainClient, error) {
	client, err := ethclient.Dial(networkURL)
	if err != nil {
		return nil, fmt.Errorf("error dialing Ethereum client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying chain ID: %v", err)
	}

	return &BlockchainClient{
		Client:  client,
		ChainID: chainID,
	}, nil
}


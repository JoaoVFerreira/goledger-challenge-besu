package blockchain

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ContractService struct {
	Client          *BlockchainClient
	ContractAddress common.Address
	ABI             abi.ABI
	PrivateKey      string
}

func NewContractService(client *BlockchainClient, contractAddress, PrivateKey, abiJSON string) (*ContractService, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("error parsing ABI: %v", err)
	}

	return &ContractService{
		Client:          client,
		ContractAddress: common.HexToAddress(contractAddress),
		ABI:             parsedABI,
		PrivateKey:      PrivateKey,
	}, nil
}

func (s *ContractService) ExecContract(methodName string, params ...interface{}) error {
	boundContract := bind.NewBoundContract(
		s.ContractAddress,
		s.ABI,
		s.Client.Client,
		s.Client.Client,
		s.Client.Client,
	)

	privKey, err := crypto.HexToECDSA(s.PrivateKey)
	if err != nil {
		return fmt.Errorf("error loading private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, s.Client.ChainID)
	if err != nil {
		return fmt.Errorf("error creating transactor: %v", err)
	}

	tx, err := boundContract.Transact(auth, methodName, params...)
	if err != nil {
		return fmt.Errorf("error executing contract method: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), s.Client.Client, tx)
	if err != nil {
		return fmt.Errorf("error waiting for transaction to be mined: %v", err)
	}

	log.Printf("Transaction mined: %v", receipt)
	return nil
}

func (s *ContractService) CallContract(methodName string, result *[]interface{}, params ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	boundContract := bind.NewBoundContract(
		s.ContractAddress,
		s.ABI,
		s.Client.Client,
		s.Client.Client,
		s.Client.Client,
	)

	if err := boundContract.Call(&caller, result, methodName, params...); err != nil {
		return fmt.Errorf("error when calling method: %v", err)
	}

	return nil
}

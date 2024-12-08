package main

import (
	"fmt"
	"ledger-api/blockchain"
	"ledger-api/database"
	"ledger-api/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const ABI_FILE_PATH = "./configs/SimpleStorageABI.json"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading the .env file")
	}

	CONTRACT_ADDRESS := os.Getenv("CONTRACT_ADDRESS")
	BLOCKCHAIN_NODE  := os.Getenv("BLOCKCHAIN_NODE")
	PRIVATE_KEY 		 := os.Getenv("PRIVATE_KEY")
	DATABASE_URL 		 := os.Getenv("DATABASE_URL")
	PORT 						 := os.Getenv("PORT")
	
	db, err := database.NewDatabaseConnection(DATABASE_URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	ABI, err := os.ReadFile(ABI_FILE_PATH)
	if err != nil {
		log.Fatalf("error loading the ABI: %v", err)
	}

	blockchainClient, err := blockchain.NewClient(BLOCKCHAIN_NODE)
	if err != nil {
		log.Fatalf("error loading the blockchain client: %v", err)
	}

	contractService, err := blockchain.NewContractService(blockchainClient, CONTRACT_ADDRESS, PRIVATE_KEY, string(ABI))
	if err != nil {
		log.Fatalf("error loading the blockchain service: %v", err)
	}

	handler := handlers.Handler{
		DB:             db,
		ContractClient: contractService,
	}

	router := gin.Default()
	router.SetTrustedProxies(nil) 
	router.GET("/simple-storage/get/value", handler.GetContractValue)
	router.GET("/simple-storage/check/value", handler.CheckContractValue)
	router.GET("/simple-storage/sync/value", handler.SyncContractValue)
	router.POST("/simple-storage/set/value", handler.SetContractValue)

	log.Printf("Server is running on port %s", PORT)
	router.Run(fmt.Sprintf(":%s", PORT))	
}

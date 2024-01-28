package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	localwallet "github.com/c0rwin/fabric-gateway-wallet/pkg/local"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	baseDir := fmt.Sprintf("%s/workspace/hyperledger/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com", homedir)

	wallet, err := localwallet.NewLocalWallet(
		localwallet.WithBaseDir(baseDir),
		localwallet.WithMSPId("Org1MSP"),
		localwallet.WithTLSFolder("peers/peer0.org1.example.com/tls"),
		localwallet.WithTLSRootCAFile("ca.crt"),
		localwallet.WithTLSHostName("peer0.org1.example.com"),
		localwallet.WithKeyStore("users/User1@org1.example.com/msp/keystore"),
		localwallet.WithSignIndentity("users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"),
	)

	if err != nil {
		panic(err)
	}

	transportCreds, err := wallet.TransportCredentials()
	if err != nil {
		panic(err)
	}
	grpcConnection, err := grpc.Dial("peer0.org1.example.com:7051", grpc.WithTransportCredentials(transportCreds))
	if err != nil {
		panic(err)
	}
	defer grpcConnection.Close()

	id, err := wallet.Identity()
	if err != nil {
		panic(err)
	}

	sign, err := wallet.Sign()
	if err != nil {
		panic(err)
	}

	gw, err := client.Connect(id,
		client.WithSign(sign),
		client.WithClientConnection(grpcConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		panic(err)
	}

	defer gw.Close()

	mychannel := gw.GetNetwork("mychannel")
	chaincode := mychannel.GetContract("mychaincode")

	_, err = chaincode.SubmitTransaction("CreateAsset", "asset2", "red", "5", "Alice", "1300")
	if err != nil {
		panic(err)
	}

	fmt.Println("Asset created successfully")

	select {
	case <-ctx.Done():
		// handle context cancellation
	}
}

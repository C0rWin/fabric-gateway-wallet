# Hyperledger Fabric Gateway SDK Wallet Implementation in Go

## Introduction

This repository provides a Go-based implementation for the wallet component of the Hyperledger Fabric Gateway SDK. It focuses on abstracting the complexities involved in handling cryptographic materials, such as loading TLS certificates, signing certificates, and private keys. These functionalities are crucial for creating user identities, signing identities, and establishing TLS connection credentials, thereby facilitating a seamless connection to the Fabric Gateway.

## Features

- **Cryptographic Material Handling**: Efficient management of TLS certificates, signing certificates, and private keys.
- **Identity Management**: Streamlined creation and management of user identities and signing identities.
- **TLS Credential Support**: Simplified generation of credentials for establishing secure TLS connections.

## Installation

Ensure you have Go installed on your system. Then, follow these steps to install the wallet implementation:

```bash
go get github.com/c0rwin/fabric-gateway-sdwallet
```

Or, clone the repository and build locally:

```bash
git clone [Your Repository URL]
cd [Your Repository Directory]
go build
```

## Usage

```go

func main() {
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
}
```

## Contributing

Contributions are welcome! If you would like to contribute:

* Fork the repository.
* Create a new feature branch.
* Commit your changes.
* Push to the branch.
* Create a new Pull Request.


## License
This project is licensed under the Apache License 2.0. For more information, see the LICENSE file in this repository.
package wallet

import (
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc/credentials"
)

// Wallet is an interface for providing credentials and signing for the client
type Wallet interface {
	// TransportCredentials returns the credentials to use for the gRPC connection
	TransportCredentials() (credentials.TransportCredentials, error)
	// Sign returns the signing implementation to use for signing transactions
	Sign() (identity.Sign, error)
	// Identity returns the identity to use for signing transactions
	Identity() (identity.Identity, error)
}

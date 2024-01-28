package localwallet

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc/credentials"
)

// LocalWalletOption is a functional option for configuring the LocalWallet
type LocalWalletOption func(*LocalWallet)

// WithBaseDir sets the base directory for the LocalWallet
func WithBaseDir(baseDir string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.baseDir = baseDir
	}
}

// WithTLSFolder sets the TLS folder for the LocalWallet
func WithTLSFolder(tlsFolder string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.tlsFolder = tlsFolder
	}
}

// WithTLSRootCAFile sets the TLS root CA file for the LocalWallet
func WithTLSRootCAFile(tlsRootCAFile string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.tlsRootCAFile = tlsRootCAFile
	}
}

// WithTLSHostName sets the TLS host name for the LocalWallet
func WithTLSHostName(tlsHostName string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.tlsHostName = tlsHostName
	}
}

// WithKeyStore sets the key store for the LocalWallet
func WithKeyStore(keyStore string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.keyStore = keyStore
	}
}

// WithSignIndentityDir sets the sign identity directory for the LocalWallet
func WithSignIndentity(signIdentity string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.signIdentity = signIdentity
	}
}

// WithMSPId sets the MSP ID for the LocalWallet
func WithMSPId(mspID string) LocalWalletOption {
	return func(wallet *LocalWallet) {
		wallet.mspID = mspID
	}
}

// LocalWallet is a wallet implementation that stores identities in the local filesystem
type LocalWallet struct {
	baseDir       string
	tlsFolder     string
	tlsRootCAFile string
	tlsHostName   string
	keyStore      string
	signIdentity  string
	mspID         string
}

func NewLocalWallet(options ...LocalWalletOption) (*LocalWallet, error) {
	wallet := &LocalWallet{}

	for _, option := range options {
		option(wallet)
	}

	if wallet.baseDir == "" {
		return nil, fmt.Errorf("baseDir is required")
	}

	if wallet.tlsFolder == "" {
		return nil, fmt.Errorf("tlsFolder is required")
	}

	if wallet.tlsRootCAFile == "" {
		return nil, fmt.Errorf("tlsRootCAFile is required")
	}

	if wallet.tlsHostName == "" {
		return nil, fmt.Errorf("tlsHostName is required")
	}

	if wallet.keyStore == "" {
		return nil, fmt.Errorf("keyStore is required")
	}

	if wallet.signIdentity == "" {
		return nil, fmt.Errorf("signIdentity is required")
	}

	if wallet.mspID == "" {
		return nil, fmt.Errorf("mspID is required")
	}

	return wallet, nil
}

// TransportCredentials returns the credentials to use for the gRPC connection
func (w *LocalWallet) TransportCredentials() (credentials.TransportCredentials, error) {
	pemTLS, err := os.ReadFile(filepath.Clean(filepath.Join(w.baseDir, w.tlsFolder, w.tlsRootCAFile)))
	if err != nil {
		return nil, err
	}

	tlsCert, err := identity.CertificateFromPEM(pemTLS)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(tlsCert)

	return credentials.NewClientTLSFromCert(certPool, w.tlsHostName), nil
}

// Sign returns the signing implementation to use for signing transactions
func (w *LocalWallet) Identity() (identity.Identity, error) {
	pem, err := os.ReadFile(filepath.Clean(filepath.Join(w.baseDir, w.signIdentity)))
	if err != nil {
		return nil, err
	}

	cert, err := identity.CertificateFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return identity.NewX509Identity(w.mspID, cert)
}

// Sign returns the signing implementation to use for signing transactions
func (w *LocalWallet) Sign() (identity.Sign, error) {
	files, err := os.ReadDir(filepath.Clean(filepath.Join(w.baseDir, w.keyStore)))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in %s", filepath.Join(w.baseDir, w.keyStore))
	}

	pem, err := os.ReadFile(filepath.Clean(filepath.Join(w.baseDir, w.keyStore, files[0].Name())))
	if err != nil {
		return nil, err
	}

	pk, err := identity.PrivateKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}

	return identity.NewPrivateKeySign(pk)
}

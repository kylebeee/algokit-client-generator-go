package testutil

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	algokit "github.com/kylebeee/algokit-utils-go"
)

// LocalNetAlgodURL is the default algod URL for localnet.
const LocalNetAlgodURL = "http://localhost:4001"

// LocalNetAlgodToken is the default algod token for localnet.
const LocalNetAlgodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

// LocalNetKMDURL is the default KMD URL for localnet.
const LocalNetKMDURL = "http://localhost:4002"

// LocalNetKMDToken is the default KMD token for localnet.
const LocalNetKMDToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

// TestFixture provides test utilities for localnet integration tests.
type TestFixture struct {
	T        *testing.T
	Algod    *algod.Client
	Algorand *algokit.AlgorandClient
	Ctx      context.Context
}

// NewTestFixture creates a new test fixture connected to localnet.
func NewTestFixture(t *testing.T) *TestFixture {
	t.Helper()

	algorand, err := algokit.LocalNet()
	if err != nil {
		t.Fatalf("failed to create localnet client: %v", err)
	}

	return &TestFixture{
		T:        t,
		Algod:    algorand.Algod(),
		Algorand: algorand,
		Ctx:      context.Background(),
	}
}

// TestAccount wraps a crypto.Account with a signer for convenience.
type TestAccount struct {
	Account crypto.Account
	Address types.Address
	Signer  transaction.TransactionSigner
}

// GenerateAccount creates a new random account and funds it from the dispenser.
func (f *TestFixture) GenerateAccount(fundAmount uint64) *TestAccount {
	f.T.Helper()

	account := crypto.GenerateAccount()
	addr := types.Address(account.Address)
	signer := transaction.BasicAccountTransactionSigner{Account: account}

	// Register signer
	f.Algorand.Account().SetSignerFromAccount(account)

	// Fund from dispenser
	if fundAmount > 0 {
		f.FundAccount(addr, fundAmount)
	}

	return &TestAccount{
		Account: account,
		Address: addr,
		Signer:  signer,
	}
}

// GetDispenser returns the default localnet dispenser account.
func (f *TestFixture) GetDispenser() *TestAccount {
	f.T.Helper()

	// The localnet dispenser is the first account from the default wallet.
	// For AlgoKit localnet, this is typically available at the KMD endpoint.
	// We use a well-known mnemonic approach or KMD API.
	// For simplicity, use the go-algorand-sdk KMD API.
	kmdClient, err := kmd.MakeClient(LocalNetKMDURL, LocalNetKMDToken)
	if err != nil {
		f.T.Fatalf("failed to create KMD client: %v", err)
	}

	// List wallets
	walletsResp, err := kmdClient.ListWallets()
	if err != nil {
		f.T.Fatalf("failed to list wallets: %v", err)
	}

	var walletID string
	for _, w := range walletsResp.Wallets {
		if w.Name == "unencrypted-default-wallet" {
			walletID = w.ID
			break
		}
	}
	if walletID == "" {
		f.T.Fatal("default wallet not found")
	}

	// Init wallet handle
	handleResp, err := kmdClient.InitWalletHandle(walletID, "")
	if err != nil {
		f.T.Fatalf("failed to init wallet handle: %v", err)
	}

	// List keys
	keysResp, err := kmdClient.ListKeys(handleResp.WalletHandleToken)
	if err != nil {
		f.T.Fatalf("failed to list keys: %v", err)
	}

	if len(keysResp.Addresses) == 0 {
		f.T.Fatal("no accounts in default wallet")
	}

	// Export key for the first account
	keyResp, err := kmdClient.ExportKey(handleResp.WalletHandleToken, "", keysResp.Addresses[0])
	if err != nil {
		f.T.Fatalf("failed to export key: %v", err)
	}

	account, err := crypto.AccountFromPrivateKey(keyResp.PrivateKey)
	if err != nil {
		f.T.Fatalf("failed to create account from private key: %v", err)
	}

	addr := types.Address(account.Address)
	signer := transaction.BasicAccountTransactionSigner{Account: account}
	f.Algorand.Account().SetSignerFromAccount(account)

	return &TestAccount{
		Account: account,
		Address: addr,
		Signer:  signer,
	}
}

// FundAccount sends ALGO from the dispenser to the given address.
func (f *TestFixture) FundAccount(addr types.Address, microAlgos uint64) {
	f.T.Helper()

	dispenser := f.GetDispenser()

	_, err := algokit.SendPayment(f.Ctx, f.Algod, algokit.PaymentParams{
		Sender:   dispenser.Address,
		Signer:   dispenser.Signer,
		Receiver: addr,
		Amount:   algokit.MicroAlgos(microAlgos),
	})
	if err != nil {
		f.T.Fatalf("failed to fund account %s: %v", addr.String(), err)
	}
}

// FundApp sends ALGO to an application's escrow address.
func (f *TestFixture) FundApp(appID uint64, microAlgos uint64) {
	f.T.Helper()

	appAddr := algokit.GetAppAddress(appID)
	dispenser := f.GetDispenser()

	_, err := algokit.SendPayment(f.Ctx, f.Algod, algokit.PaymentParams{
		Sender:   dispenser.Address,
		Signer:   dispenser.Signer,
		Receiver: appAddr,
		Amount:   algokit.MicroAlgos(microAlgos),
	})
	if err != nil {
		f.T.Fatalf("failed to fund app %d: %v", appID, err)
	}
}

// CreateASA creates a test asset and returns its ID.
func (f *TestFixture) CreateASA(creator *TestAccount, total uint64, decimals uint32, unitName, assetName string) uint64 {
	f.T.Helper()

	result, err := algokit.NewAssetManager(f.Algod).Create(f.Ctx, algokit.AssetCreateParams{
		Sender:    creator.Address,
		Signer:    creator.Signer,
		Total:     total,
		Decimals:  decimals,
		UnitName:  unitName,
		AssetName: assetName,
	})
	if err != nil {
		f.T.Fatalf("failed to create ASA: %v", err)
	}

	return result.Confirmation.AssetIndex
}

// OptInASA opts an account into an asset.
func (f *TestFixture) OptInASA(account *TestAccount, assetID uint64) {
	f.T.Helper()

	_, err := algokit.NewAssetManager(f.Algod).OptIn(f.Ctx, algokit.AssetOptInParams{
		Sender:  account.Address,
		Signer:  account.Signer,
		AssetID: assetID,
	})
	if err != nil {
		f.T.Fatalf("failed to opt in to ASA %d: %v", assetID, err)
	}
}

// TransferASA transfers an asset between accounts.
func (f *TestFixture) TransferASA(sender *TestAccount, receiver types.Address, assetID uint64, amount uint64) {
	f.T.Helper()

	_, err := algokit.NewAssetManager(f.Algod).Transfer(f.Ctx, algokit.AssetTransferParams{
		Sender:   sender.Address,
		Signer:   sender.Signer,
		AssetID:  assetID,
		Receiver: receiver,
		Amount:   amount,
	})
	if err != nil {
		f.T.Fatalf("failed to transfer ASA %d: %v", assetID, err)
	}
}

// GetAccountBalance returns the account's microAlgo balance.
func (f *TestFixture) GetAccountBalance(addr types.Address) uint64 {
	f.T.Helper()

	info, err := f.Algod.AccountInformation(addr.String()).Do(f.Ctx)
	if err != nil {
		f.T.Fatalf("failed to get account info: %v", err)
	}

	return info.Amount
}

// MakePaymentTxn creates a payment transaction (for use as method args).
func (f *TestFixture) MakePaymentTxn(sender *TestAccount, receiver types.Address, amount uint64) transaction.TransactionWithSigner {
	f.T.Helper()

	txn, err := algokit.MakePaymentTxn(f.Ctx, f.Algod, algokit.PaymentParams{
		Sender:   sender.Address,
		Receiver: receiver,
		Amount:   algokit.MicroAlgos(amount),
	})
	if err != nil {
		f.T.Fatalf("failed to create payment txn: %v", err)
	}

	return transaction.TransactionWithSigner{
		Txn:    txn,
		Signer: sender.Signer,
	}
}

// MakeAssetTransferTxn creates an asset transfer transaction (for use as method args).
func (f *TestFixture) MakeAssetTransferTxn(sender *TestAccount, receiver types.Address, assetID uint64, amount uint64) transaction.TransactionWithSigner {
	f.T.Helper()

	sp, err := f.Algod.SuggestedParams().Do(f.Ctx)
	if err != nil {
		f.T.Fatalf("failed to get suggested params: %v", err)
	}

	txn, err := transaction.MakeAssetTransferTxn(
		sender.Address.String(),
		receiver.String(),
		amount,
		nil, // note
		sp,
		"", // close to
		assetID,
	)
	if err != nil {
		f.T.Fatalf("failed to create asset transfer txn: %v", err)
	}

	return transaction.TransactionWithSigner{
		Txn:    txn,
		Signer: sender.Signer,
	}
}

// AdvanceRounds sends N dummy transactions to advance the round number.
func (f *TestFixture) AdvanceRounds(n int) {
	f.T.Helper()

	dispenser := f.GetDispenser()
	for i := 0; i < n; i++ {
		note := make([]byte, 8)
		_, _ = rand.Read(note)
		_, err := algokit.SendPayment(f.Ctx, f.Algod, algokit.PaymentParams{
			Sender:   dispenser.Address,
			Signer:   dispenser.Signer,
			Receiver: dispenser.Address,
			Amount:   algokit.MicroAlgos(0),
			Note:     note,
		})
		if err != nil {
			f.T.Fatalf("failed to advance round: %v", err)
		}
	}
}

// RequireError checks that an error was returned and optionally contains a substring.
func RequireError(t *testing.T, err error, msgSubstring ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if len(msgSubstring) > 0 {
		errStr := err.Error()
		for _, sub := range msgSubstring {
			if !contains(errStr, sub) {
				t.Errorf("error %q does not contain %q", errStr, sub)
			}
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// Logf logs a test message.
func (f *TestFixture) Logf(format string, args ...interface{}) {
	f.T.Helper()
	f.T.Logf(format, args...)
}

// AlgoAmount returns micro algos from a whole algo amount.
func AlgoAmount(algos float64) uint64 {
	return uint64(algos * 1_000_000)
}

// FormatAlgos formats micro algos as a human-readable ALGO string.
func FormatAlgos(microAlgos uint64) string {
	algos := float64(microAlgos) / 1_000_000.0
	return fmt.Sprintf("%.6f ALGO", algos)
}

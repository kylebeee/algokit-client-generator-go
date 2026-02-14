package tests

import (
	"context"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/prizebox"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// prizeBoxTestEnv holds the shared state for PrizeBox tests.
type prizeBoxTestEnv struct {
	fixture  *testutil.TestFixture
	ctx      context.Context
	owner    *testutil.TestAccount
	user1    *testutil.TestAccount
	client   *prizebox.Client
	factory  *prizebox.Factory
}

func setupPrizeBoxTest(t *testing.T) *prizeBoxTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	user1 := fixture.GenerateAccount(testutil.AlgoAmount(10))

	factory, err := prizebox.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create prizebox factory: %v", err)
	}

	// Deploy with create(owner) method call
	client, _, err := factory.Create(ctx, prizebox.FactoryCreateParams{
		Args: prizebox.CreateArgs{
			Owner: owner.Address,
		},
	})
	if err != nil {
		t.Fatalf("failed to deploy prizebox contract: %v", err)
	}

	t.Logf("PrizeBox contract deployed with appID: %d", client.AppID())

	// Fund the contract for MBR
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(1))

	return &prizeBoxTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		user1:   user1,
		client:  client,
		factory: factory,
	}
}

func TestPrizeBoxDeployment(t *testing.T) {
	env := setupPrizeBoxTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("PrizeBox deployed with appID: %d", env.client.AppID())
}

func TestPrizeBoxOptIn(t *testing.T) {
	env := setupPrizeBoxTest(t)

	// Create a test ASA
	assetID := env.fixture.CreateASA(env.owner, 1_000_000, 0, "Prize", "PTK")

	// MBR payment for ASA opt-in (100,000 microAlgos = assetOptInMinBalance)
	optInPayment := env.fixture.MakePaymentTxn(env.owner, env.client.AppAddress(), 100_000)

	err := env.client.SendOptin(env.ctx, algokit.CallParams[prizebox.OptinArgs]{
		Args: prizebox.OptinArgs{
			Payment: optInPayment,
			Asset:   assetID,
		},
		Sender:          env.owner.Address,
		Signer:          env.owner.Signer,
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000, // inner ASA opt-in txn
	})
	if err != nil {
		t.Fatalf("failed to opt into ASA: %v", err)
	}

	t.Logf("PrizeBox opted into ASA %d", assetID)
}

func TestPrizeBoxOptInNotOwner(t *testing.T) {
	env := setupPrizeBoxTest(t)

	assetID := env.fixture.CreateASA(env.owner, 1_000_000, 0, "Prize", "PTK")

	// Try opt-in as non-owner
	optInPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), 100_000)

	err := env.client.SendOptin(env.ctx, algokit.CallParams[prizebox.OptinArgs]{
		Args: prizebox.OptinArgs{
			Payment: optInPayment,
			Asset:   assetID,
		},
		Sender:          env.user1.Address,
		Signer:          env.user1.Signer,
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000,
	})
	if err == nil {
		t.Fatal("expected error when non-owner tries to opt in")
	}

	t.Logf("Non-owner opt-in correctly rejected: %v", err)
}

func TestPrizeBoxTransfer(t *testing.T) {
	env := setupPrizeBoxTest(t)

	// Transfer ownership to user1
	err := env.client.SendTransfer(env.ctx, algokit.CallParams[prizebox.TransferArgs]{
		Args: prizebox.TransferArgs{
			NewOwner: env.user1.Address,
		},
		Sender: env.owner.Address,
		Signer: env.owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to transfer ownership: %v", err)
	}

	t.Logf("Transferred ownership to user1")
}

func TestPrizeBoxTransferNotOwner(t *testing.T) {
	env := setupPrizeBoxTest(t)

	// Try transferring as non-owner
	err := env.client.SendTransfer(env.ctx, algokit.CallParams[prizebox.TransferArgs]{
		Args: prizebox.TransferArgs{
			NewOwner: env.user1.Address,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
	})
	if err == nil {
		t.Fatal("expected error when non-owner tries to transfer")
	}

	t.Logf("Non-owner transfer correctly rejected: %v", err)
}

func TestPrizeBoxOptInAndWithdrawASA(t *testing.T) {
	env := setupPrizeBoxTest(t)

	// Create ASA and opt in
	assetID := env.fixture.CreateASA(env.owner, 1_000_000, 0, "Prize", "PTK")

	optInPayment := env.fixture.MakePaymentTxn(env.owner, env.client.AppAddress(), 100_000)

	err := env.client.SendOptin(env.ctx, algokit.CallParams[prizebox.OptinArgs]{
		Args: prizebox.OptinArgs{
			Payment: optInPayment,
			Asset:   assetID,
		},
		Sender:          env.owner.Address,
		Signer:          env.owner.Signer,
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000,
	})
	if err != nil {
		t.Fatalf("failed to opt into ASA: %v", err)
	}

	// Transfer some ASA to the contract
	env.fixture.TransferASA(env.owner, env.client.AppAddress(), assetID, 500)

	// Withdraw partial amount (not closing out)
	err = env.client.SendWithdraw(env.ctx, algokit.CallParams[prizebox.WithdrawArgs]{
		Args: prizebox.WithdrawArgs{
			Assets: [][]interface{}{
				{assetID, uint64(200)}, // (asset, amount) tuple
			},
		},
		Sender:          env.owner.Address,
		Signer:          env.owner.Signer,
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000, // inner axfer txn
	})
	if err != nil {
		t.Fatalf("failed to withdraw ASA: %v", err)
	}

	t.Logf("Withdrew 200 of ASA %d from PrizeBox", assetID)
}

func TestPrizeBoxOptInAndCloseOutASA(t *testing.T) {
	env := setupPrizeBoxTest(t)

	// Create ASA, opt in, and transfer to contract
	assetID := env.fixture.CreateASA(env.owner, 1_000_000, 0, "Prize", "PTK")

	optInPayment := env.fixture.MakePaymentTxn(env.owner, env.client.AppAddress(), 100_000)

	err := env.client.SendOptin(env.ctx, algokit.CallParams[prizebox.OptinArgs]{
		Args: prizebox.OptinArgs{
			Payment: optInPayment,
			Asset:   assetID,
		},
		Sender:          env.owner.Address,
		Signer:          env.owner.Signer,
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000,
	})
	if err != nil {
		t.Fatalf("failed to opt into ASA: %v", err)
	}

	env.fixture.TransferASA(env.owner, env.client.AppAddress(), assetID, 500)

	// Withdraw full amount (close out - amount == balance)
	err = env.client.SendWithdraw(env.ctx, algokit.CallParams[prizebox.WithdrawArgs]{
		Args: prizebox.WithdrawArgs{
			Assets: [][]interface{}{
				{assetID, uint64(500)}, // close out: amount == balance
			},
		},
		Sender:          env.owner.Address,
		Signer:          env.owner.Signer,
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000, // inner axfer with close-to
	})
	if err != nil {
		t.Fatalf("failed to close out ASA: %v", err)
	}

	t.Logf("Closed out ASA %d from PrizeBox", assetID)
}

func TestPrizeBoxFactory(t *testing.T) {
	env := setupPrizeBoxTest(t)

	client2, _, err := env.factory.Create(env.ctx, prizebox.FactoryCreateParams{
		Args: prizebox.CreateArgs{
			Owner: env.owner.Address,
		},
	})
	if err != nil {
		t.Fatalf("failed to deploy second prizebox: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second PrizeBox deployed with appID: %d", client2.AppID())
}

// ensure imports are used
var _ = types.Address{}

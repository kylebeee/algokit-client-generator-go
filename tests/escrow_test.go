package tests

import (
	"context"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/escrow"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// escrowTestEnv holds the shared state for Escrow tests.
type escrowTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	creator *testutil.TestAccount
	user1   *testutil.TestAccount
	client  *escrow.Client
	factory *escrow.Factory
}

func setupEscrowTest(t *testing.T) *escrowTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	creator := fixture.GenerateAccount(testutil.AlgoAmount(100))
	user1 := fixture.GenerateAccount(testutil.AlgoAmount(10))

	factory, err := escrow.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: creator.Address,
		DefaultSigner: creator.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create escrow factory: %v", err)
	}

	// Deploy with create(creator) - pass creator address as bytes
	client, _, err := factory.Create(ctx, escrow.FactoryCreateParams{
		Args: escrow.CreateArgs{
			Creator: creator.Address[:],
		},
	})
	if err != nil {
		t.Fatalf("failed to deploy escrow contract: %v", err)
	}

	t.Logf("Escrow contract deployed with appID: %d", client.AppID())

	// Fund the contract for MBR
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(1))

	return &escrowTestEnv{
		fixture: fixture,
		ctx:     ctx,
		creator: creator,
		user1:   user1,
		client:  client,
		factory: factory,
	}
}

func TestEscrowDeployment(t *testing.T) {
	env := setupEscrowTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("Escrow deployed with appID: %d", env.client.AppID())
}

func TestEscrowRekey(t *testing.T) {
	env := setupEscrowTest(t)

	// Rekey the escrow to user1
	err := env.client.SendRekey(env.ctx, algokit.CallParams[escrow.RekeyArgs]{
		Args: escrow.RekeyArgs{
			RekeyTo: env.user1.Address,
		},
		Sender:   env.creator.Address,
		Signer:   env.creator.Signer,
		ExtraFee: 1000, // inner rekey txn
	})
	if err != nil {
		t.Fatalf("failed to rekey escrow: %v", err)
	}

	t.Logf("Escrow rekeyed to user1")
}

func TestEscrowRekeyNotCreator(t *testing.T) {
	env := setupEscrowTest(t)

	// Try rekeying as non-creator
	err := env.client.SendRekey(env.ctx, algokit.CallParams[escrow.RekeyArgs]{
		Args: escrow.RekeyArgs{
			RekeyTo: env.user1.Address,
		},
		Sender:   env.user1.Address,
		Signer:   env.user1.Signer,
		ExtraFee: 1000,
	})
	if err == nil {
		t.Fatal("expected error when non-creator tries to rekey")
	}

	t.Logf("Non-creator rekey correctly rejected: %v", err)
}

func TestEscrowFactory(t *testing.T) {
	env := setupEscrowTest(t)

	client2, _, err := env.factory.Create(env.ctx, escrow.FactoryCreateParams{
		Args: escrow.CreateArgs{
			Creator: env.creator.Address[:],
		},
	})
	if err != nil {
		t.Fatalf("failed to deploy second escrow: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second Escrow deployed with appID: %d", client2.AppID())
}

// ensure imports are used
var _ = types.Address{}

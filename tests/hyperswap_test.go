package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/hyperswap"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// hyperswapTestEnv holds the shared state for HyperSwap tests.
type hyperswapTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *hyperswap.Client
	factory *hyperswap.Factory
}

func setupHyperSwapTest(t *testing.T) *hyperswapTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := hyperswap.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create hyperswap factory: %v", err)
	}

	// Deploy with create(version, akitaDAO) - pass 0 as dummy akitaDAO
	client, _, err := factory.Create(ctx, hyperswap.FactoryCreateParams{
		Args: hyperswap.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy hyperswap contract: %v", err)
	}

	t.Logf("HyperSwap contract deployed with appID: %d", client.AppID())

	// Fund the contract
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &hyperswapTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestHyperSwapDeployment(t *testing.T) {
	env := setupHyperSwapTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("HyperSwap deployed with appID: %d", env.client.AppID())
}

func TestHyperSwapMBR(t *testing.T) {
	env := setupHyperSwapTest(t)

	result, err := env.client.SendMBR(env.ctx)
	if err != nil {
		t.Fatalf("failed to call mbr: %v", err)
	}

	t.Logf("HyperSwap MBR - Offers: %d, Participants: %d, Hashes: %d, MM Root: %d, MM Data: %d",
		result.Return.Offers, result.Return.Participants, result.Return.Hashes,
		result.Return.Mm.Root, result.Return.Mm.Data)

	if result.Return.Offers == 0 {
		t.Error("expected non-zero offers MBR")
	}
	if result.Return.Participants == 0 {
		t.Error("expected non-zero participants MBR")
	}
}

func TestHyperSwapFactory(t *testing.T) {
	env := setupHyperSwapTest(t)

	client2, _, err := env.factory.Create(env.ctx, hyperswap.FactoryCreateParams{
		Args: hyperswap.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second hyperswap: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second HyperSwap deployed with appID: %d", client2.AppID())
}

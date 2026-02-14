package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/auctionfactory"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// auctionFactoryTestEnv holds the shared state for AuctionFactory tests.
type auctionFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *auctionfactory.Client
	factory *auctionfactory.Factory
}

func setupAuctionFactoryTest(t *testing.T) *auctionFactoryTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := auctionfactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create auction factory: %v", err)
	}

	// Deploy with create(version, childVersion, akitaDAO, akitaDAOEscrow)
	client, _, err := factory.Create(ctx, auctionfactory.FactoryCreateParams{
		Args: auctionfactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy auction factory contract: %v", err)
	}

	t.Logf("AuctionFactory contract deployed with appID: %d", client.AppID())

	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &auctionFactoryTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestAuctionFactoryDeployment(t *testing.T) {
	env := setupAuctionFactoryTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("AuctionFactory deployed with appID: %d", env.client.AppID())
}

func TestAuctionFactoryMBR(t *testing.T) {
	env := setupAuctionFactoryTest(t)

	result, err := env.client.SendMBR(env.ctx)
	if err != nil {
		t.Fatalf("failed to call mbr: %v", err)
	}

	t.Logf("AuctionFactory MBR - Bids: %d, Weights: %d, BidsByAddress: %d, Locations: %d",
		result.Return.Bids, result.Return.Weights,
		result.Return.BidsByAddress, result.Return.Locations)

	if result.Return.Bids == 0 {
		t.Error("expected non-zero bids MBR")
	}
}

func TestAuctionFactoryNewAuctionCost(t *testing.T) {
	env := setupAuctionFactoryTest(t)

	result, err := env.client.SendNewAuctionCost(env.ctx, algokit.CallParams[auctionfactory.NewAuctionCostArgs]{
		Args: auctionfactory.NewAuctionCostArgs{
			IsPrizeBox:       false,
			BidAssetID:       0,     // ALGO
			WeightsListCount: 10,
		},
	})
	if err != nil {
		t.Fatalf("failed to call newAuctionCost: %v", err)
	}

	t.Logf("AuctionFactory newAuctionCost: %d microAlgos", result.Return)

	if result.Return == 0 {
		t.Error("expected non-zero auction cost")
	}
}

func TestAuctionFactoryFactory(t *testing.T) {
	env := setupAuctionFactoryTest(t)

	client2, _, err := env.factory.Create(env.ctx, auctionfactory.FactoryCreateParams{
		Args: auctionfactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second auction factory: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second AuctionFactory deployed with appID: %d", client2.AppID())
}

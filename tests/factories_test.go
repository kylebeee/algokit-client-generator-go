package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/marketplace"
	"github.com/kylebeee/algokit-client-generator-go/generated/pollfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/rafflefactory"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- PollFactory Tests ---

type pollFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *pollfactory.Client
	factory *pollfactory.Factory
}

func setupPollFactoryTest(t *testing.T) *pollFactoryTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := pollfactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create poll factory: %v", err)
	}

	client, _, err := factory.Create(ctx, pollfactory.FactoryCreateParams{
		Args: pollfactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy poll factory contract: %v", err)
	}

	t.Logf("PollFactory contract deployed with appID: %d", client.AppID())

	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &pollFactoryTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestPollFactoryDeployment(t *testing.T) {
	env := setupPollFactoryTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("PollFactory deployed with appID: %d", env.client.AppID())
}

func TestPollFactoryNewPollCost(t *testing.T) {
	env := setupPollFactoryTest(t)

	result, err := env.client.SendNewPollCost(env.ctx)
	if err != nil {
		t.Fatalf("failed to call newPollCost: %v", err)
	}

	t.Logf("PollFactory newPollCost: %d microAlgos", result.Return)

	if result.Return == 0 {
		t.Error("expected non-zero poll cost")
	}
}

func TestPollFactoryFactory(t *testing.T) {
	env := setupPollFactoryTest(t)

	client2, _, err := env.factory.Create(env.ctx, pollfactory.FactoryCreateParams{
		Args: pollfactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second poll factory: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second PollFactory deployed with appID: %d", client2.AppID())
}

// --- Marketplace Tests ---

type marketplaceTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *marketplace.Client
	factory *marketplace.Factory
}

func setupMarketplaceTest(t *testing.T) *marketplaceTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := marketplace.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create marketplace factory: %v", err)
	}

	client, _, err := factory.Create(ctx, marketplace.FactoryCreateParams{
		Args: marketplace.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy marketplace contract: %v", err)
	}

	t.Logf("Marketplace contract deployed with appID: %d", client.AppID())

	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &marketplaceTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestMarketplaceDeployment(t *testing.T) {
	env := setupMarketplaceTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("Marketplace deployed with appID: %d", env.client.AppID())
}

func TestMarketplaceFactory(t *testing.T) {
	env := setupMarketplaceTest(t)

	client2, _, err := env.factory.Create(env.ctx, marketplace.FactoryCreateParams{
		Args: marketplace.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second marketplace: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second Marketplace deployed with appID: %d", client2.AppID())
}

// --- RaffleFactory Tests ---

type raffleFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *rafflefactory.Client
	factory *rafflefactory.Factory
}

func setupRaffleFactoryTest(t *testing.T) *raffleFactoryTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := rafflefactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create raffle factory: %v", err)
	}

	client, _, err := factory.Create(ctx, rafflefactory.FactoryCreateParams{
		Args: rafflefactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy raffle factory contract: %v", err)
	}

	t.Logf("RaffleFactory contract deployed with appID: %d", client.AppID())

	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &raffleFactoryTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestRaffleFactoryDeployment(t *testing.T) {
	env := setupRaffleFactoryTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("RaffleFactory deployed with appID: %d", env.client.AppID())
}

func TestRaffleFactoryMBR(t *testing.T) {
	env := setupRaffleFactoryTest(t)

	result, err := env.client.SendMBR(env.ctx)
	if err != nil {
		t.Fatalf("failed to call mbr: %v", err)
	}

	t.Logf("RaffleFactory MBR - Entries: %d, Weights: %d, EntriesByAddress: %d",
		result.Return.Entries, result.Return.Weights, result.Return.EntriesByAddress)

	if result.Return.Entries == 0 {
		t.Error("expected non-zero entries MBR")
	}
}

func TestRaffleFactoryFactory(t *testing.T) {
	env := setupRaffleFactoryTest(t)

	client2, _, err := env.factory.Create(env.ctx, rafflefactory.FactoryCreateParams{
		Args: rafflefactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second raffle factory: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second RaffleFactory deployed with appID: %d", client2.AppID())
}

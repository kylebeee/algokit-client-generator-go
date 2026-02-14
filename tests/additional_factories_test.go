package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/escrowfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/prizeboxfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/stakingpoolfactory"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- PrizeBoxFactory ---

type prizeBoxFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *prizeboxfactory.Client
	factory *prizeboxfactory.Factory
}

func setupPrizeBoxFactoryTest(t *testing.T) *prizeBoxFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := prizeboxfactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, prizeboxfactory.FactoryCreateParams{
		Args: prizeboxfactory.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &prizeBoxFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestPrizeBoxFactoryDeployment(t *testing.T) {
	env := setupPrizeBoxFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("PrizeBoxFactory deployed with appID: %d", env.client.AppID())
}

func TestPrizeBoxFactoryComposer(t *testing.T) {
	env := setupPrizeBoxFactoryTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("PrizeBoxFactory composer created successfully")
}

func TestPrizeBoxFactoryFactory(t *testing.T) {
	env := setupPrizeBoxFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, prizeboxfactory.FactoryCreateParams{
		Args: prizeboxfactory.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second PrizeBoxFactory deployed with appID: %d", client2.AppID())
}

// --- EscrowFactory ---

type escrowFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *escrowfactory.Client
	factory *escrowfactory.Factory
}

func setupEscrowFactoryTest(t *testing.T) *escrowFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := escrowfactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	return &escrowFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestEscrowFactoryDeployment(t *testing.T) {
	env := setupEscrowFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("EscrowFactory deployed with appID: %d", env.client.AppID())
}

func TestEscrowFactoryComposer(t *testing.T) {
	env := setupEscrowFactoryTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("EscrowFactory composer created successfully")
}

func TestEscrowFactoryFactory(t *testing.T) {
	env := setupEscrowFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second EscrowFactory deployed with appID: %d", client2.AppID())
}

// --- StakingPoolFactory ---

type stakingPoolFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *stakingpoolfactory.Client
	factory *stakingpoolfactory.Factory
}

func setupStakingPoolFactoryTest(t *testing.T) *stakingPoolFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := stakingpoolfactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, stakingpoolfactory.FactoryCreateParams{
		Args: stakingpoolfactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &stakingPoolFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestStakingPoolFactoryDeployment(t *testing.T) {
	env := setupStakingPoolFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("StakingPoolFactory deployed with appID: %d", env.client.AppID())
}

func TestStakingPoolFactoryComposer(t *testing.T) {
	env := setupStakingPoolFactoryTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("StakingPoolFactory composer created successfully")
}

func TestStakingPoolFactoryFactory(t *testing.T) {
	env := setupStakingPoolFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, stakingpoolfactory.FactoryCreateParams{
		Args: stakingpoolfactory.CreateArgs{
			Version:        "1.0.0",
			ChildVersion:   "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second StakingPoolFactory deployed with appID: %d", client2.AppID())
}

package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/subscriptions"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// subscriptionsTestEnv holds the shared state for Subscriptions tests.
type subscriptionsTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *subscriptions.Client
	factory *subscriptions.Factory
}

func setupSubscriptionsTest(t *testing.T) *subscriptionsTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := subscriptions.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create subscriptions factory: %v", err)
	}

	// Deploy with create(version, akitaDAO, akitaDAOEscrow) - pass 0s as dummy
	client, _, err := factory.Create(ctx, subscriptions.FactoryCreateParams{
		Args: subscriptions.CreateArgs{
			Version:        "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy subscriptions contract: %v", err)
	}

	t.Logf("Subscriptions contract deployed with appID: %d", client.AppID())

	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &subscriptionsTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestSubscriptionsDeployment(t *testing.T) {
	env := setupSubscriptionsTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("Subscriptions deployed with appID: %d", env.client.AppID())
}

func TestSubscriptionsBlockCost(t *testing.T) {
	env := setupSubscriptionsTest(t)

	result, err := env.client.SendBlockCost(env.ctx)
	if err != nil {
		t.Fatalf("failed to get block cost: %v", err)
	}

	t.Logf("Block cost: %d microAlgos", result.Return)

	if result.Return == 0 {
		t.Error("expected non-zero block cost")
	}
}

func TestSubscriptionsFactory(t *testing.T) {
	env := setupSubscriptionsTest(t)

	client2, _, err := env.factory.Create(env.ctx, subscriptions.FactoryCreateParams{
		Args: subscriptions.CreateArgs{
			Version:        "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second subscriptions: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second Subscriptions deployed with appID: %d", client2.AppID())
}

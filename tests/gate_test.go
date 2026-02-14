package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/gate"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// gateTestEnv holds the shared state for Gate tests.
type gateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *gate.Client
	factory *gate.Factory
}

func setupGateTest(t *testing.T) *gateTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := gate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create gate factory: %v", err)
	}

	// Deploy with create(version, akitaDAO) - pass 0 as dummy akitaDAO
	client, _, err := factory.Create(ctx, gate.FactoryCreateParams{
		Args: gate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy gate contract: %v", err)
	}

	t.Logf("Gate contract deployed with appID: %d", client.AppID())

	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &gateTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestGateDeployment(t *testing.T) {
	env := setupGateTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("Gate deployed with appID: %d", env.client.AppID())
}

func TestGateFactory(t *testing.T) {
	env := setupGateTest(t)

	client2, _, err := env.factory.Create(env.ctx, gate.FactoryCreateParams{
		Args: gate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second gate: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second Gate deployed with appID: %d", client2.AppID())
}

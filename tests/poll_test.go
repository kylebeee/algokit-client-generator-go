package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/daostub"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockrandomnessbeacon"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- DaoStub Tests ---

type daoStubTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *daostub.Client
	factory *daostub.Factory
}

func setupDaoStubTest(t *testing.T) *daoStubTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := daostub.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create daostub factory: %v", err)
	}

	// DaoStub uses bare create (no method args)
	client, _, err := factory.Create(ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy daostub contract: %v", err)
	}

	t.Logf("DaoStub contract deployed with appID: %d", client.AppID())

	return &daoStubTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestDaoStubDeployment(t *testing.T) {
	env := setupDaoStubTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("DaoStub deployed with appID: %d", env.client.AppID())
}

func TestDaoStubIsValidUpgrade(t *testing.T) {
	env := setupDaoStubTest(t)

	// Test isValidUpgrade with a dummy lease and app ID
	var lease [32]byte
	copy(lease[:], []byte("test-lease-for-upgrade-check!"))

	result, err := env.client.SendIsValidUpgrade(env.ctx, algokit.CallParams[daostub.IsValidUpgradeArgs]{
		Args: daostub.IsValidUpgradeArgs{
			Lease:            lease,
			AppBeingUpgraded: env.client.AppID(),
		},
	})
	if err != nil {
		t.Fatalf("failed to call isValidUpgrade: %v", err)
	}

	// DaoStub returns a bool - verify the call succeeded and we got a typed result
	t.Logf("DaoStub isValidUpgrade: %v", result.Return)
}

func TestDaoStubFactory(t *testing.T) {
	env := setupDaoStubTest(t)

	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second daostub: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second DaoStub deployed with appID: %d", client2.AppID())
}

// --- MockRandomnessBeacon Tests ---

type mockBeaconTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockrandomnessbeacon.Client
	factory *mockrandomnessbeacon.Factory
}

func setupMockBeaconTest(t *testing.T) *mockBeaconTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))

	factory, err := mockrandomnessbeacon.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create mock beacon factory: %v", err)
	}

	// MockRandomnessBeacon uses bare create
	client, _, err := factory.Create(ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy mock beacon contract: %v", err)
	}

	t.Logf("MockRandomnessBeacon contract deployed with appID: %d", client.AppID())

	return &mockBeaconTestEnv{
		fixture: fixture,
		ctx:     ctx,
		owner:   owner,
		client:  client,
		factory: factory,
	}
}

func TestMockBeaconDeployment(t *testing.T) {
	env := setupMockBeaconTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("MockRandomnessBeacon deployed with appID: %d", env.client.AppID())
}

func TestMockBeaconGet(t *testing.T) {
	env := setupMockBeaconTest(t)

	// Call get(round, userData) to get random bytes
	result, err := env.client.SendGet(env.ctx, algokit.CallParams[mockrandomnessbeacon.GetArgs]{
		Args: mockrandomnessbeacon.GetArgs{
			Round:    42,
			UserData: []byte("test-user-data"),
		},
	})
	if err != nil {
		t.Fatalf("failed to call get: %v", err)
	}

	// Mock may return empty or deterministic bytes - verify the call succeeded
	t.Logf("MockRandomnessBeacon get returned %d bytes: %x", len(result.Return), result.Return)
}

func TestMockBeaconFactory(t *testing.T) {
	env := setupMockBeaconTest(t)

	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second mock beacon: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second MockRandomnessBeacon deployed with appID: %d", client2.AppID())
}

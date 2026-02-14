package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/akitareferrergate"
	"github.com/kylebeee/algokit-client-generator-go/generated/assetgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/merkleaddressgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/merkleassetgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/nfdgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/nfdrootgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/pollgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/socialactivitygate"
	"github.com/kylebeee/algokit-client-generator-go/generated/socialfollowercountgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/socialfollowerindexgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/socialimpactgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/socialmoderatorgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/stakingamountgate"
	"github.com/kylebeee/algokit-client-generator-go/generated/stakingpowergate"
	"github.com/kylebeee/algokit-client-generator-go/generated/subscriptiongate"
	"github.com/kylebeee/algokit-client-generator-go/generated/subscriptionstreakgate"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// ============================================================================
// AssetGate
// ============================================================================

type assetGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *assetgate.Client
	factory *assetgate.Factory
}

func setupAssetGateTest(t *testing.T) *assetGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := assetgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, assetgate.FactoryCreateParams{
		Args: assetgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &assetGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAssetGateDeployment(t *testing.T) {
	env := setupAssetGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AssetGate deployed with appID: %d", env.client.AppID())
}

func TestAssetGateSendOpUp(t *testing.T) {
	env := setupAssetGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("AssetGate opUp succeeded")
}

func TestAssetGateFactory(t *testing.T) {
	env := setupAssetGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, assetgate.FactoryCreateParams{
		Args: assetgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second AssetGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// MerkleAddressGate
// ============================================================================

type merkleAddressGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *merkleaddressgate.Client
	factory *merkleaddressgate.Factory
}

func setupMerkleAddressGateTest(t *testing.T) *merkleAddressGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := merkleaddressgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, merkleaddressgate.FactoryCreateParams{
		Args: merkleaddressgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &merkleAddressGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMerkleAddressGateDeployment(t *testing.T) {
	env := setupMerkleAddressGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MerkleAddressGate deployed with appID: %d", env.client.AppID())
}

func TestMerkleAddressGateSendOpUp(t *testing.T) {
	env := setupMerkleAddressGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("MerkleAddressGate opUp succeeded")
}

func TestMerkleAddressGateFactory(t *testing.T) {
	env := setupMerkleAddressGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, merkleaddressgate.FactoryCreateParams{
		Args: merkleaddressgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MerkleAddressGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// MerkleAssetGate
// ============================================================================

type merkleAssetGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *merkleassetgate.Client
	factory *merkleassetgate.Factory
}

func setupMerkleAssetGateTest(t *testing.T) *merkleAssetGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := merkleassetgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, merkleassetgate.FactoryCreateParams{
		Args: merkleassetgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &merkleAssetGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMerkleAssetGateDeployment(t *testing.T) {
	env := setupMerkleAssetGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MerkleAssetGate deployed with appID: %d", env.client.AppID())
}

func TestMerkleAssetGateSendOpUp(t *testing.T) {
	env := setupMerkleAssetGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("MerkleAssetGate opUp succeeded")
}

func TestMerkleAssetGateFactory(t *testing.T) {
	env := setupMerkleAssetGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, merkleassetgate.FactoryCreateParams{
		Args: merkleassetgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MerkleAssetGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// NfdGate
// ============================================================================

type nfdGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *nfdgate.Client
	factory *nfdgate.Factory
}

func setupNfdGateTest(t *testing.T) *nfdGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := nfdgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, nfdgate.FactoryCreateParams{
		Args: nfdgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &nfdGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestNfdGateDeployment(t *testing.T) {
	env := setupNfdGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("NfdGate deployed with appID: %d", env.client.AppID())
}

func TestNfdGateSendOpUp(t *testing.T) {
	env := setupNfdGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("NfdGate opUp succeeded")
}

func TestNfdGateFactory(t *testing.T) {
	env := setupNfdGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, nfdgate.FactoryCreateParams{
		Args: nfdgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second NfdGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// NfdRootGate
// ============================================================================

type nfdRootGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *nfdrootgate.Client
	factory *nfdrootgate.Factory
}

func setupNfdRootGateTest(t *testing.T) *nfdRootGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := nfdrootgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, nfdrootgate.FactoryCreateParams{
		Args: nfdrootgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &nfdRootGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestNfdRootGateDeployment(t *testing.T) {
	env := setupNfdRootGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("NfdRootGate deployed with appID: %d", env.client.AppID())
}

func TestNfdRootGateSendOpUp(t *testing.T) {
	env := setupNfdRootGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("NfdRootGate opUp succeeded")
}

func TestNfdRootGateFactory(t *testing.T) {
	env := setupNfdRootGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, nfdrootgate.FactoryCreateParams{
		Args: nfdrootgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second NfdRootGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// PollGate
// ============================================================================

type pollGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *pollgate.Client
	factory *pollgate.Factory
}

func setupPollGateTest(t *testing.T) *pollGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := pollgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, pollgate.FactoryCreateParams{
		Args: pollgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &pollGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestPollGateDeployment(t *testing.T) {
	env := setupPollGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("PollGate deployed with appID: %d", env.client.AppID())
}

func TestPollGateSendOpUp(t *testing.T) {
	env := setupPollGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("PollGate opUp succeeded")
}

func TestPollGateFactory(t *testing.T) {
	env := setupPollGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, pollgate.FactoryCreateParams{
		Args: pollgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second PollGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SocialActivityGate
// ============================================================================

type socialActivityGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *socialactivitygate.Client
	factory *socialactivitygate.Factory
}

func setupSocialActivityGateTest(t *testing.T) *socialActivityGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := socialactivitygate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, socialactivitygate.FactoryCreateParams{
		Args: socialactivitygate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &socialActivityGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSocialActivityGateDeployment(t *testing.T) {
	env := setupSocialActivityGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SocialActivityGate deployed with appID: %d", env.client.AppID())
}

func TestSocialActivityGateSendOpUp(t *testing.T) {
	env := setupSocialActivityGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SocialActivityGate opUp succeeded")
}

func TestSocialActivityGateFactory(t *testing.T) {
	env := setupSocialActivityGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, socialactivitygate.FactoryCreateParams{
		Args: socialactivitygate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SocialActivityGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SocialFollowerCountGate
// ============================================================================

type socialFollowerCountGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *socialfollowercountgate.Client
	factory *socialfollowercountgate.Factory
}

func setupSocialFollowerCountGateTest(t *testing.T) *socialFollowerCountGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := socialfollowercountgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, socialfollowercountgate.FactoryCreateParams{
		Args: socialfollowercountgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &socialFollowerCountGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSocialFollowerCountGateDeployment(t *testing.T) {
	env := setupSocialFollowerCountGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SocialFollowerCountGate deployed with appID: %d", env.client.AppID())
}

func TestSocialFollowerCountGateSendOpUp(t *testing.T) {
	env := setupSocialFollowerCountGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SocialFollowerCountGate opUp succeeded")
}

func TestSocialFollowerCountGateFactory(t *testing.T) {
	env := setupSocialFollowerCountGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, socialfollowercountgate.FactoryCreateParams{
		Args: socialfollowercountgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SocialFollowerCountGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SocialFollowerIndexGate
// ============================================================================

type socialFollowerIndexGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *socialfollowerindexgate.Client
	factory *socialfollowerindexgate.Factory
}

func setupSocialFollowerIndexGateTest(t *testing.T) *socialFollowerIndexGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := socialfollowerindexgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, socialfollowerindexgate.FactoryCreateParams{
		Args: socialfollowerindexgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &socialFollowerIndexGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSocialFollowerIndexGateDeployment(t *testing.T) {
	env := setupSocialFollowerIndexGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SocialFollowerIndexGate deployed with appID: %d", env.client.AppID())
}

func TestSocialFollowerIndexGateSendOpUp(t *testing.T) {
	env := setupSocialFollowerIndexGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SocialFollowerIndexGate opUp succeeded")
}

func TestSocialFollowerIndexGateFactory(t *testing.T) {
	env := setupSocialFollowerIndexGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, socialfollowerindexgate.FactoryCreateParams{
		Args: socialfollowerindexgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SocialFollowerIndexGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SocialImpactGate
// ============================================================================

type socialImpactGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *socialimpactgate.Client
	factory *socialimpactgate.Factory
}

func setupSocialImpactGateTest(t *testing.T) *socialImpactGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := socialimpactgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, socialimpactgate.FactoryCreateParams{
		Args: socialimpactgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &socialImpactGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSocialImpactGateDeployment(t *testing.T) {
	env := setupSocialImpactGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SocialImpactGate deployed with appID: %d", env.client.AppID())
}

func TestSocialImpactGateSendOpUp(t *testing.T) {
	env := setupSocialImpactGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SocialImpactGate opUp succeeded")
}

func TestSocialImpactGateFactory(t *testing.T) {
	env := setupSocialImpactGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, socialimpactgate.FactoryCreateParams{
		Args: socialimpactgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SocialImpactGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SocialModeratorGate
// ============================================================================

type socialModeratorGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *socialmoderatorgate.Client
	factory *socialmoderatorgate.Factory
}

func setupSocialModeratorGateTest(t *testing.T) *socialModeratorGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := socialmoderatorgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, socialmoderatorgate.FactoryCreateParams{
		Args: socialmoderatorgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &socialModeratorGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSocialModeratorGateDeployment(t *testing.T) {
	env := setupSocialModeratorGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SocialModeratorGate deployed with appID: %d", env.client.AppID())
}

func TestSocialModeratorGateSendOpUp(t *testing.T) {
	env := setupSocialModeratorGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SocialModeratorGate opUp succeeded")
}

func TestSocialModeratorGateFactory(t *testing.T) {
	env := setupSocialModeratorGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, socialmoderatorgate.FactoryCreateParams{
		Args: socialmoderatorgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SocialModeratorGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// StakingAmountGate
// ============================================================================

type stakingAmountGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *stakingamountgate.Client
	factory *stakingamountgate.Factory
}

func setupStakingAmountGateTest(t *testing.T) *stakingAmountGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := stakingamountgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, stakingamountgate.FactoryCreateParams{
		Args: stakingamountgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &stakingAmountGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestStakingAmountGateDeployment(t *testing.T) {
	env := setupStakingAmountGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("StakingAmountGate deployed with appID: %d", env.client.AppID())
}

func TestStakingAmountGateSendOpUp(t *testing.T) {
	env := setupStakingAmountGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("StakingAmountGate opUp succeeded")
}

func TestStakingAmountGateFactory(t *testing.T) {
	env := setupStakingAmountGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, stakingamountgate.FactoryCreateParams{
		Args: stakingamountgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second StakingAmountGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// StakingPowerGate
// ============================================================================

type stakingPowerGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *stakingpowergate.Client
	factory *stakingpowergate.Factory
}

func setupStakingPowerGateTest(t *testing.T) *stakingPowerGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := stakingpowergate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, stakingpowergate.FactoryCreateParams{
		Args: stakingpowergate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &stakingPowerGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestStakingPowerGateDeployment(t *testing.T) {
	env := setupStakingPowerGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("StakingPowerGate deployed with appID: %d", env.client.AppID())
}

func TestStakingPowerGateSendOpUp(t *testing.T) {
	env := setupStakingPowerGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("StakingPowerGate opUp succeeded")
}

func TestStakingPowerGateFactory(t *testing.T) {
	env := setupStakingPowerGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, stakingpowergate.FactoryCreateParams{
		Args: stakingpowergate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second StakingPowerGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SubscriptionGate
// ============================================================================

type subscriptionGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *subscriptiongate.Client
	factory *subscriptiongate.Factory
}

func setupSubscriptionGateTest(t *testing.T) *subscriptionGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := subscriptiongate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, subscriptiongate.FactoryCreateParams{
		Args: subscriptiongate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &subscriptionGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSubscriptionGateDeployment(t *testing.T) {
	env := setupSubscriptionGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SubscriptionGate deployed with appID: %d", env.client.AppID())
}

func TestSubscriptionGateSendOpUp(t *testing.T) {
	env := setupSubscriptionGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SubscriptionGate opUp succeeded")
}

func TestSubscriptionGateFactory(t *testing.T) {
	env := setupSubscriptionGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, subscriptiongate.FactoryCreateParams{
		Args: subscriptiongate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SubscriptionGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// SubscriptionStreakGate
// ============================================================================

type subscriptionStreakGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *subscriptionstreakgate.Client
	factory *subscriptionstreakgate.Factory
}

func setupSubscriptionStreakGateTest(t *testing.T) *subscriptionStreakGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := subscriptionstreakgate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, subscriptionstreakgate.FactoryCreateParams{
		Args: subscriptionstreakgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &subscriptionStreakGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSubscriptionStreakGateDeployment(t *testing.T) {
	env := setupSubscriptionStreakGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SubscriptionStreakGate deployed with appID: %d", env.client.AppID())
}

func TestSubscriptionStreakGateSendOpUp(t *testing.T) {
	env := setupSubscriptionStreakGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("SubscriptionStreakGate opUp succeeded")
}

func TestSubscriptionStreakGateFactory(t *testing.T) {
	env := setupSubscriptionStreakGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, subscriptionstreakgate.FactoryCreateParams{
		Args: subscriptionstreakgate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SubscriptionStreakGate deployed with appID: %d", client2.AppID())
}

// ============================================================================
// AkitaReferrerGate
// ============================================================================

type akitaReferrerGateTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitareferrergate.Client
	factory *akitareferrergate.Factory
}

func setupAkitaReferrerGateTest(t *testing.T) *akitaReferrerGateTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitareferrergate.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitareferrergate.FactoryCreateParams{
		Args: akitareferrergate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaReferrerGateTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaReferrerGateDeployment(t *testing.T) {
	env := setupAkitaReferrerGateTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaReferrerGate deployed with appID: %d", env.client.AppID())
}

func TestAkitaReferrerGateSendOpUp(t *testing.T) {
	env := setupAkitaReferrerGateTest(t)
	err := env.client.SendOpUp(env.ctx)
	if err != nil {
		t.Fatalf("failed to call opUp: %v", err)
	}
	t.Log("AkitaReferrerGate opUp succeeded")
}

func TestAkitaReferrerGateFactory(t *testing.T) {
	env := setupAkitaReferrerGateTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitareferrergate.FactoryCreateParams{
		Args: akitareferrergate.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 1,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second AkitaReferrerGate deployed with appID: %d", client2.AppID())
}

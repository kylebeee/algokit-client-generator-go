package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/mockabstractedaccountfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockakitadao"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockakitasocial"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockauctionfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockmarketplace"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockpollfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockprizeboxfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockrafflefactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/mockstakingpoolfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/mocksubscriptions"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- MockAuctionFactory Tests ---

type mockAuctionFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockauctionfactory.Client
	factory *mockauctionfactory.Factory
}

func setupMockAuctionFactoryTest(t *testing.T) *mockAuctionFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockauctionfactory.NewFactory(algokit.AppFactoryParams{
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
	return &mockAuctionFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockAuctionFactoryDeployment(t *testing.T) {
	env := setupMockAuctionFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockAuctionFactory deployed with appID: %d", env.client.AppID())
}

func TestMockAuctionFactorySendPing(t *testing.T) {
	env := setupMockAuctionFactoryTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockAuctionFactory ping returned: %d", result.Return)
}

func TestMockAuctionFactoryFactory(t *testing.T) {
	env := setupMockAuctionFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockAuctionFactory deployed with appID: %d", client2.AppID())
}

// --- MockPollFactory Tests ---

type mockPollFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockpollfactory.Client
	factory *mockpollfactory.Factory
}

func setupMockPollFactoryTest(t *testing.T) *mockPollFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockpollfactory.NewFactory(algokit.AppFactoryParams{
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
	return &mockPollFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockPollFactoryDeployment(t *testing.T) {
	env := setupMockPollFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockPollFactory deployed with appID: %d", env.client.AppID())
}

func TestMockPollFactorySendPing(t *testing.T) {
	env := setupMockPollFactoryTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockPollFactory ping returned: %d", result.Return)
}

func TestMockPollFactoryFactory(t *testing.T) {
	env := setupMockPollFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockPollFactory deployed with appID: %d", client2.AppID())
}

// --- MockRaffleFactory Tests ---

type mockRaffleFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockrafflefactory.Client
	factory *mockrafflefactory.Factory
}

func setupMockRaffleFactoryTest(t *testing.T) *mockRaffleFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockrafflefactory.NewFactory(algokit.AppFactoryParams{
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
	return &mockRaffleFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockRaffleFactoryDeployment(t *testing.T) {
	env := setupMockRaffleFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockRaffleFactory deployed with appID: %d", env.client.AppID())
}

func TestMockRaffleFactorySendPing(t *testing.T) {
	env := setupMockRaffleFactoryTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockRaffleFactory ping returned: %d", result.Return)
}

func TestMockRaffleFactoryFactory(t *testing.T) {
	env := setupMockRaffleFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockRaffleFactory deployed with appID: %d", client2.AppID())
}

// --- MockPrizeBoxFactory Tests ---

type mockPrizeBoxFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockprizeboxfactory.Client
	factory *mockprizeboxfactory.Factory
}

func setupMockPrizeBoxFactoryTest(t *testing.T) *mockPrizeBoxFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockprizeboxfactory.NewFactory(algokit.AppFactoryParams{
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
	return &mockPrizeBoxFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockPrizeBoxFactoryDeployment(t *testing.T) {
	env := setupMockPrizeBoxFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockPrizeBoxFactory deployed with appID: %d", env.client.AppID())
}

func TestMockPrizeBoxFactorySendPing(t *testing.T) {
	env := setupMockPrizeBoxFactoryTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockPrizeBoxFactory ping returned: %d", result.Return)
}

func TestMockPrizeBoxFactoryFactory(t *testing.T) {
	env := setupMockPrizeBoxFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockPrizeBoxFactory deployed with appID: %d", client2.AppID())
}

// --- MockStakingPoolFactory Tests ---

type mockStakingPoolFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockstakingpoolfactory.Client
	factory *mockstakingpoolfactory.Factory
}

func setupMockStakingPoolFactoryTest(t *testing.T) *mockStakingPoolFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockstakingpoolfactory.NewFactory(algokit.AppFactoryParams{
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
	return &mockStakingPoolFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockStakingPoolFactoryDeployment(t *testing.T) {
	env := setupMockStakingPoolFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockStakingPoolFactory deployed with appID: %d", env.client.AppID())
}

func TestMockStakingPoolFactorySendPing(t *testing.T) {
	env := setupMockStakingPoolFactoryTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockStakingPoolFactory ping returned: %d", result.Return)
}

func TestMockStakingPoolFactoryFactory(t *testing.T) {
	env := setupMockStakingPoolFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockStakingPoolFactory deployed with appID: %d", client2.AppID())
}

// --- MockAbstractedAccountFactory Tests ---

type mockAbstractedAccountFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockabstractedaccountfactory.Client
	factory *mockabstractedaccountfactory.Factory
}

func setupMockAbstractedAccountFactoryTest(t *testing.T) *mockAbstractedAccountFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockabstractedaccountfactory.NewFactory(algokit.AppFactoryParams{
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
	return &mockAbstractedAccountFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockAbstractedAccountFactoryDeployment(t *testing.T) {
	env := setupMockAbstractedAccountFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockAbstractedAccountFactory deployed with appID: %d", env.client.AppID())
}

func TestMockAbstractedAccountFactorySendPing(t *testing.T) {
	env := setupMockAbstractedAccountFactoryTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockAbstractedAccountFactory ping returned: %d", result.Return)
}

func TestMockAbstractedAccountFactoryFactory(t *testing.T) {
	env := setupMockAbstractedAccountFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockAbstractedAccountFactory deployed with appID: %d", client2.AppID())
}

// --- MockAkitaDAO Tests ---

type mockAkitaDAOTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockakitadao.Client
	factory *mockakitadao.Factory
}

func setupMockAkitaDAOTest(t *testing.T) *mockAkitaDAOTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockakitadao.NewFactory(algokit.AppFactoryParams{
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
	return &mockAkitaDAOTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockAkitaDAODeployment(t *testing.T) {
	env := setupMockAkitaDAOTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockAkitaDAO deployed with appID: %d", env.client.AppID())
}

func TestMockAkitaDAOSendPing(t *testing.T) {
	env := setupMockAkitaDAOTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockAkitaDAO ping returned: %d", result.Return)
}

func TestMockAkitaDAOFactory(t *testing.T) {
	env := setupMockAkitaDAOTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockAkitaDAO deployed with appID: %d", client2.AppID())
}

// --- MockAkitaSocial Tests ---

type mockAkitaSocialTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockakitasocial.Client
	factory *mockakitasocial.Factory
}

func setupMockAkitaSocialTest(t *testing.T) *mockAkitaSocialTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockakitasocial.NewFactory(algokit.AppFactoryParams{
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
	return &mockAkitaSocialTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockAkitaSocialDeployment(t *testing.T) {
	env := setupMockAkitaSocialTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockAkitaSocial deployed with appID: %d", env.client.AppID())
}

func TestMockAkitaSocialSendPing(t *testing.T) {
	env := setupMockAkitaSocialTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockAkitaSocial ping returned: %d", result.Return)
}

func TestMockAkitaSocialFactory(t *testing.T) {
	env := setupMockAkitaSocialTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockAkitaSocial deployed with appID: %d", client2.AppID())
}

// --- MockMarketplace Tests ---

type mockMarketplaceTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mockmarketplace.Client
	factory *mockmarketplace.Factory
}

func setupMockMarketplaceTest(t *testing.T) *mockMarketplaceTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mockmarketplace.NewFactory(algokit.AppFactoryParams{
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
	return &mockMarketplaceTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockMarketplaceDeployment(t *testing.T) {
	env := setupMockMarketplaceTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockMarketplace deployed with appID: %d", env.client.AppID())
}

func TestMockMarketplaceSendPing(t *testing.T) {
	env := setupMockMarketplaceTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockMarketplace ping returned: %d", result.Return)
}

func TestMockMarketplaceFactory(t *testing.T) {
	env := setupMockMarketplaceTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockMarketplace deployed with appID: %d", client2.AppID())
}

// --- MockSubscriptions Tests ---

type mockSubscriptionsTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *mocksubscriptions.Client
	factory *mocksubscriptions.Factory
}

func setupMockSubscriptionsTest(t *testing.T) *mockSubscriptionsTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := mocksubscriptions.NewFactory(algokit.AppFactoryParams{
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
	return &mockSubscriptionsTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMockSubscriptionsDeployment(t *testing.T) {
	env := setupMockSubscriptionsTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MockSubscriptions deployed with appID: %d", env.client.AppID())
}

func TestMockSubscriptionsSendPing(t *testing.T) {
	env := setupMockSubscriptionsTest(t)
	result, err := env.client.SendPing(env.ctx)
	if err != nil {
		t.Fatalf("failed to call ping: %v", err)
	}
	t.Logf("MockSubscriptions ping returned: %d", result.Return)
}

func TestMockSubscriptionsFactory(t *testing.T) {
	env := setupMockSubscriptionsTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second MockSubscriptions deployed with appID: %d", client2.AppID())
}

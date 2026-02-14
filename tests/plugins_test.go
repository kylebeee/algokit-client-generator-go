package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/akitadaoplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/akitasocialplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/auctionplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/dualstakeplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/gateplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/hyperswapplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/marketplaceplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/nfdplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/paysiloplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/paysilofactoryplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/pollplugincontract"
	"github.com/kylebeee/algokit-client-generator-go/generated/raffleplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/revenuemanagerplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/rewardsplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/stakingplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/stakingpoolplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/subscriptionsplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/updateakitadaoplugin"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- AkitaDAOPlugin ---

type akitaDaoPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitadaoplugin.Client
	factory *akitadaoplugin.Factory
}

func setupAkitaDaoPluginTest(t *testing.T) *akitaDaoPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitadaoplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitadaoplugin.FactoryCreateParams{
		Args:       akitadaoplugin.CreateArgs{DaoAppID: 0},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaDaoPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaDaoPluginDeployment(t *testing.T) {
	env := setupAkitaDaoPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaDaoPlugin deployed with appID: %d", env.client.AppID())
}

func TestAkitaDaoPluginComposer(t *testing.T) {
	env := setupAkitaDaoPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaDaoPlugin composer created")
}

func TestAkitaDaoPluginFactory(t *testing.T) {
	env := setupAkitaDaoPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitadaoplugin.FactoryCreateParams{
		Args:       akitadaoplugin.CreateArgs{DaoAppID: 0},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second AkitaDaoPlugin deployed with appID: %d", client2.AppID())
}

// --- AuctionPlugin ---

type auctionPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *auctionplugin.Client
	factory *auctionplugin.Factory
}

func setupAuctionPluginTest(t *testing.T) *auctionPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := auctionplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, auctionplugin.FactoryCreateParams{
		Args: auctionplugin.CreateArgs{
			Version:  "1.0.0",
			Factory:  0,
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &auctionPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAuctionPluginDeployment(t *testing.T) {
	env := setupAuctionPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AuctionPlugin deployed with appID: %d", env.client.AppID())
}

func TestAuctionPluginComposer(t *testing.T) {
	env := setupAuctionPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AuctionPlugin composer created")
}

func TestAuctionPluginFactory(t *testing.T) {
	env := setupAuctionPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, auctionplugin.FactoryCreateParams{
		Args: auctionplugin.CreateArgs{
			Version:  "1.0.0",
			Factory:  0,
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
	t.Logf("Second AuctionPlugin deployed with appID: %d", client2.AppID())
}

// --- DualStakePlugin ---

type dualStakePluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *dualstakeplugin.Client
	factory *dualstakeplugin.Factory
}

func setupDualStakePluginTest(t *testing.T) *dualStakePluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := dualstakeplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, dualstakeplugin.FactoryCreateParams{
		Args: dualstakeplugin.CreateArgs{Registry: 0},
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &dualStakePluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestDualStakePluginDeployment(t *testing.T) {
	env := setupDualStakePluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("DualStakePlugin deployed with appID: %d", env.client.AppID())
}

func TestDualStakePluginComposer(t *testing.T) {
	env := setupDualStakePluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("DualStakePlugin composer created")
}

func TestDualStakePluginFactory(t *testing.T) {
	env := setupDualStakePluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, dualstakeplugin.FactoryCreateParams{
		Args: dualstakeplugin.CreateArgs{Registry: 0},
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second DualStakePlugin deployed with appID: %d", client2.AppID())
}

// --- GatePlugin ---

type gatePluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *gateplugin.Client
	factory *gateplugin.Factory
}

func setupGatePluginTest(t *testing.T) *gatePluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := gateplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, gateplugin.FactoryCreateParams{
		Args: gateplugin.CreateArgs{GateAppID: 0},
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &gatePluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestGatePluginDeployment(t *testing.T) {
	env := setupGatePluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("GatePlugin deployed with appID: %d", env.client.AppID())
}

func TestGatePluginComposer(t *testing.T) {
	env := setupGatePluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("GatePlugin composer created")
}

func TestGatePluginFactory(t *testing.T) {
	env := setupGatePluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, gateplugin.FactoryCreateParams{
		Args: gateplugin.CreateArgs{GateAppID: 0},
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second GatePlugin deployed with appID: %d", client2.AppID())
}

// --- HyperSwapPlugin ---

type hyperSwapPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *hyperswapplugin.Client
	factory *hyperswapplugin.Factory
}

func setupHyperSwapPluginTest(t *testing.T) *hyperSwapPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := hyperswapplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, hyperswapplugin.FactoryCreateParams{
		Args: hyperswapplugin.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &hyperSwapPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestHyperSwapPluginDeployment(t *testing.T) {
	env := setupHyperSwapPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("HyperSwapPlugin deployed with appID: %d", env.client.AppID())
}

func TestHyperSwapPluginComposer(t *testing.T) {
	env := setupHyperSwapPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("HyperSwapPlugin composer created")
}

func TestHyperSwapPluginFactory(t *testing.T) {
	env := setupHyperSwapPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, hyperswapplugin.FactoryCreateParams{
		Args: hyperswapplugin.CreateArgs{
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
	t.Logf("Second HyperSwapPlugin deployed with appID: %d", client2.AppID())
}

// --- MarketplacePlugin ---

type marketplacePluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *marketplaceplugin.Client
	factory *marketplaceplugin.Factory
}

func setupMarketplacePluginTest(t *testing.T) *marketplacePluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := marketplaceplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, marketplaceplugin.FactoryCreateParams{
		Args: marketplaceplugin.CreateApplicationArgs{
			Version:  "1.0.0",
			Factory:  0,
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &marketplacePluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestMarketplacePluginDeployment(t *testing.T) {
	env := setupMarketplacePluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("MarketplacePlugin deployed with appID: %d", env.client.AppID())
}

func TestMarketplacePluginComposer(t *testing.T) {
	env := setupMarketplacePluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("MarketplacePlugin composer created")
}

func TestMarketplacePluginFactory(t *testing.T) {
	env := setupMarketplacePluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, marketplaceplugin.FactoryCreateParams{
		Args: marketplaceplugin.CreateApplicationArgs{
			Version:  "1.0.0",
			Factory:  0,
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
	t.Logf("Second MarketplacePlugin deployed with appID: %d", client2.AppID())
}

// --- NFDPlugin ---

type nfdPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *nfdplugin.Client
	factory *nfdplugin.Factory
}

func setupNfdPluginTest(t *testing.T) *nfdPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := nfdplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, nfdplugin.FactoryCreateParams{
		Args:       nfdplugin.CreateArgs{Registry: 0},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &nfdPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestNfdPluginDeployment(t *testing.T) {
	env := setupNfdPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("NFDPlugin deployed with appID: %d", env.client.AppID())
}

func TestNfdPluginComposer(t *testing.T) {
	env := setupNfdPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("NFDPlugin composer created")
}

func TestNfdPluginFactory(t *testing.T) {
	env := setupNfdPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, nfdplugin.FactoryCreateParams{
		Args:       nfdplugin.CreateArgs{Registry: 0},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second NFDPlugin deployed with appID: %d", client2.AppID())
}

// --- PaySiloPlugin ---

type paySiloPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *paysiloplugin.Client
	factory *paysiloplugin.Factory
}

func setupPaySiloPluginTest(t *testing.T) *paySiloPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := paysiloplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, paysiloplugin.FactoryCreateParams{
		Args: paysiloplugin.CreateArgs{Recipient: owner.Address},
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &paySiloPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestPaySiloPluginDeployment(t *testing.T) {
	env := setupPaySiloPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("PaySiloPlugin deployed with appID: %d", env.client.AppID())
}

func TestPaySiloPluginComposer(t *testing.T) {
	env := setupPaySiloPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("PaySiloPlugin composer created")
}

func TestPaySiloPluginFactory(t *testing.T) {
	env := setupPaySiloPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, paysiloplugin.FactoryCreateParams{
		Args: paysiloplugin.CreateArgs{Recipient: env.owner.Address},
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second PaySiloPlugin deployed with appID: %d", client2.AppID())
}

// --- PaySiloFactoryPlugin ---

type paySiloFactoryPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *paysilofactoryplugin.Client
	factory *paysilofactoryplugin.Factory
}

func setupPaySiloFactoryPluginTest(t *testing.T) *paySiloFactoryPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := paysilofactoryplugin.NewFactory(algokit.AppFactoryParams{
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
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &paySiloFactoryPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestPaySiloFactoryPluginDeployment(t *testing.T) {
	env := setupPaySiloFactoryPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("PaySiloFactoryPlugin deployed with appID: %d", env.client.AppID())
}

func TestPaySiloFactoryPluginComposer(t *testing.T) {
	env := setupPaySiloFactoryPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("PaySiloFactoryPlugin composer created")
}

func TestPaySiloFactoryPluginFactory(t *testing.T) {
	env := setupPaySiloFactoryPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second PaySiloFactoryPlugin deployed with appID: %d", client2.AppID())
}

// --- PollPluginContract ---

type pollPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *pollplugincontract.Client
	factory *pollplugincontract.Factory
}

func setupPollPluginTest(t *testing.T) *pollPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := pollplugincontract.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, pollplugincontract.FactoryCreateParams{
		Args: pollplugincontract.CreateArgs{
			Version:  "1.0.0",
			Factory:  0,
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &pollPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestPollPluginDeployment(t *testing.T) {
	env := setupPollPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("PollPluginContract deployed with appID: %d", env.client.AppID())
}

func TestPollPluginComposer(t *testing.T) {
	env := setupPollPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("PollPluginContract composer created")
}

func TestPollPluginFactory(t *testing.T) {
	env := setupPollPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, pollplugincontract.FactoryCreateParams{
		Args: pollplugincontract.CreateArgs{
			Version:  "1.0.0",
			Factory:  0,
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
	t.Logf("Second PollPluginContract deployed with appID: %d", client2.AppID())
}

// --- RafflePlugin ---

type rafflePluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *raffleplugin.Client
	factory *raffleplugin.Factory
}

func setupRafflePluginTest(t *testing.T) *rafflePluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := raffleplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, raffleplugin.FactoryCreateParams{
		Args: raffleplugin.CreateArgs{
			Version: "1.0.0",
			Factory: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &rafflePluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestRafflePluginDeployment(t *testing.T) {
	env := setupRafflePluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("RafflePlugin deployed with appID: %d", env.client.AppID())
}

func TestRafflePluginComposer(t *testing.T) {
	env := setupRafflePluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("RafflePlugin composer created")
}

func TestRafflePluginFactory(t *testing.T) {
	env := setupRafflePluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, raffleplugin.FactoryCreateParams{
		Args: raffleplugin.CreateArgs{
			Version: "1.0.0",
			Factory: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second RafflePlugin deployed with appID: %d", client2.AppID())
}

// --- RevenueManagerPlugin ---

type revenueManagerPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *revenuemanagerplugin.Client
	factory *revenuemanagerplugin.Factory
}

func setupRevenueManagerPluginTest(t *testing.T) *revenueManagerPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := revenuemanagerplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, revenuemanagerplugin.FactoryCreateParams{
		Args: revenuemanagerplugin.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &revenueManagerPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestRevenueManagerPluginDeployment(t *testing.T) {
	env := setupRevenueManagerPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("RevenueManagerPlugin deployed with appID: %d", env.client.AppID())
}

func TestRevenueManagerPluginComposer(t *testing.T) {
	env := setupRevenueManagerPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("RevenueManagerPlugin composer created")
}

func TestRevenueManagerPluginFactory(t *testing.T) {
	env := setupRevenueManagerPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, revenuemanagerplugin.FactoryCreateParams{
		Args: revenuemanagerplugin.CreateArgs{
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
	t.Logf("Second RevenueManagerPlugin deployed with appID: %d", client2.AppID())
}

// --- RewardsPlugin ---

type rewardsPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *rewardsplugin.Client
	factory *rewardsplugin.Factory
}

func setupRewardsPluginTest(t *testing.T) *rewardsPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := rewardsplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, rewardsplugin.FactoryCreateParams{
		Args: rewardsplugin.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &rewardsPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestRewardsPluginDeployment(t *testing.T) {
	env := setupRewardsPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("RewardsPlugin deployed with appID: %d", env.client.AppID())
}

func TestRewardsPluginComposer(t *testing.T) {
	env := setupRewardsPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("RewardsPlugin composer created")
}

func TestRewardsPluginFactory(t *testing.T) {
	env := setupRewardsPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, rewardsplugin.FactoryCreateParams{
		Args: rewardsplugin.CreateArgs{
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
	t.Logf("Second RewardsPlugin deployed with appID: %d", client2.AppID())
}

// --- StakingPlugin ---

type stakingPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *stakingplugin.Client
	factory *stakingplugin.Factory
}

func setupStakingPluginTest(t *testing.T) *stakingPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := stakingplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, stakingplugin.FactoryCreateParams{
		Args: stakingplugin.CreateArgs{
			AkitaDao: 0,
			Version:  "1.0.0",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &stakingPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestStakingPluginDeployment(t *testing.T) {
	env := setupStakingPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("StakingPlugin deployed with appID: %d", env.client.AppID())
}

func TestStakingPluginComposer(t *testing.T) {
	env := setupStakingPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("StakingPlugin composer created")
}

func TestStakingPluginFactory(t *testing.T) {
	env := setupStakingPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, stakingplugin.FactoryCreateParams{
		Args: stakingplugin.CreateArgs{
			AkitaDao: 0,
			Version:  "1.0.0",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second StakingPlugin deployed with appID: %d", client2.AppID())
}

// --- StakingPoolPlugin ---

type stakingPoolPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *stakingpoolplugin.Client
	factory *stakingpoolplugin.Factory
}

func setupStakingPoolPluginTest(t *testing.T) *stakingPoolPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := stakingpoolplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, stakingpoolplugin.FactoryCreateParams{
		Args: stakingpoolplugin.CreateArgs{
			Version:  "1.0.0",
			Factory:  0,
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &stakingPoolPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestStakingPoolPluginDeployment(t *testing.T) {
	env := setupStakingPoolPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("StakingPoolPlugin deployed with appID: %d", env.client.AppID())
}

func TestStakingPoolPluginComposer(t *testing.T) {
	env := setupStakingPoolPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("StakingPoolPlugin composer created")
}

func TestStakingPoolPluginFactory(t *testing.T) {
	env := setupStakingPoolPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, stakingpoolplugin.FactoryCreateParams{
		Args: stakingpoolplugin.CreateArgs{
			Version:  "1.0.0",
			Factory:  0,
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
	t.Logf("Second StakingPoolPlugin deployed with appID: %d", client2.AppID())
}

// --- SubscriptionsPlugin ---

type subscriptionsPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *subscriptionsplugin.Client
	factory *subscriptionsplugin.Factory
}

func setupSubscriptionsPluginTest(t *testing.T) *subscriptionsPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := subscriptionsplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, subscriptionsplugin.FactoryCreateParams{
		Args: subscriptionsplugin.CreateArgs{
			AkitaDao: 0,
			Version:  "1.0.0",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &subscriptionsPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestSubscriptionsPluginDeployment(t *testing.T) {
	env := setupSubscriptionsPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("SubscriptionsPlugin deployed with appID: %d", env.client.AppID())
}

func TestSubscriptionsPluginComposer(t *testing.T) {
	env := setupSubscriptionsPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("SubscriptionsPlugin composer created")
}

func TestSubscriptionsPluginFactory(t *testing.T) {
	env := setupSubscriptionsPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, subscriptionsplugin.FactoryCreateParams{
		Args: subscriptionsplugin.CreateArgs{
			AkitaDao: 0,
			Version:  "1.0.0",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second SubscriptionsPlugin deployed with appID: %d", client2.AppID())
}

// --- UpdateAkitaDAOPlugin ---

type updateAkitaDaoPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *updateakitadaoplugin.Client
	factory *updateakitadaoplugin.Factory
}

func setupUpdateAkitaDaoPluginTest(t *testing.T) *updateAkitaDaoPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := updateakitadaoplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, updateakitadaoplugin.FactoryCreateParams{
		Args: updateakitadaoplugin.CreateArgs{
			AkitaDao:     0,
			ClearProgram: []byte{},
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &updateAkitaDaoPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestUpdateAkitaDaoPluginDeployment(t *testing.T) {
	env := setupUpdateAkitaDaoPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("UpdateAkitaDaoPlugin deployed with appID: %d", env.client.AppID())
}

func TestUpdateAkitaDaoPluginComposer(t *testing.T) {
	env := setupUpdateAkitaDaoPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("UpdateAkitaDaoPlugin composer created")
}

func TestUpdateAkitaDaoPluginFactory(t *testing.T) {
	env := setupUpdateAkitaDaoPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, updateakitadaoplugin.FactoryCreateParams{
		Args: updateakitadaoplugin.CreateArgs{
			AkitaDao:     0,
			ClearProgram: []byte{},
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second UpdateAkitaDaoPlugin deployed with appID: %d", client2.AppID())
}

// --- AkitaSocialPlugin ---

type akitaSocialPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitasocialplugin.Client
	factory *akitasocialplugin.Factory
}

func setupAkitaSocialPluginTest(t *testing.T) *akitaSocialPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitasocialplugin.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitasocialplugin.FactoryCreateParams{
		Args: akitasocialplugin.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
			Escrow:   0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaSocialPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaSocialPluginDeployment(t *testing.T) {
	env := setupAkitaSocialPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaSocialPlugin deployed with appID: %d", env.client.AppID())
}

func TestAkitaSocialPluginComposer(t *testing.T) {
	env := setupAkitaSocialPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaSocialPlugin composer created")
}

func TestAkitaSocialPluginFactory(t *testing.T) {
	env := setupAkitaSocialPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitasocialplugin.FactoryCreateParams{
		Args: akitasocialplugin.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
			Escrow:   0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second AkitaSocialPlugin deployed with appID: %d", client2.AppID())
}

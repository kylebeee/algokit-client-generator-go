package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/akitadaotypes"
	"github.com/kylebeee/algokit-client-generator-go/generated/asamintplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/optinplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/payplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/testcloseoutplugin"
	"github.com/kylebeee/algokit-client-generator-go/generated/testproxyrekeyplugin"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- OptInPlugin ---

type optinPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *optinplugin.Client
	factory *optinplugin.Factory
}

func setupOptinPluginTest(t *testing.T) *optinPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := optinplugin.NewFactory(algokit.AppFactoryParams{
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
	return &optinPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestOptinPluginDeployment(t *testing.T) {
	env := setupOptinPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("OptinPlugin deployed with appID: %d", env.client.AppID())
}

func TestOptinPluginComposer(t *testing.T) {
	env := setupOptinPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("OptinPlugin composer created successfully")
}

func TestOptinPluginFactory(t *testing.T) {
	env := setupOptinPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second OptinPlugin deployed with appID: %d", client2.AppID())
}

// --- PayPlugin ---

type payPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *payplugin.Client
	factory *payplugin.Factory
}

func setupPayPluginTest(t *testing.T) *payPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := payplugin.NewFactory(algokit.AppFactoryParams{
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
	return &payPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestPayPluginDeployment(t *testing.T) {
	env := setupPayPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("PayPlugin deployed with appID: %d", env.client.AppID())
}

func TestPayPluginComposer(t *testing.T) {
	env := setupPayPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("PayPlugin composer created successfully")
}

func TestPayPluginFactory(t *testing.T) {
	env := setupPayPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second PayPlugin deployed with appID: %d", client2.AppID())
}

// --- TestCloseOutPlugin ---

type testCloseOutPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *testcloseoutplugin.Client
	factory *testcloseoutplugin.Factory
}

func setupTestCloseOutPluginTest(t *testing.T) *testCloseOutPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := testcloseoutplugin.NewFactory(algokit.AppFactoryParams{
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
	return &testCloseOutPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestTestCloseOutPluginDeployment(t *testing.T) {
	env := setupTestCloseOutPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("TestCloseOutPlugin deployed with appID: %d", env.client.AppID())
}

func TestTestCloseOutPluginComposer(t *testing.T) {
	env := setupTestCloseOutPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("TestCloseOutPlugin composer created successfully")
}

func TestTestCloseOutPluginFactory(t *testing.T) {
	env := setupTestCloseOutPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second TestCloseOutPlugin deployed with appID: %d", client2.AppID())
}

// --- TestProxyRekeyPlugin ---

type testProxyRekeyPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *testproxyrekeyplugin.Client
	factory *testproxyrekeyplugin.Factory
}

func setupTestProxyRekeyPluginTest(t *testing.T) *testProxyRekeyPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := testproxyrekeyplugin.NewFactory(algokit.AppFactoryParams{
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
	return &testProxyRekeyPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestTestProxyRekeyPluginDeployment(t *testing.T) {
	env := setupTestProxyRekeyPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("TestProxyRekeyPlugin deployed with appID: %d", env.client.AppID())
}

func TestTestProxyRekeyPluginComposer(t *testing.T) {
	env := setupTestProxyRekeyPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("TestProxyRekeyPlugin composer created successfully")
}

func TestTestProxyRekeyPluginFactory(t *testing.T) {
	env := setupTestProxyRekeyPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second TestProxyRekeyPlugin deployed with appID: %d", client2.AppID())
}

// --- ASAMintPlugin ---

type asamintPluginTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *asamintplugin.Client
	factory *asamintplugin.Factory
}

func setupAsamintPluginTest(t *testing.T) *asamintPluginTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := asamintplugin.NewFactory(algokit.AppFactoryParams{
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
	return &asamintPluginTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAsamintPluginDeployment(t *testing.T) {
	env := setupAsamintPluginTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("ASAMintPlugin deployed with appID: %d", env.client.AppID())
}

func TestAsamintPluginComposer(t *testing.T) {
	env := setupAsamintPluginTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("ASAMintPlugin composer created successfully")
}

func TestAsamintPluginFactory(t *testing.T) {
	env := setupAsamintPluginTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second ASAMintPlugin deployed with appID: %d", client2.AppID())
}

// --- AkitaDAOTypes ---

type akitaDaoTypesTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitadaotypes.Client
	factory *akitadaotypes.Factory
}

func setupAkitaDaoTypesTest(t *testing.T) *akitaDaoTypesTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitadaotypes.NewFactory(algokit.AppFactoryParams{
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
	return &akitaDaoTypesTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaDaoTypesDeployment(t *testing.T) {
	env := setupAkitaDaoTypesTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaDAOTypes deployed with appID: %d", env.client.AppID())
}

func TestAkitaDaoTypesComposer(t *testing.T) {
	env := setupAkitaDaoTypesTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaDAOTypes composer created successfully")
}

func TestAkitaDaoTypesFactory(t *testing.T) {
	env := setupAkitaDaoTypesTest(t)
	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second AkitaDAOTypes deployed with appID: %d", client2.AppID())
}

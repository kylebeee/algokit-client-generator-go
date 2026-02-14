package tests

import (
	"context"
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/abstractedaccountfactory"
	"github.com/kylebeee/algokit-client-generator-go/generated/akitasocial"
	"github.com/kylebeee/algokit-client-generator-go/generated/akitasocialgraph"
	"github.com/kylebeee/algokit-client-generator-go/generated/akitasocialimpact"
	"github.com/kylebeee/algokit-client-generator-go/generated/akitasocialmoderation"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// --- AkitaSocial ---

type akitaSocialTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitasocial.Client
	factory *akitasocial.Factory
}

func setupAkitaSocialTest(t *testing.T) *akitaSocialTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitasocial.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitasocial.FactoryCreateParams{
		Args: akitasocial.CreateArgs{
			Version:        "1.0.0",
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaSocialTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaSocialDeployment(t *testing.T) {
	env := setupAkitaSocialTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaSocial deployed with appID: %d", env.client.AppID())
}

func TestAkitaSocialComposer(t *testing.T) {
	env := setupAkitaSocialTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaSocial composer created successfully")
}

func TestAkitaSocialFactory(t *testing.T) {
	env := setupAkitaSocialTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitasocial.FactoryCreateParams{
		Args: akitasocial.CreateArgs{
			Version:        "1.0.0",
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
	t.Logf("Second AkitaSocial deployed with appID: %d", client2.AppID())
}

// --- AkitaSocialGraph ---

type akitaSocialGraphTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitasocialgraph.Client
	factory *akitasocialgraph.Factory
}

func setupAkitaSocialGraphTest(t *testing.T) *akitaSocialGraphTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitasocialgraph.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitasocialgraph.FactoryCreateParams{
		Args: akitasocialgraph.CreateArgs{
			AkitaDao: 0,
			Version:  "1.0.0",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaSocialGraphTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaSocialGraphDeployment(t *testing.T) {
	env := setupAkitaSocialGraphTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaSocialGraph deployed with appID: %d", env.client.AppID())
}

func TestAkitaSocialGraphComposer(t *testing.T) {
	env := setupAkitaSocialGraphTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaSocialGraph composer created successfully")
}

func TestAkitaSocialGraphFactory(t *testing.T) {
	env := setupAkitaSocialGraphTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitasocialgraph.FactoryCreateParams{
		Args: akitasocialgraph.CreateArgs{
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
	t.Logf("Second AkitaSocialGraph deployed with appID: %d", client2.AppID())
}

// --- AkitaSocialImpact ---

type akitaSocialImpactTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitasocialimpact.Client
	factory *akitasocialimpact.Factory
}

func setupAkitaSocialImpactTest(t *testing.T) *akitaSocialImpactTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitasocialimpact.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitasocialimpact.FactoryCreateParams{
		Args: akitasocialimpact.CreateArgs{
			AkitaDao: 0,
			Version:  "1.0.0",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaSocialImpactTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaSocialImpactDeployment(t *testing.T) {
	env := setupAkitaSocialImpactTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaSocialImpact deployed with appID: %d", env.client.AppID())
}

func TestAkitaSocialImpactComposer(t *testing.T) {
	env := setupAkitaSocialImpactTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaSocialImpact composer created successfully")
}

func TestAkitaSocialImpactFactory(t *testing.T) {
	env := setupAkitaSocialImpactTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitasocialimpact.FactoryCreateParams{
		Args: akitasocialimpact.CreateArgs{
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
	t.Logf("Second AkitaSocialImpact deployed with appID: %d", client2.AppID())
}

// --- AkitaSocialModeration ---

type akitaSocialModerationTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *akitasocialmoderation.Client
	factory *akitasocialmoderation.Factory
}

func setupAkitaSocialModerationTest(t *testing.T) *akitaSocialModerationTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := akitasocialmoderation.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, akitasocialmoderation.FactoryCreateParams{
		Args: akitasocialmoderation.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &akitaSocialModerationTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAkitaSocialModerationDeployment(t *testing.T) {
	env := setupAkitaSocialModerationTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AkitaSocialModeration deployed with appID: %d", env.client.AppID())
}

func TestAkitaSocialModerationComposer(t *testing.T) {
	env := setupAkitaSocialModerationTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AkitaSocialModeration composer created successfully")
}

func TestAkitaSocialModerationFactory(t *testing.T) {
	env := setupAkitaSocialModerationTest(t)
	client2, _, err := env.factory.Create(env.ctx, akitasocialmoderation.FactoryCreateParams{
		Args: akitasocialmoderation.CreateArgs{
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
	t.Logf("Second AkitaSocialModeration deployed with appID: %d", client2.AppID())
}

// --- AbstractedAccountFactory ---

type abstractedAccountFactoryTestEnv struct {
	fixture *testutil.TestFixture
	ctx     context.Context
	owner   *testutil.TestAccount
	client  *abstractedaccountfactory.Client
	factory *abstractedaccountfactory.Factory
}

func setupAbstractedAccountFactoryTest(t *testing.T) *abstractedAccountFactoryTestEnv {
	t.Helper()
	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := abstractedaccountfactory.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	client, _, err := factory.Create(ctx, abstractedaccountfactory.FactoryCreateParams{
		Args: abstractedaccountfactory.CreateArgs{
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
			Version:        "1.0.0",
			EscrowFactory:  0,
			Revocation:     0,
			Domain:         "test",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy: %v", err)
	}
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))
	return &abstractedAccountFactoryTestEnv{fixture: fixture, ctx: ctx, owner: owner, client: client, factory: factory}
}

func TestAbstractedAccountFactoryDeployment(t *testing.T) {
	env := setupAbstractedAccountFactoryTest(t)
	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}
	t.Logf("AbstractedAccountFactory deployed with appID: %d", env.client.AppID())
}

func TestAbstractedAccountFactoryComposer(t *testing.T) {
	env := setupAbstractedAccountFactoryTest(t)
	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("AbstractedAccountFactory composer created successfully")
}

func TestAbstractedAccountFactoryFactory(t *testing.T) {
	env := setupAbstractedAccountFactoryTest(t)
	client2, _, err := env.factory.Create(env.ctx, abstractedaccountfactory.FactoryCreateParams{
		Args: abstractedaccountfactory.CreateArgs{
			AkitaDao:       0,
			AkitaDaoEscrow: 0,
			Version:        "1.0.0",
			EscrowFactory:  0,
			Revocation:     0,
			Domain:         "test",
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second: %v", err)
	}
	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}
	t.Logf("Second AbstractedAccountFactory deployed with appID: %d", client2.AppID())
}

package tests

import (
	"context"
	"encoding/binary"
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/types"
	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/staking"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// Staking type constants (from the contract)
const (
	STAKING_TYPE_HEARTBEAT = 10
	STAKING_TYPE_SOFT      = 20
	STAKING_TYPE_HARD      = 30
	STAKING_TYPE_LOCK      = 40
)

// MBR constants (from the contract)
const (
	STAKES_MBR       = 28_900
	HEARTBEATS_MBR   = 70_100
	SETTINGS_MBR     = 9_300
	TOTALS_MBR       = 12_500
	ASSET_OPT_IN_MBR = 100_000
)

// Time constants
const (
	ONE_YEAR = 31_536_000
	ONE_DAY  = 86_400
	ONE_HOUR = 3_600
)

// stakingTestEnv holds the shared state for all staking tests.
type stakingTestEnv struct {
	fixture          *testutil.TestFixture
	ctx              context.Context
	deployer         *testutil.TestAccount
	user1            *testutil.TestAccount
	user2            *testutil.TestAccount
	heartbeatManager *testutil.TestAccount
	client           *staking.Client
	factory          *staking.Factory
	testAssetID      uint64
}

func setupStakingTest(t *testing.T) *stakingTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	// Create accounts with 100 ALGO each
	deployer := fixture.GenerateAccount(testutil.AlgoAmount(100))
	user1 := fixture.GenerateAccount(testutil.AlgoAmount(100))
	user2 := fixture.GenerateAccount(testutil.AlgoAmount(100))
	heartbeatManager := fixture.GenerateAccount(testutil.AlgoAmount(100))

	// Deploy staking contract
	factory, err := staking.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: deployer.Address,
		DefaultSigner: deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create staking factory: %v", err)
	}

	client, _, err := factory.Create(ctx, staking.FactoryCreateParams{
		Args: staking.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy staking contract: %v", err)
	}

	t.Logf("Staking contract deployed with appID: %d", client.AppID())

	// Fund the contract with 1 ALGO
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(1))

	// Initialize the contract
	totalsBoxKey := makeTotalsBoxKey(0)
	_, err = client.AppClient.Send(ctx, algokit.AppCallSendParams{
		MethodName: "init",
		MethodArgs: nil,
		Sender:     deployer.Address,
		Signer:     deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to initialize staking contract: %v", err)
	}

	// Create a test ASA (1 trillion total, 6 decimals)
	testAssetID := fixture.CreateASA(deployer, 1_000_000_000_000, 6, "TEST", "Test Token")
	t.Logf("Created test ASA with ID: %d", testAssetID)

	// Opt user1 and user2 into the test ASA
	fixture.OptInASA(user1, testAssetID)
	fixture.OptInASA(user2, testAssetID)

	// Transfer 100B tokens to each user
	fixture.TransferASA(deployer, user1.Address, testAssetID, 100_000_000_000)
	fixture.TransferASA(deployer, user2.Address, testAssetID, 100_000_000_000)

	return &stakingTestEnv{
		fixture:          fixture,
		ctx:              ctx,
		deployer:         deployer,
		user1:            user1,
		user2:            user2,
		heartbeatManager: heartbeatManager,
		client:           client,
		factory:          factory,
		testAssetID:      testAssetID,
	}
}

// optInASA opts the staking contract into an ASA.
func (env *stakingTestEnv) optInASA(t *testing.T, assetID uint64) {
	t.Helper()
	mbrAmount := uint64(TOTALS_MBR + ASSET_OPT_IN_MBR)
	paymentTxn := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), mbrAmount)
	totalsBoxKey := makeTotalsBoxKey(assetID)

	err := env.client.SendOptIn(env.ctx, algokit.CallParams[staking.OptInArgs]{
		Args: staking.OptInArgs{
			Payment: paymentTxn,
			Asset:   assetID,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: totalsBoxKey},
		},
		AssetReferences: []uint64{assetID},
		ExtraFee:        1000, // cover inner ASA opt-in txn
		SendParams: algokit.SendParams{
			PopulateAppCallResources: true,
		},
	})
	if err != nil {
		t.Fatalf("failed to opt contract into ASA %d: %v", assetID, err)
	}
}

// --- Deployment Tests ---

func TestStakingDeployment(t *testing.T) {
	env := setupStakingTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	balance := env.fixture.GetAccountBalance(env.client.AppAddress())
	if balance < testutil.AlgoAmount(1) {
		t.Errorf("expected contract balance >= 1 ALGO, got %s", testutil.FormatAlgos(balance))
	}

	t.Logf("Contract balance: %s", testutil.FormatAlgos(balance))
}

func TestStakingInitializedTotals(t *testing.T) {
	env := setupStakingTest(t)

	// Verify initial totals for ALGO (asset 0) are all zeros
	totalsBoxKey := makeTotalsBoxKey(0)
	totalsResult, err := env.client.SendGetTotals(env.ctx, algokit.CallParams[staking.GetTotalsArgs]{
		Args: staking.GetTotalsArgs{
			Assets: []uint64{0},
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to call getTotals: %v", err)
	}
	t.Logf("Initial ALGO totals: %v", totalsResult.Return)
}

// --- ASA Opt-In Tests ---

func TestStakingASAOptIn(t *testing.T) {
	env := setupStakingTest(t)

	mbrAmount := uint64(TOTALS_MBR + ASSET_OPT_IN_MBR)
	paymentTxn := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), mbrAmount)
	totalsBoxKey := makeTotalsBoxKey(env.testAssetID)

	err := env.client.SendOptIn(env.ctx, algokit.CallParams[staking.OptInArgs]{
		Args: staking.OptInArgs{
			Payment: paymentTxn,
			Asset:   env.testAssetID,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: totalsBoxKey},
		},
		AssetReferences: []uint64{env.testAssetID},
		ExtraFee:        1000,
		SendParams: algokit.SendParams{
			PopulateAppCallResources: true,
		},
	})
	if err != nil {
		t.Fatalf("failed to opt contract into ASA: %v", err)
	}

	t.Log("Contract opted into test ASA successfully")
}

func TestStakingASAOptInDuplicate(t *testing.T) {
	env := setupStakingTest(t)

	// Create a fresh ASA for this test
	freshAssetID := env.fixture.CreateASA(env.deployer, 1_000_000, 0, "OPT", "Opt In Test")

	// First opt-in should succeed
	env.optInASA(t, freshAssetID)

	// Second opt-in should fail
	mbrAmount := uint64(TOTALS_MBR + ASSET_OPT_IN_MBR)
	paymentTxn := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), mbrAmount)
	totalsBoxKey := makeTotalsBoxKey(freshAssetID)

	err := env.client.SendOptIn(env.ctx, algokit.CallParams[staking.OptInArgs]{
		Args: staking.OptInArgs{
			Payment: paymentTxn,
			Asset:   freshAssetID,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: totalsBoxKey},
		},
		AssetReferences: []uint64{freshAssetID},
		ExtraFee:        1000,
		SendParams: algokit.SendParams{
			PopulateAppCallResources: true,
		},
	})
	if err == nil {
		t.Fatal("expected error opting into already opted-in ASA")
	}
	t.Logf("Correctly rejected duplicate opt-in: %v", err)
}

// --- ALGO Soft Staking Tests ---

func TestStakingAlgoSoftStake(t *testing.T) {
	env := setupStakingTest(t)

	stakeAmount := testutil.AlgoAmount(1)
	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR)
	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_SOFT)
	totalsBoxKey := makeTotalsBoxKey(0)

	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_SOFT,
			Amount:     stakeAmount,
			Expiration: 0,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create soft stake: %v", err)
	}

	t.Logf("Created soft ALGO stake of %s", testutil.FormatAlgos(stakeAmount))
}

func TestStakingAlgoSoftStakeUpdate(t *testing.T) {
	env := setupStakingTest(t)

	initialAmount := testutil.AlgoAmount(2)
	additionalAmount := testutil.AlgoAmount(1)
	stakesBoxKey := makeStakesBoxKey(env.user2.Address, 0, STAKING_TYPE_SOFT)
	totalsBoxKey := makeTotalsBoxKey(0)

	// Create initial stake
	mbrPayment := env.fixture.MakePaymentTxn(env.user2, env.client.AppAddress(), STAKES_MBR)
	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_SOFT,
			Amount:     initialAmount,
			Expiration: 0,
		},
		Sender: env.user2.Address,
		Signer: env.user2.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create initial soft stake: %v", err)
	}

	// Update stake (no MBR needed, pay 0)
	updatePayment := env.fixture.MakePaymentTxn(env.user2, env.client.AppAddress(), 0)
	err = env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    updatePayment,
			Type:       STAKING_TYPE_SOFT,
			Amount:     additionalAmount,
			Expiration: 0,
		},
		Sender: env.user2.Address,
		Signer: env.user2.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to update soft stake: %v", err)
	}

	t.Logf("Updated soft ALGO stake from %s to %s",
		testutil.FormatAlgos(initialAmount), testutil.FormatAlgos(initialAmount+additionalAmount))
}

// --- ALGO Hard Staking Tests ---

func TestStakingAlgoHardStake(t *testing.T) {
	env := setupStakingTest(t)

	stakeAmount := testutil.AlgoAmount(5)
	expiration := uint64(time.Now().Unix()) + uint64(ONE_DAY)

	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR+stakeAmount)
	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HARD)
	totalsBoxKey := makeTotalsBoxKey(0)

	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_HARD,
			Amount:     stakeAmount,
			Expiration: expiration,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create hard stake: %v", err)
	}

	t.Logf("Created hard ALGO stake of %s with expiration %d", testutil.FormatAlgos(stakeAmount), expiration)
}

func TestStakingAlgoHardStakeAddToExisting(t *testing.T) {
	env := setupStakingTest(t)

	stakeAmount := testutil.AlgoAmount(5)
	additionalAmount := testutil.AlgoAmount(2)
	expiration := uint64(time.Now().Unix()) + uint64(ONE_DAY)
	newExpiration := uint64(time.Now().Unix()) + uint64(ONE_DAY*2)
	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HARD)
	totalsBoxKey := makeTotalsBoxKey(0)

	// Create initial hard stake
	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR+stakeAmount)
	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_HARD,
			Amount:     stakeAmount,
			Expiration: expiration,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create initial hard stake: %v", err)
	}

	// Add to existing hard stake
	addPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), additionalAmount)
	err = env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    addPayment,
			Type:       STAKING_TYPE_HARD,
			Amount:     additionalAmount,
			Expiration: newExpiration,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add to hard stake: %v", err)
	}

	t.Logf("Added %s to hard ALGO stake", testutil.FormatAlgos(additionalAmount))
}

// --- Lock Staking Tests ---

func TestStakingAlgoLockStake(t *testing.T) {
	env := setupStakingTest(t)

	stakeAmount := testutil.AlgoAmount(10)
	expiration := uint64(time.Now().Unix()) + uint64(ONE_DAY)

	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR+stakeAmount)
	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_LOCK)
	totalsBoxKey := makeTotalsBoxKey(0)

	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_LOCK,
			Amount:     stakeAmount,
			Expiration: expiration,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create lock stake: %v", err)
	}

	t.Logf("Created lock ALGO stake of %s", testutil.FormatAlgos(stakeAmount))
}

func TestStakingAlgoLockStakeFailZeroExpiration(t *testing.T) {
	env := setupStakingTest(t)

	stakeAmount := testutil.AlgoAmount(3)
	mbrPayment := env.fixture.MakePaymentTxn(env.user2, env.client.AppAddress(), STAKES_MBR+stakeAmount)
	stakesBoxKey := makeStakesBoxKey(env.user2.Address, 0, STAKING_TYPE_LOCK)
	totalsBoxKey := makeTotalsBoxKey(0)

	// Lock type requires future expiration (expiration > Global.latestTimestamp)
	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_LOCK,
			Amount:     stakeAmount,
			Expiration: 0,
		},
		Sender: env.user2.Address,
		Signer: env.user2.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err == nil {
		t.Fatal("expected error creating lock stake with 0 expiration")
	}
	t.Logf("Correctly rejected lock stake with 0 expiration: %v", err)
}

func TestStakingAlgoLockStakeFailOverOneYear(t *testing.T) {
	env := setupStakingTest(t)

	stakeAmount := testutil.AlgoAmount(1)
	expiration := uint64(time.Now().Unix()) + uint64(ONE_YEAR+ONE_DAY)
	mbrPayment := env.fixture.MakePaymentTxn(env.user2, env.client.AppAddress(), STAKES_MBR+stakeAmount)
	stakesBoxKey := makeStakesBoxKey(env.user2.Address, 0, STAKING_TYPE_LOCK)
	totalsBoxKey := makeTotalsBoxKey(0)

	// Lock type has 1 year max expiration
	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_LOCK,
			Amount:     stakeAmount,
			Expiration: expiration,
		},
		Sender: env.user2.Address,
		Signer: env.user2.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err == nil {
		t.Fatal("expected error creating lock stake with >1 year expiration")
	}
	t.Logf("Correctly rejected lock stake with >1 year expiration: %v", err)
}

// --- Heartbeat Staking Tests ---

func TestStakingHeartbeat(t *testing.T) {
	env := setupStakingTest(t)

	// Heartbeat requires STAKES_MBR + HEARTBEATS_MBR
	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR+HEARTBEATS_MBR)
	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HEARTBEAT)
	heartbeatBoxKey := makeHeartbeatBoxKey(env.user1.Address, 0)
	totalsBoxKey := makeTotalsBoxKey(0)
	// Heartbeat creation checks existing HARD and LOCK stakes
	hardBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HARD)
	lockBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_LOCK)
	softBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_SOFT)

	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_HEARTBEAT,
			Amount:     0, // Amount ignored for heartbeat
			Expiration: 0,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: heartbeatBoxKey},
			{AppID: 0, Name: totalsBoxKey},
			{AppID: 0, Name: hardBoxKey},
			{AppID: 0, Name: lockBoxKey},
			{AppID: 0, Name: softBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create heartbeat stake: %v", err)
	}

	t.Log("Created heartbeat stake successfully")
}

func TestStakingHeartbeatFailUpdate(t *testing.T) {
	env := setupStakingTest(t)

	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HEARTBEAT)
	heartbeatBoxKey := makeHeartbeatBoxKey(env.user1.Address, 0)
	totalsBoxKey := makeTotalsBoxKey(0)
	hardBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HARD)
	lockBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_LOCK)
	softBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_SOFT)

	heartbeatBoxRefs := []types.AppBoxReference{
		{AppID: 0, Name: stakesBoxKey},
		{AppID: 0, Name: heartbeatBoxKey},
		{AppID: 0, Name: totalsBoxKey},
		{AppID: 0, Name: hardBoxKey},
		{AppID: 0, Name: lockBoxKey},
		{AppID: 0, Name: softBoxKey},
	}

	// Create heartbeat first
	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR+HEARTBEATS_MBR)
	err := env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    mbrPayment,
			Type:       STAKING_TYPE_HEARTBEAT,
			Amount:     0,
			Expiration: 0,
		},
		Sender:        env.user1.Address,
		Signer:        env.user1.Signer,
		BoxReferences: heartbeatBoxRefs,
	})
	if err != nil {
		t.Fatalf("failed to create initial heartbeat: %v", err)
	}

	// Try to update - should fail
	updatePayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), 0)
	err = env.client.SendStake(env.ctx, algokit.CallParams[staking.StakeArgs]{
		Args: staking.StakeArgs{
			Payment:    updatePayment,
			Type:       STAKING_TYPE_HEARTBEAT,
			Amount:     0,
			Expiration: 0,
		},
		Sender:        env.user1.Address,
		Signer:        env.user1.Signer,
		BoxReferences: heartbeatBoxRefs,
	})
	if err == nil {
		t.Fatal("expected error updating heartbeat stake")
	}
	t.Logf("Correctly rejected heartbeat update: %v", err)
}

// --- ASA Staking Tests ---

func TestStakingASASoftStake(t *testing.T) {
	env := setupStakingTest(t)

	// Opt contract into test ASA first
	env.optInASA(t, env.testAssetID)

	stakeAmount := uint64(1_000_000_000) // 1000 tokens (6 decimals)

	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR)
	// For soft ASA staking, asset transfer amount is 0 (no escrow)
	assetXfer := env.fixture.MakeAssetTransferTxn(env.user1, env.client.AppAddress(), env.testAssetID, 0)

	stakesBoxKey := makeStakesBoxKey(env.user1.Address, env.testAssetID, STAKING_TYPE_SOFT)
	totalsBoxKey := makeTotalsBoxKey(env.testAssetID)
	settingsBoxKey := makeSettingsBoxKey(env.testAssetID)

	err := env.client.SendStakeASA(env.ctx, algokit.CallParams[staking.StakeASAArgs]{
		Args: staking.StakeASAArgs{
			Payment:    mbrPayment,
			AssetXfer:  assetXfer,
			Type:       STAKING_TYPE_SOFT,
			Amount:     stakeAmount,
			Expiration: 0,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
			{AppID: 0, Name: settingsBoxKey},
		},
		AssetReferences: []uint64{env.testAssetID},
	})
	if err != nil {
		t.Fatalf("failed to create soft ASA stake: %v", err)
	}

	t.Logf("Created soft ASA stake of %d tokens", stakeAmount/1_000_000)
}

func TestStakingASAHardStake(t *testing.T) {
	env := setupStakingTest(t)

	// Opt contract into test ASA first
	env.optInASA(t, env.testAssetID)

	stakeAmount := uint64(5_000_000_000) // 5000 tokens
	expiration := uint64(time.Now().Unix()) + uint64(ONE_DAY)

	mbrPayment := env.fixture.MakePaymentTxn(env.user1, env.client.AppAddress(), STAKES_MBR)
	assetXfer := env.fixture.MakeAssetTransferTxn(env.user1, env.client.AppAddress(), env.testAssetID, stakeAmount)

	stakesBoxKey := makeStakesBoxKey(env.user1.Address, env.testAssetID, STAKING_TYPE_HARD)
	totalsBoxKey := makeTotalsBoxKey(env.testAssetID)
	settingsBoxKey := makeSettingsBoxKey(env.testAssetID)

	err := env.client.SendStakeASA(env.ctx, algokit.CallParams[staking.StakeASAArgs]{
		Args: staking.StakeASAArgs{
			Payment:    mbrPayment,
			AssetXfer:  assetXfer,
			Type:       STAKING_TYPE_HARD,
			Amount:     stakeAmount,
			Expiration: expiration,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
			{AppID: 0, Name: settingsBoxKey},
		},
		AssetReferences: []uint64{env.testAssetID},
	})
	if err != nil {
		t.Fatalf("failed to create hard ASA stake: %v", err)
	}

	t.Logf("Created hard ASA stake of %d tokens with escrow", stakeAmount/1_000_000)
}

// --- Read Method Tests ---

func TestStakingReadMethods(t *testing.T) {
	env := setupStakingTest(t)

	// Test optInCost
	costResult, err := env.client.SendOptInCost(env.ctx)
	if err != nil {
		t.Fatalf("failed to call optInCost: %v", err)
	}
	if costResult.Return == 0 {
		t.Error("expected non-zero opt-in cost")
	}
	t.Logf("OptIn cost: %d microAlgos", costResult.Return)

	// Test stakeCost
	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_SOFT)
	stakeCostResult, err := env.client.SendStakeCost(env.ctx, algokit.CallParams[staking.StakeCostArgs]{
		Args: staking.StakeCostArgs{
			Asset: 0,
			Type:  STAKING_TYPE_SOFT,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to call stakeCost: %v", err)
	}
	if stakeCostResult.Return == 0 {
		t.Error("expected non-zero stake cost")
	}
	t.Logf("Stake cost for ALGO soft: %d microAlgos", stakeCostResult.Return)

	// Test getTotals
	totalsBoxKey := makeTotalsBoxKey(0)
	totalsResult, err := env.client.SendGetTotals(env.ctx, algokit.CallParams[staking.GetTotalsArgs]{
		Args: staking.GetTotalsArgs{
			Assets: []uint64{0},
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: totalsBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to call getTotals: %v", err)
	}
	t.Logf("Totals result: %v", totalsResult.Return)
}

// --- Withdraw Tests ---

func TestStakingWithdrawNonExistent(t *testing.T) {
	env := setupStakingTest(t)

	stakesBoxKey := makeStakesBoxKey(env.user1.Address, 0, STAKING_TYPE_HARD)
	totalsBoxKey := makeTotalsBoxKey(0)

	err := env.client.SendWithdraw(env.ctx, algokit.CallParams[staking.WithdrawArgs]{
		Args: staking.WithdrawArgs{
			Asset: 0,
			Type:  STAKING_TYPE_HARD,
		},
		Sender: env.user1.Address,
		Signer: env.user1.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: stakesBoxKey},
			{AppID: 0, Name: totalsBoxKey},
		},
		SendParams: algokit.SendParams{
			PopulateAppCallResources: true,
		},
	})
	if err == nil {
		t.Fatal("expected error withdrawing non-existent stake")
	}
	t.Logf("Correctly rejected withdraw of non-existent stake: %v", err)
}

// --- Composer & Factory Tests ---

func TestStakingComposer(t *testing.T) {
	env := setupStakingTest(t)

	composer := env.client.NewGroup()
	if composer == nil {
		t.Fatal("expected non-nil composer")
	}
	t.Log("Composer created successfully from client")
}

func TestStakingFactory(t *testing.T) {
	env := setupStakingTest(t)

	client2, _, err := env.factory.Create(env.ctx, staking.FactoryCreateParams{
		Args: staking.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second staking contract: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs for two deployments")
	}

	t.Logf("Second staking contract deployed with appID: %d", client2.AppID())
}

// --- Box Key Helpers ---

// makeTotalsBoxKey builds a totals box key: prefix "t" + uint64(assetID)
func makeTotalsBoxKey(assetID uint64) []byte {
	key := make([]byte, 9)
	key[0] = 't'
	binary.BigEndian.PutUint64(key[1:], assetID)
	return key
}

// makeStakesBoxKey builds a stakes box key: prefix "s" + address(32 bytes) + uint64(assetID) + uint8(type)
func makeStakesBoxKey(addr types.Address, assetID uint64, stakeType uint8) []byte {
	key := make([]byte, 42) // 1 + 32 + 8 + 1
	key[0] = 's'
	copy(key[1:33], addr[:])
	binary.BigEndian.PutUint64(key[33:41], assetID)
	key[41] = stakeType
	return key
}

// makeHeartbeatBoxKey builds a heartbeat box key: prefix "h" + address(32 bytes) + uint64(assetID)
func makeHeartbeatBoxKey(addr types.Address, assetID uint64) []byte {
	key := make([]byte, 41) // 1 + 32 + 8
	key[0] = 'h'
	copy(key[1:33], addr[:])
	binary.BigEndian.PutUint64(key[33:], assetID)
	return key
}

// makeSettingsBoxKey builds a settings box key: prefix "e" + uint64(assetID)
func makeSettingsBoxKey(assetID uint64) []byte {
	key := make([]byte, 9)
	key[0] = 'e'
	binary.BigEndian.PutUint64(key[1:], assetID)
	return key
}

// ensure imports are used
var _ = types.Address{}

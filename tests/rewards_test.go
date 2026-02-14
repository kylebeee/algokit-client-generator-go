package tests

import (
	"context"
	"encoding/binary"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/rewards"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// rewardsTestEnv holds the shared state for Rewards tests.
type rewardsTestEnv struct {
	fixture  *testutil.TestFixture
	ctx      context.Context
	deployer *testutil.TestAccount
	user1    *testutil.TestAccount
	user2    *testutil.TestAccount
	client   *rewards.Client
	factory  *rewards.Factory
}

func setupRewardsTest(t *testing.T) *rewardsTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	deployer := fixture.GenerateAccount(testutil.AlgoAmount(100))
	user1 := fixture.GenerateAccount(testutil.AlgoAmount(10))
	user2 := fixture.GenerateAccount(testutil.AlgoAmount(10))

	factory, err := rewards.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: deployer.Address,
		DefaultSigner: deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create rewards factory: %v", err)
	}

	// Deploy with create(version, akitaDAO) - pass 0 as dummy akitaDAO
	client, _, err := factory.Create(ctx, rewards.FactoryCreateParams{
		Args: rewards.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy rewards contract: %v", err)
	}

	t.Logf("Rewards contract deployed with appID: %d", client.AppID())

	// Fund the contract
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(10))

	return &rewardsTestEnv{
		fixture:  fixture,
		ctx:      ctx,
		deployer: deployer,
		user1:    user1,
		user2:    user2,
		client:   client,
		factory:  factory,
	}
}

// makeDisbursementBoxKey builds box key: prefix 'd' + uint64(disbursementID)
func makeDisbursementBoxKey(id uint64) []byte {
	key := make([]byte, 9)
	key[0] = 'd'
	binary.BigEndian.PutUint64(key[1:], id)
	return key
}

// makeUserAllocationBoxKey builds box key: prefix 'u' + ABI-encoded UserAllocationsKey
// UserAllocationsKey = (address, uint64, uint64)
func makeUserAllocationBoxKey(address types.Address, asset uint64, disbursementID uint64) []byte {
	// ABI tuple: address(32) + uint64(8) + uint64(8) = 48 bytes, all static
	key := make([]byte, 1+32+8+8)
	key[0] = 'u'
	copy(key[1:33], address[:])
	binary.BigEndian.PutUint64(key[33:41], asset)
	binary.BigEndian.PutUint64(key[41:49], disbursementID)
	return key
}

func TestRewardsDeployment(t *testing.T) {
	env := setupRewardsTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("Rewards deployed with appID: %d", env.client.AppID())
}

func TestRewardsMBR(t *testing.T) {
	env := setupRewardsTest(t)

	result, err := env.client.SendMBR(env.ctx, algokit.CallParams[rewards.MBRArgs]{
		Args: rewards.MBRArgs{
			Title: "Test Reward",
			Note:  "A test disbursement",
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to call mbr: %v", err)
	}

	t.Logf("MBR - Disbursements: %d, UserAllocations: %d",
		result.Return.Disbursements, result.Return.UserAllocations)

	if result.Return.Disbursements == 0 {
		t.Error("expected non-zero disbursement MBR")
	}
	if result.Return.UserAllocations == 0 {
		t.Error("expected non-zero user allocation MBR")
	}
}

func TestRewardsCreateDisbursement(t *testing.T) {
	env := setupRewardsTest(t)

	title := "Test Reward"
	note := "Testing"

	// Get MBR cost
	mbrResult, err := env.client.SendMBR(env.ctx, algokit.CallParams[rewards.MBRArgs]{
		Args: rewards.MBRArgs{
			Title: title,
			Note:  note,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get MBR: %v", err)
	}

	disbursementBoxKey := makeDisbursementBoxKey(1) // first disbursement = ID 1 (counter starts at 0, incremented before use)

	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), mbrResult.Return.Disbursements)

	result, err := env.client.SendCreateDisbursement(env.ctx, algokit.CallParams[rewards.CreateDisbursementArgs]{
		Args: rewards.CreateDisbursementArgs{
			MBRPayment:   mbrPayment,
			Title:        title,
			TimeToUnlock: 0,
			Expiration:   0,
			Note:         note,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: disbursementBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create disbursement: %v", err)
	}

	t.Logf("Created disbursement with ID: %d", result.Return)
}

func TestRewardsEditDisbursement(t *testing.T) {
	env := setupRewardsTest(t)

	title := "Initial"
	note := ""

	// Get MBR and create disbursement first
	mbrResult, err := env.client.SendMBR(env.ctx, algokit.CallParams[rewards.MBRArgs]{
		Args:   rewards.MBRArgs{Title: title, Note: note},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get MBR: %v", err)
	}

	disbursementBoxKey := makeDisbursementBoxKey(1)
	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), mbrResult.Return.Disbursements)

	createResult, err := env.client.SendCreateDisbursement(env.ctx, algokit.CallParams[rewards.CreateDisbursementArgs]{
		Args: rewards.CreateDisbursementArgs{
			MBRPayment:   mbrPayment,
			Title:        title,
			TimeToUnlock: 0,
			Expiration:   0,
			Note:         note,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: disbursementBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create disbursement: %v", err)
	}

	// Edit the disbursement
	err = env.client.SendEditDisbursement(env.ctx, algokit.CallParams[rewards.EditDisbursementArgs]{
		Args: rewards.EditDisbursementArgs{
			ID:           createResult.Return,
			Title:        "Updated",
			TimeToUnlock: 100,
			Expiration:   1000,
			Note:         "updated note",
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: disbursementBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to edit disbursement: %v", err)
	}

	t.Logf("Edited disbursement %d", createResult.Return)
}

func TestRewardsCreateUserAllocations(t *testing.T) {
	env := setupRewardsTest(t)

	title := "Airdrop"
	note := ""

	// Get MBR and create disbursement
	mbrResult, err := env.client.SendMBR(env.ctx, algokit.CallParams[rewards.MBRArgs]{
		Args:   rewards.MBRArgs{Title: title, Note: note},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get MBR: %v", err)
	}

	disbursementBoxKey := makeDisbursementBoxKey(1)
	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), mbrResult.Return.Disbursements)

	createResult, err := env.client.SendCreateDisbursement(env.ctx, algokit.CallParams[rewards.CreateDisbursementArgs]{
		Args: rewards.CreateDisbursementArgs{
			MBRPayment:   mbrPayment,
			Title:        title,
			TimeToUnlock: 0,
			Expiration:   0,
			Note:         note,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: disbursementBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create disbursement: %v", err)
	}

	// Create user allocations - allocate ALGO to user1 and user2
	// Allocations format: [(address, uint64_amount), ...]
	user1AllocBoxKey := makeUserAllocationBoxKey(env.user1.Address, 0, createResult.Return)
	user2AllocBoxKey := makeUserAllocationBoxKey(env.user2.Address, 0, createResult.Return)

	// Payment = sum of ALGO allocations + MBR per allocation
	algoAmount := uint64(1_000_000 + 500_000) // 1.5 ALGO total allocation
	mbrCost := mbrResult.Return.UserAllocations * 2  // 2 allocations
	allocPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), algoAmount+mbrCost)

	err = env.client.SendCreateUserAllocations(env.ctx, algokit.CallParams[rewards.CreateUserAllocationsArgs]{
		Args: rewards.CreateUserAllocationsArgs{
			Payment: allocPayment,
			ID:      createResult.Return,
			Allocations: [][]interface{}{
				{env.user1.Address, uint64(1_000_000)},
				{env.user2.Address, uint64(500_000)},
			},
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: disbursementBoxKey},
			{AppID: 0, Name: user1AllocBoxKey},
			{AppID: 0, Name: user2AllocBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to create user allocations: %v", err)
	}

	t.Logf("Created user allocations for disbursement %d", createResult.Return)
}

func TestRewardsFactory(t *testing.T) {
	env := setupRewardsTest(t)

	client2, _, err := env.factory.Create(env.ctx, rewards.FactoryCreateParams{
		Args: rewards.CreateArgs{
			Version:  "1.0.0",
			AkitaDao: 0,
		},
		ExtraPages: 3,
	})
	if err != nil {
		t.Fatalf("failed to deploy second rewards: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second Rewards deployed with appID: %d", client2.AppID())
}

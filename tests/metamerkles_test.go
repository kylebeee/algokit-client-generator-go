package tests

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/metamerkles"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// metamerklesTestEnv holds the shared state for MetaMerkles tests.
type metamerklesTestEnv struct {
	fixture  *testutil.TestFixture
	ctx      context.Context
	deployer *testutil.TestAccount
	user1    *testutil.TestAccount
	client   *metamerkles.Client
	factory  *metamerkles.Factory
}

func setupMetaMerklesTest(t *testing.T) *metamerklesTestEnv {
	t.Helper()

	fixture := testutil.NewTestFixture(t)
	ctx := fixture.Ctx

	// Need extra ALGO for addType (100 ALGO payment required)
	deployer := fixture.GenerateAccount(testutil.AlgoAmount(300))
	user1 := fixture.GenerateAccount(testutil.AlgoAmount(10))

	factory, err := metamerkles.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: deployer.Address,
		DefaultSigner: deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create metamerkles factory: %v", err)
	}

	// Deploy with create() method call (no args - MethodName set internally by factory)
	client, _, err := factory.Create(ctx, algokit.AppFactoryCreateParams{
		ExtraPages: 2,
	})
	if err != nil {
		t.Fatalf("failed to deploy metamerkles contract: %v", err)
	}

	t.Logf("MetaMerkles contract deployed with appID: %d", client.AppID())

	// Fund the contract generously for inner txn MBR returns
	fixture.FundApp(client.AppID(), testutil.AlgoAmount(5))

	return &metamerklesTestEnv{
		fixture:  fixture,
		ctx:      ctx,
		deployer: deployer,
		user1:    user1,
		client:   client,
		factory:  factory,
	}
}

// makeTypeBoxKey builds the box key for the types box map: prefix 't' + uint64(typeID)
func makeTypeBoxKey(typeID uint64) []byte {
	key := make([]byte, 9)
	key[0] = 't'
	binary.BigEndian.PutUint64(key[1:], typeID)
	return key
}

// makeRootBoxKey builds the ABI-encoded box key for the roots box map.
// RootKey = (address, string) → prefix 'r' + ABI tuple encoding
func makeRootBoxKey(address types.Address, name string) []byte {
	// ABI tuple: address(32 fixed) + uint16(offset=34) + uint16(len) + name bytes
	total := 1 + 32 + 2 + 2 + len(name)
	key := make([]byte, total)
	key[0] = 'r'
	copy(key[1:33], address[:])
	binary.BigEndian.PutUint16(key[33:35], 34)              // offset to string data
	binary.BigEndian.PutUint16(key[35:37], uint16(len(name))) // string length
	copy(key[37:], []byte(name))
	return key
}

// makeDataBoxKey builds the ABI-encoded box key for the data box map.
// DataKey = (bytes<16>, string, string) → prefix 'd' + ABI tuple encoding
func makeDataBoxKey(address types.Address, name string, dataKey string) []byte {
	// ABI tuple: bytes16(16 fixed) + uint16(offset1) + uint16(offset2)
	//   + uint16(nameLen) + name + uint16(keyLen) + key
	staticSize := 16 + 2 + 2 // = 20
	offset1 := uint16(staticSize)
	offset2 := uint16(staticSize + 2 + len(name))
	total := 1 + staticSize + 2 + len(name) + 2 + len(dataKey)
	key := make([]byte, total)
	key[0] = 'd'
	copy(key[1:17], address[:16])                                    // truncated address
	binary.BigEndian.PutUint16(key[17:19], offset1)                  // offset to first string
	binary.BigEndian.PutUint16(key[19:21], offset2)                  // offset to second string
	binary.BigEndian.PutUint16(key[21:23], uint16(len(name)))        // first string length
	copy(key[23:23+len(name)], []byte(name))                         // first string data
	pos := 23 + len(name)
	binary.BigEndian.PutUint16(key[pos:pos+2], uint16(len(dataKey))) // second string length
	copy(key[pos+2:], []byte(dataKey))                               // second string data
	return key
}

// addTypeToContract registers a new type and returns the type ID (auto-incremented from 0)
func addTypeToContract(t *testing.T, env *metamerklesTestEnv, typeID uint64) {
	t.Helper()

	typeBoxKey := makeTypeBoxKey(typeID)

	// addType requires exactly 100 ALGO payment
	typePayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), 100_000_000)

	err := env.client.SendAddType(env.ctx, algokit.CallParams[metamerkles.AddTypeArgs]{
		Args: metamerkles.AddTypeArgs{
			Payment:     typePayment,
			Description: "test_type",
			SchemaList:  []uint8{13}, // SchemaPartUint64 = 13
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: typeBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add type %d: %v", typeID, err)
	}
	t.Logf("Registered type %d", typeID)
}

func TestMetaMerklesDeployment(t *testing.T) {
	env := setupMetaMerklesTest(t)

	if env.client.AppID() == 0 {
		t.Fatal("expected non-zero app ID")
	}

	t.Logf("MetaMerkles deployed with appID: %d", env.client.AppID())
}

func TestMetaMerklesRootCosts(t *testing.T) {
	env := setupMetaMerklesTest(t)

	// Check cost to add a root
	result, err := env.client.SendRootCosts(env.ctx, algokit.CallParams[metamerkles.RootCostsArgs]{
		Args: metamerkles.RootCostsArgs{
			Name: "test_root",
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to call rootCosts: %v", err)
	}

	if result.Return == 0 {
		t.Error("expected non-zero root costs")
	}
	t.Logf("Root costs for 'test_root': %d microAlgos", result.Return)
}

func TestMetaMerklesAddType(t *testing.T) {
	env := setupMetaMerklesTest(t)

	// Register type 0 (auto-incremented)
	addTypeToContract(t, env, 0)

	t.Log("Successfully registered type 0")
}

func TestMetaMerklesAddRoot(t *testing.T) {
	env := setupMetaMerklesTest(t)

	// Must register a type first - addRoot checks this.types(type).exists
	addTypeToContract(t, env, 0)

	// Create a simple merkle root (sha256 of "hello")
	root := sha256.Sum256([]byte("hello"))
	rootName := "test_root"

	// Get the cost first
	costResult, err := env.client.SendRootCosts(env.ctx, algokit.CallParams[metamerkles.RootCostsArgs]{
		Args: metamerkles.RootCostsArgs{Name: rootName},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get root costs: %v", err)
	}

	// Box keys
	rootBoxKey := makeRootBoxKey(env.deployer.Address, rootName)
	typeBoxKey := makeTypeBoxKey(0)
	dataBoxKey := makeDataBoxKey(env.deployer.Address, rootName, "l.type")

	// MBR payment for root creation
	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), costResult.Return)

	err = env.client.SendAddRoot(env.ctx, algokit.CallParams[metamerkles.AddRootArgs]{
		Args: metamerkles.AddRootArgs{
			Payment: mbrPayment,
			Name:    rootName,
			Root:    root,
			Type:    0,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: typeBoxKey},
			{AppID: 0, Name: rootBoxKey},
			{AppID: 0, Name: dataBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	t.Logf("Added root '%s' with hash %x", rootName, root)
}

func TestMetaMerklesUpdateRoot(t *testing.T) {
	env := setupMetaMerklesTest(t)

	// Register type first
	addTypeToContract(t, env, 0)

	rootName := "updatable"
	root1 := sha256.Sum256([]byte("version1"))
	root2 := sha256.Sum256([]byte("version2"))

	// Get cost and add initial root
	costResult, err := env.client.SendRootCosts(env.ctx, algokit.CallParams[metamerkles.RootCostsArgs]{
		Args: metamerkles.RootCostsArgs{Name: rootName},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get root costs: %v", err)
	}

	rootBoxKey := makeRootBoxKey(env.deployer.Address, rootName)
	typeBoxKey := makeTypeBoxKey(0)
	dataBoxKey := makeDataBoxKey(env.deployer.Address, rootName, "l.type")

	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), costResult.Return)

	err = env.client.SendAddRoot(env.ctx, algokit.CallParams[metamerkles.AddRootArgs]{
		Args: metamerkles.AddRootArgs{
			Payment: mbrPayment,
			Name:    rootName,
			Root:    root1,
			Type:    0,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: typeBoxKey},
			{AppID: 0, Name: rootBoxKey},
			{AppID: 0, Name: dataBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add initial root: %v", err)
	}

	// Update the root - only accesses root box
	err = env.client.SendUpdateRoot(env.ctx, algokit.CallParams[metamerkles.UpdateRootArgs]{
		Args: metamerkles.UpdateRootArgs{
			Name:    rootName,
			NewRoot: root2,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: rootBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to update root: %v", err)
	}

	t.Logf("Updated root '%s' from %x to %x", rootName, root1[:4], root2[:4])
}

func TestMetaMerklesDeleteRoot(t *testing.T) {
	env := setupMetaMerklesTest(t)

	// Register type first
	addTypeToContract(t, env, 0)

	rootName := "deletable"
	root := sha256.Sum256([]byte("temp"))

	// Add root
	costResult, err := env.client.SendRootCosts(env.ctx, algokit.CallParams[metamerkles.RootCostsArgs]{
		Args: metamerkles.RootCostsArgs{Name: rootName},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get root costs: %v", err)
	}

	rootBoxKey := makeRootBoxKey(env.deployer.Address, rootName)
	typeBoxKey := makeTypeBoxKey(0)
	dataBoxKey := makeDataBoxKey(env.deployer.Address, rootName, "l.type")

	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), costResult.Return)

	err = env.client.SendAddRoot(env.ctx, algokit.CallParams[metamerkles.AddRootArgs]{
		Args: metamerkles.AddRootArgs{
			Payment: mbrPayment,
			Name:    rootName,
			Root:    root,
			Type:    0,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: typeBoxKey},
			{AppID: 0, Name: rootBoxKey},
			{AppID: 0, Name: dataBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Delete the root - accesses root + data boxes, issues inner payment (needs ExtraFee)
	err = env.client.SendDeleteRoot(env.ctx, algokit.CallParams[metamerkles.DeleteRootArgs]{
		Args: metamerkles.DeleteRootArgs{
			Name: rootName,
		},
		Sender:    env.deployer.Address,
		Signer:    env.deployer.Signer,
		StaticFee: 10000, // inner payment txn fee pooling
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: rootBoxKey},
			{AppID: 0, Name: dataBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to delete root: %v", err)
	}

	t.Logf("Successfully deleted root '%s'", rootName)
}

func TestMetaMerklesDataCosts(t *testing.T) {
	env := setupMetaMerklesTest(t)

	result, err := env.client.SendDataCosts(env.ctx, algokit.CallParams[metamerkles.DataCostsArgs]{
		Args: metamerkles.DataCostsArgs{
			Name:  "data_root",
			Key:   "key1",
			Value: "value1",
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to call dataCosts: %v", err)
	}

	if result.Return == 0 {
		t.Error("expected non-zero data costs")
	}
	t.Logf("Data costs for key 'key1': %d microAlgos", result.Return)
}

func TestMetaMerklesAddData(t *testing.T) {
	env := setupMetaMerklesTest(t)

	// Register type and add a root first
	addTypeToContract(t, env, 0)

	rootName := "data_test"
	root := sha256.Sum256([]byte("data_root"))

	// Add root
	costResult, err := env.client.SendRootCosts(env.ctx, algokit.CallParams[metamerkles.RootCostsArgs]{
		Args: metamerkles.RootCostsArgs{Name: rootName},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get root costs: %v", err)
	}

	rootBoxKey := makeRootBoxKey(env.deployer.Address, rootName)
	typeBoxKey := makeTypeBoxKey(0)
	typeDataBoxKey := makeDataBoxKey(env.deployer.Address, rootName, "l.type")

	mbrPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), costResult.Return)
	err = env.client.SendAddRoot(env.ctx, algokit.CallParams[metamerkles.AddRootArgs]{
		Args: metamerkles.AddRootArgs{
			Payment: mbrPayment,
			Name:    rootName,
			Root:    root,
			Type:    0,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: typeBoxKey},
			{AppID: 0, Name: rootBoxKey},
			{AppID: 0, Name: typeDataBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Now add data to the root
	dataKey := "royalty"
	dataValue := "5"

	dataCostResult, err := env.client.SendDataCosts(env.ctx, algokit.CallParams[metamerkles.DataCostsArgs]{
		Args: metamerkles.DataCostsArgs{
			Name:  rootName,
			Key:   dataKey,
			Value: dataValue,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
	})
	if err != nil {
		t.Fatalf("failed to get data costs: %v", err)
	}

	dataBoxKey := makeDataBoxKey(env.deployer.Address, rootName, dataKey)
	dataPayment := env.fixture.MakePaymentTxn(env.deployer, env.client.AppAddress(), dataCostResult.Return)

	err = env.client.SendAddData(env.ctx, algokit.CallParams[metamerkles.AddDataArgs]{
		Args: metamerkles.AddDataArgs{
			Payment: dataPayment,
			Name:    rootName,
			Key:     dataKey,
			Value:   dataValue,
		},
		Sender: env.deployer.Address,
		Signer: env.deployer.Signer,
		BoxReferences: []types.AppBoxReference{
			{AppID: 0, Name: rootBoxKey},
			{AppID: 0, Name: dataBoxKey},
		},
	})
	if err != nil {
		t.Fatalf("failed to add data: %v", err)
	}

	t.Logf("Added data key '%s' = '%s' to root '%s'", dataKey, dataValue, rootName)
}

func TestMetaMerklesFactory(t *testing.T) {
	env := setupMetaMerklesTest(t)

	client2, _, err := env.factory.Create(env.ctx, algokit.AppFactoryCreateParams{
		ExtraPages: 2,
	})
	if err != nil {
		t.Fatalf("failed to deploy second metamerkles: %v", err)
	}

	if client2.AppID() == env.client.AppID() {
		t.Error("expected different app IDs")
	}

	t.Logf("Second MetaMerkles deployed with appID: %d", client2.AppID())
}

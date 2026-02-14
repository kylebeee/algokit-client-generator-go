package tests

import (
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/abstractedaccount"
	"github.com/kylebeee/algokit-client-generator-go/generated/auction"
	"github.com/kylebeee/algokit-client-generator-go/generated/listing"
	"github.com/kylebeee/algokit-client-generator-go/generated/poll"
	"github.com/kylebeee/algokit-client-generator-go/generated/raffle"
	"github.com/kylebeee/algokit-client-generator-go/generated/stakingpool"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// Factory-deployed contracts require CallerApplicationID to be non-zero,
// meaning they can only be created by another contract (their parent factory).
// Direct deployment will fail. These tests verify that:
// 1. The typed client code compiles correctly
// 2. Direct deployment correctly fails (proving they need CallerApplicationID)

func TestAuctionDirectDeployFails(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := auction.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	_, _, err = factory.Create(fixture.Ctx, auction.FactoryCreateParams{
		Args: auction.CreateArgs{},
	})
	if err == nil {
		t.Fatal("expected direct deploy to fail for factory-deployed contract")
	}
	t.Logf("Auction correctly rejected direct deploy: %v", err)
}

func TestPollDirectDeployFails(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := poll.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	_, _, err = factory.Create(fixture.Ctx, poll.FactoryCreateParams{
		Args: poll.CreateArgs{},
	})
	if err == nil {
		t.Fatal("expected direct deploy to fail for factory-deployed contract")
	}
	t.Logf("Poll correctly rejected direct deploy: %v", err)
}

func TestRaffleDirectDeployFails(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := raffle.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	_, _, err = factory.Create(fixture.Ctx, raffle.FactoryCreateParams{
		Args: raffle.CreateArgs{},
	})
	if err == nil {
		t.Fatal("expected direct deploy to fail for factory-deployed contract")
	}
	t.Logf("Raffle correctly rejected direct deploy: %v", err)
}

func TestListingDirectDeployFails(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := listing.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	_, _, err = factory.Create(fixture.Ctx, listing.FactoryCreateParams{
		Args: listing.CreateArgs{},
	})
	if err == nil {
		t.Fatal("expected direct deploy to fail for factory-deployed contract")
	}
	t.Logf("Listing correctly rejected direct deploy: %v", err)
}

func TestStakingPoolDirectDeployFails(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := stakingpool.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	_, _, err = factory.Create(fixture.Ctx, stakingpool.FactoryCreateParams{
		Args: stakingpool.CreateArgs{},
	})
	if err == nil {
		t.Fatal("expected direct deploy to fail for factory-deployed contract")
	}
	t.Logf("StakingPool correctly rejected direct deploy: %v", err)
}

func TestAbstractedAccountDirectDeployFails(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	factory, err := abstractedaccount.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	_, _, err = factory.Create(fixture.Ctx, abstractedaccount.FactoryCreateParams{
		Args: abstractedaccount.CreateArgs{},
	})
	if err == nil {
		t.Fatal("expected direct deploy to fail for factory-deployed contract")
	}
	t.Logf("AbstractedAccount correctly rejected direct deploy: %v", err)
}

// Compilation verification tests - ensure all factory child types are usable

func TestAuctionTypesCompile(t *testing.T) {
	_ = auction.CreateArgs{}
	_ = auction.FactoryCreateParams{}
	t.Log("Auction types compile correctly")
}

func TestPollTypesCompile(t *testing.T) {
	_ = poll.CreateArgs{}
	_ = poll.FactoryCreateParams{}
	t.Log("Poll types compile correctly")
}

func TestRaffleTypesCompile(t *testing.T) {
	_ = raffle.CreateArgs{}
	_ = raffle.FactoryCreateParams{}
	t.Log("Raffle types compile correctly")
}

func TestListingTypesCompile(t *testing.T) {
	_ = listing.CreateArgs{}
	_ = listing.FactoryCreateParams{}
	t.Log("Listing types compile correctly")
}

func TestStakingPoolTypesCompile(t *testing.T) {
	_ = stakingpool.CreateArgs{}
	_ = stakingpool.FactoryCreateParams{}
	t.Log("StakingPool types compile correctly")
}

func TestAbstractedAccountTypesCompile(t *testing.T) {
	_ = abstractedaccount.CreateArgs{}
	_ = abstractedaccount.FactoryCreateParams{}
	t.Log("AbstractedAccount types compile correctly")
}

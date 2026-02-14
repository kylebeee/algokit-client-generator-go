package tests

import (
	"testing"

	algokit "github.com/kylebeee/algokit-utils-go"

	"github.com/kylebeee/algokit-client-generator-go/generated/akitadao"
	"github.com/kylebeee/algokit-client-generator-go/tests/testutil"
)

// AkitaDAO has a complex CreateArgs with nested structs and [][]interface{} fields
// (RevenueSplits). The go-algorand-sdk ABI encoder cannot infer tuple types from
// empty [][]interface{} slices, so deployment with zero-value args fails.
// These tests verify factory creation and type compilation instead.

func TestAkitaDaoFactoryCreation(t *testing.T) {
	fixture := testutil.NewTestFixture(t)
	owner := fixture.GenerateAccount(testutil.AlgoAmount(100))
	_, err := akitadao.NewFactory(algokit.AppFactoryParams{
		Algorand:      fixture.Algorand,
		DefaultSender: owner.Address,
		DefaultSigner: owner.Signer,
	})
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}
	t.Log("AkitaDAO factory created successfully")
}

func TestAkitaDaoTypedCreateArgs(t *testing.T) {
	// Verify that complex nested struct types compile and can be used
	_ = akitadao.CreateArgs{
		Version:          "1.0.0",
		Akta:             0,
		ContentPolicy:    [36]byte{},
		MinRewardsImpact: 0,
		Apps:             akitadao.AkitaDaoApps{},
		Fees:             akitadao.AkitaDaoFees{},
		ProposalSettings: akitadao.Object752a5b25{},
		RevenueSplits:    [][]interface{}{},
	}
	t.Log("AkitaDAO typed CreateArgs compile correctly")
}

func TestAkitaDaoNestedTypes(t *testing.T) {
	// Verify all nested types in AkitaDAO are accessible and compile
	_ = akitadao.AkitaDaoApps{
		Staking:       1,
		Rewards:       2,
		Pool:          3,
		PrizeBox:      4,
		Subscriptions: 5,
		Gate:          6,
		Auction:       7,
		HyperSwap:     8,
		Raffle:        9,
		MetaMerkles:   10,
		Marketplace:   11,
		Wallet:        12,
	}
	_ = akitadao.ProposalSettings{
		Fee:           100,
		Power:         200,
		Duration:      300,
		Participation: 400,
		Approval:      500,
	}
	_ = akitadao.Object752a5b25{
		UpgradeApp:          akitadao.ProposalSettings{},
		AddPlugin:           akitadao.ProposalSettings{},
		RemoveExecutePlugin: akitadao.ProposalSettings{},
		RemovePlugin:        akitadao.ProposalSettings{},
		AddAllowance:        akitadao.ProposalSettings{},
		RemoveAllowance:     akitadao.ProposalSettings{},
		NewEscrow:           akitadao.ProposalSettings{},
		ToggleEscrowLock:    akitadao.ProposalSettings{},
		UpdateFields:        akitadao.ProposalSettings{},
	}
	_ = akitadao.FactoryCreateParams{}
	t.Log("AkitaDAO nested types compile correctly")
}

func TestAkitaDaoAppSpec(t *testing.T) {
	spec, err := akitadao.GetAppSpec()
	if err != nil {
		t.Fatalf("failed to parse app spec: %v", err)
	}
	if spec.Name != "AkitaDAO" {
		t.Errorf("expected contract name 'AkitaDAO', got '%s'", spec.Name)
	}
	t.Logf("AkitaDAO app spec parsed: %s", spec.Name)
}

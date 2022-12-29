package watchtower

import (
	"fmt"

	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/wallet"
	"github.com/rocket-pool/smartnode/shared/utils/api"
	"github.com/rocket-pool/smartnode/shared/utils/log"
)

// Respond to challenges task
type respondChallenges struct {
	c   *cli.Context
	log log.ColorLogger
	cfg *config.RocketPoolConfig
	w   *wallet.Wallet
	rp  *rocketpool.RocketPool
}

// Create respond to challenges task
func newRespondChallenges(c *cli.Context, logger log.ColorLogger) (*respondChallenges, error) {

	// Get services
	cfg, err := services.GetConfig(c)
	if err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Return task
	return &respondChallenges{
		c:   c,
		log: logger,
		cfg: cfg,
		w:   w,
		rp:  rp,
	}, nil

}

// Respond to challenges
func (t *respondChallenges) run() error {

	// Wait for eth client to sync
	if err := services.WaitEthClientSynced(t.c, true); err != nil {
		return err
	}

	// Get node account
	nodeAccount, err := t.w.GetNodeAccount()
	if err != nil {
		return err
	}

	// Check node trusted status
	nodeTrusted, err := trustednode.GetMemberExists(t.rp, nodeAccount.Address, nil)
	if err != nil {
		return err
	}
	if !nodeTrusted {
		return nil
	}

	// Log
	t.log.Println("Checking for challenges to respond to...")

	// Check for active challenges
	isChallenged, err := trustednode.GetMemberIsChallenged(t.rp, nodeAccount.Address, nil)
	if err != nil {
		return err
	}
	if !isChallenged {
		return nil
	}

	// Log
	t.log.Printlnf("Node %s has an active challenge against it.", nodeAccount.Address.Hex())


}

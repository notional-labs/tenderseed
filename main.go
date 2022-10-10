package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/pex"
	"github.com/tendermint/tendermint/version"
)

var (
	configDir = ".tinyseed"
	logger    = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func main() {
	seedConfig := DefaultConfig()

	chains := GetChains()

	var allchains []Chain
	// Get all chains that seeds
	for _, chain := range chains.Chains {
		current := GetChain(chain)
		allchains = append(allchains, current)
	}

	// Seed each chain
	for i, chain := range allchains {
		// increment the port number
		address := "tcp://0.0.0.0:" + fmt.Sprint(9000+i)

		// make folders and files
		nodeKey := MakeFolders(chain, seedConfig)

		// allpeers is a slice of peers
		var allpeers []string
		// make the struct of peers into a string
		for _, peer := range chain.Peers.PersistentPeers {
			thispeer := peer.ID + "@" + peer.Address

			// check to make sure there's a colon in the address
			if !strings.Contains(thispeer, ":") {
				continue
			}

			// ensure that there is only one ampersand in the address
			if strings.Count(thispeer, "@") != 1 {
				continue
			}
			allpeers = append(allpeers, thispeer) //nolint:staticcheck
		}

		// set the configuration
		seedConfig.ChainID = chain.ChainID
		seedConfig.Seeds = allpeers
		seedConfig.ListenAddress = address

		// give the user addresses where we are seeding
		logger.Info("Starting Seed Node for " + chain.ChainID + " on " + string(nodeKey.ID()) + "@0.0.0.0:" + fmt.Sprint(9000+i))

		go Start(seedConfig, &nodeKey)
		time.Sleep(1 * time.Second)

	}
}

// make folders and files
func MakeFolders(chain Chain, seedConfig *Config) (nodeKey p2p.NodeKey) {
	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	// init config directory & files
	homeDir := filepath.Join(userHomeDir, configDir+"/"+chain.ChainID)
	configFilePath := filepath.Join(homeDir, "config.toml")
	addrBookFilePath := filepath.Join(homeDir, "addrbook.json")
	nodeKeyFilePath := filepath.Join(homeDir, "node_key.json")

	// Make folders
	for _, path := range []string{configFilePath, addrBookFilePath, nodeKeyFilePath} {
		err := os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			panic(err)
		}
	}

	nk, err := p2p.LoadOrGenNodeKey(nodeKeyFilePath)
	if err != nil {
		panic(err)
	}

	return *nk
}

// Start starts a Tenderseed
func Start(seedConfig *Config, nodeKey *p2p.NodeKey) {
	chainID := seedConfig.ChainID
	cfg := config.DefaultP2PConfig()
	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	filteredLogger := log.NewFilter(logger, log.AllowInfo())

	logger.Info("Configuration",
		"chain", chainID,
		"key", nodeKey.ID(),
		"node listen", seedConfig.ListenAddress,
		"max-inbound", seedConfig.MaxNumInboundPeers,
		"max-outbound", seedConfig.MaxNumOutboundPeers,
	)

	protocolVersion := p2p.NewProtocolVersion(
		version.P2PProtocol,
		version.BlockProtocol,
		0,
	)

	// NodeInfo gets info on your node
	nodeInfo := p2p.DefaultNodeInfo{
		ProtocolVersion: protocolVersion,
		DefaultNodeID:   nodeKey.ID(),
		ListenAddr:      seedConfig.ListenAddress,
		Network:         chainID,
		Version:         "0.6.9",
		Channels:        []byte{pex.PexChannel},
		Moniker:         fmt.Sprintf("%s-seed", chainID),
	}

	addr, err := p2p.NewNetAddressString(p2p.IDAddressString(nodeInfo.DefaultNodeID, nodeInfo.ListenAddr))
	if err != nil {
		panic(err)
	}

	transport := p2p.NewMultiplexTransport(nodeInfo, *nodeKey, p2p.MConnConfig(cfg))
	if err := transport.Listen(*addr); err != nil {
		panic(err)
	}

	addrBookFilePath := filepath.Join(userHomeDir, configDir, seedConfig.AddrBookFile)
	book := pex.NewAddrBook(addrBookFilePath, seedConfig.AddrBookStrict)
	//	book.SetLogger(filteredLogger.With("module", "book"))

	pexReactor := pex.NewReactor(book, &pex.ReactorConfig{
		SeedMode:                     true,
		Seeds:                        seedConfig.Seeds,
		SeedDisconnectWaitPeriod:     5 * time.Second, // default is 28 hours, we just want to harvest as many addresses as possible
		PersistentPeersMaxDialPeriod: 0,               // use exponential back-off
	})
	//	pexReactor.SetLogger(filteredLogger.With("module", "pex"))

	sw := p2p.NewSwitch(cfg, transport)
	sw.SetLogger(filteredLogger.With("module", "switch"))
	sw.SetNodeKey(nodeKey)
	sw.SetAddrBook(book)
	sw.AddReactor("pex", pexReactor)

	// last
	sw.SetNodeInfo(nodeInfo)

	err = sw.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		// Fire periodically
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			peersout, peersin, dialing := sw.NumPeers()
			fmt.Println(seedConfig.ChainID, peersout, " outbound peers, ", peersin, " inbound peers, and ", dialing, " dialing peers")
		}
	}()

	sw.Wait()
	// if we block here, we just get the first chain.
}

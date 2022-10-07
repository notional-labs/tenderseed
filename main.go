package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	userHomeDir, err := homedir.Dir()
	seedConfig := DefaultConfig()
	if err != nil {
		panic(err)
	}

	chains := GetChains()

	var allchains []Chain
	// Get all chains that seeds
	for _, chain := range chains.Chains {
		current := GetChain(chain)
		allchains = append(allchains, current)
		if err != nil {
			panic(err)
		}
	}

	nodeKeyFilePath := filepath.Join(userHomeDir, configDir, "config", seedConfig.NodeKeyFile)
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFilePath)
	if err != nil {
		panic(err)
	}

	// Seed each chain

	for i, chain := range allchains {
		// increment the port number
		port := 7000 + i
		address := "tcp://0.0.0.0:" + fmt.Sprint(port)

		peers := chain.Peers.PersistentPeers
		seeds := chain.Peers.Seeds
		// make the struct of seeds into a string
		var allseeds []string
		for _, seed := range seeds {
			allseeds = append(allseeds, seed.ID+"@"+seed.Address) //nolint:staticcheck
		}

		// allpeers is a slice of peers
		var allpeers []string
		// make the struct of peers into a string
		for _, peer := range peers {
			allpeers = append(allpeers, peer.ID+"@"+peer.Address) //nolint:staticcheck
		}

		// set the configuration
		seedConfig.ChainID = chain.ChainID
		seedConfig.Seeds = append(seedConfig.Peers, seedConfig.Seeds...)
		seedConfig.ListenAddress = address

		// init config directory & files
		homeDir := filepath.Join(userHomeDir, configDir+"/"+chain.ChainID, "config")
		configFilePath := filepath.Join(homeDir, "config.toml")
		nodeKeyFilePath := filepath.Join(homeDir, seedConfig.NodeKeyFile)
		addrBookFilePath := filepath.Join(homeDir, seedConfig.AddrBookFile)

		// Make folders
		err = os.MkdirAll(filepath.Dir(nodeKeyFilePath), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(filepath.Dir(addrBookFilePath), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(filepath.Dir(configFilePath), os.ModePerm)
		if err != nil {
			panic(err)
		}

		// give the user addresses where we are seeding
		logger.Info("Starting Seed Node for " + chain.ChainID + " on " + string(nodeKey.ID()) + "@0.0.0.0:" + fmt.Sprint(port))
		Start(seedConfig, nodeKey)
	}
}

// Start starts a Tenderseed
func Start(seedConfig *Config, nodeKey *p2p.NodeKey) {
	chainID := seedConfig.ChainID
	cfg := config.DefaultP2PConfig()
	cfg.AllowDuplicateIP = true

	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	filteredLogger := log.NewFilter(logger, log.AllowInfo())

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

	addrBookFilePath := filepath.Join(userHomeDir, configDir, "config", seedConfig.AddrBookFile)
	book := pex.NewAddrBook(addrBookFilePath, seedConfig.AddrBookStrict)
	book.SetLogger(filteredLogger.With("module", "book"))

	pexReactor := pex.NewReactor(book, &pex.ReactorConfig{
		SeedMode:                     true,
		Seeds:                        seedConfig.Seeds,
		SeedDisconnectWaitPeriod:     1 * time.Second, // default is 28 hours, we just want to harvest as many addresses as possible
		PersistentPeersMaxDialPeriod: 0,               // use exponential back-off
	})
	pexReactor.SetLogger(filteredLogger.With("module", "pex"))

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
			logger.Info(seedConfig.ChainID, "has peers", sw.Peers().List())
			book.Save()
		}
	}()
}

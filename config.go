package main

// Config defines the configuration format
type Config struct {
	ListenAddress       string   `toml:"laddr" comment:"Address to listen for incoming connections"`
	ChainID             string   `toml:"chain_id" comment:"network identifier (todo move to cli flag argument? keeps the config network agnostic)"`
	NodeKeyFile         string   `toml:"node_key_file" comment:"path to node_key (relative to tendermint-seed home directory or an absolute path)"`
	AddrBookFile        string   `toml:"addr_book_file" comment:"path to address book (relative to tendermint-seed home directory or an absolute path)"`
	AddrBookStrict      bool     `toml:"addr_book_strict" comment:"Set true for strict routability rules\n Set false for private or local networks"`
	MaxNumInboundPeers  int      `toml:"max_num_inbound_peers" comment:"maximum number of inbound connections"`
	MaxNumOutboundPeers int      `toml:"max_num_outbound_peers" comment:"maximum number of outbound connections"`
	Seeds               []string `toml:"seeds" comment:"seed nodes we can use to discover peers"`
	Peers               []string `toml:"persistent_peers" comment:"persistent peers we will always keep connected to"`
}

// DefaultConfig returns a seed config initialized with default values
func DefaultConfig() *Config {
	return &Config{
		ListenAddress:       "tcp://0.0.0.0:6969",
		NodeKeyFile:         "node_key.json",
		AddrBookFile:        "addrbook.json",
		AddrBookStrict:      true,
		MaxNumInboundPeers:  3000,
		MaxNumOutboundPeers: 100,
	}
}

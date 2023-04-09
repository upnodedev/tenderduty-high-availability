package ha

type HaState struct {
	ChainName      string
	State          string
	Status         string
	Jailed         bool
	JailedNotified bool
	Retry          int
}

type HaConfig struct {
	ListenPort      string `yaml:"listen_port"`
	NodeIndex       int    `yaml:"node_index"`
	AnotherEndpoint int    `yaml:"another_endpoint"`
}

type HaChainConfig struct {
	ServiceName string `yaml:"service_name"`
}

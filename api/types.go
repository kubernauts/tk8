package api

type AllClusters map[string][]Cluster

type Cluster interface {
	CreateCluster() error
	DestroyCluster() error
	GetCluster(name string) (Cluster, error)
}

type ConfigStore interface {
	CreateConfig(Cluster) error
	DeleteConfig() error
	UpdateConfig() error
	ValidateConfig() error
	CheckConfigExists() (bool, error)
	GetConfig() ([]byte, error)
	GetConfigs() (AllClusters, error)
}

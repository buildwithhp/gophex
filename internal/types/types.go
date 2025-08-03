package types

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type         string // mysql, postgresql, mongodb
	ConfigType   string // cluster, multi-cluster, read-write
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	ReadHost     string   // for read-write setup
	WriteHost    string   // for read-write setup
	ClusterNodes []string // for multi-cluster
	SSLMode      string
	AuthSource   string // for MongoDB
	ReplicaSet   string // for MongoDB
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	Database int
}

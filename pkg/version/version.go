package version

var (
	Version = "1.0.0"
	Commit  = "dev"
	Date    = "unknown"
)

func GetVersion() string {
	return Version
}

func GetFullVersion() string {
	return Version + " (" + Commit + ") built on " + Date
}

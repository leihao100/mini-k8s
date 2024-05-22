package config

const version string = "V1.0"
const (
	etcdHost string = "localhost"
	etcdPort string = ":2379"
)

func Version() string {
	return version
}
func EtcdHost() string {
	return etcdHost
}
func EtcdPort() string {
	return etcdPort
}

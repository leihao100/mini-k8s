package config

const version string = "V1.0"
const (
	etcdHost string = "localhost"
	etcdPort string = ":2379"
)
const (
	apiServerHost string = "http://localhost"
	apiServerPort string = ":8080"
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
func ApiServerHost() string {
	return apiServerHost
}
func ApiServerPort() string {
	return apiServerPort
}

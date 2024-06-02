package config

const version string = "V1.0"
const (
	etcdHost string = "10.119.13.182"
	etcdPort string = ":2379"
)

var (
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
func SetApiServerHost(host string) {
	apiServerHost = host
}
func SetApiServerPort(port string) {
	apiServerPort = port
}

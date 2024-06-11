package url

const ApiServerURL string = ""
const CleanURL string = "/api/v1/clear"

type URL struct {
	version string
	kind    string
}

func (url *URL) Init(version string, kind string) {
	url.version = version
	url.kind = kind
}

func (url *URL) CreateURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/create"
	return res
}

func (url *URL) DeleteURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/delete"
	return res
}

func (url *URL) GetURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/get"
	return res
}

func (url *URL) WatchURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/watch"
	return res
}

func (url *URL) StatusURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/status"
	return res
}

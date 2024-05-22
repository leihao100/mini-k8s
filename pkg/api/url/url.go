package url

const apiServerURL string = ""
const clearURL string = "/api/v1/clear"

type URL struct {
	version string
	kind    string
}

func (url *URL) Init(version string, kind string) {
	url.version = version
	url.kind = kind
}
func (url *URL) createURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/create"
	return res
}

func (url *URL) deleteURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/delete"
	return res
}

func (url *URL) getURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/get"
	return res
}

func (url *URL) watchURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/watch"
	return res
}

func (url *URL) statusURL() (res string) {
	res = "/api" + "/" + url.version + "/" + url.kind + "/status"
	return res
}

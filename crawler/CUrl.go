package crawler

type CUrl struct {
	url  string
	path string
}

func (aria *CUrl) Cmd() string {
	return `curl`
}

func (aria *CUrl) Args() []string {
	return []string{aria.url, `-o`, aria.path}
}

func (aria *CUrl) SetUrl(url string) ICrawler {
	aria.url = url
	return aria
}

func (aria *CUrl) SetPath(path string) ICrawler {
	aria.path = path
	return aria
}

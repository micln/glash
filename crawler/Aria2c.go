package crawler

type Aria2c struct {
	url  string
	path string
}

func (aria *Aria2c) Cmd() string {
	return `aria2c`
}

func (aria *Aria2c) Args() []string {
	return []string{aria.url, `-o`, aria.path}
}

func (aria *Aria2c) SetUrl(url string) ICrawler {
	aria.url = url
	return aria
}

func (aria *Aria2c) SetPath(path string) ICrawler {
	aria.path = path
	return aria
}

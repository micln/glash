package crawler

type ICrawler interface {
	SetUrl(string) ICrawler
	SetPath(string) ICrawler

	Cmd() string
	Args() []string
}

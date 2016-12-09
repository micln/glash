package crawler

func Which(who string) ICrawler {
	switch who {

	case `aria2`:
		return &Aria2c{}

	default:
		return &CUrl{}

	}
}

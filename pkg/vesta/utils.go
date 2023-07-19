package vesta

func Screenshot(url string) ([]byte, error) {
	v := New().Nav(url)
	defer v.Cancel()
	return v.GetScreen()
}

package scrapper_test

import (
	"testing"

	"github.com/snake0207/scrap/scrapper"
)

func TestGetPages(t *testing.T) {
	count := scrapper.GetPages("python")

	if count != 10 {
		t.Error("Pages count wrong")
	}
}
package client

import "testing"

func TestToken(t *testing.T) {
	s := GetTimestamp()

	token := CreateRequestToken("m198sOkJEn37DjqZ32lpRu76xmw288xSQ9", "1373209025")
	if token != "9301c956749167186ee713e4f3a3d90446e84d8d19a4ca8ea9b4b314d1c51b7b" {
		t.Fail()
	}
}

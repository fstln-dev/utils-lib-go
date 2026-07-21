package graphql

import "testing"

func TestIdByGraphQL(t *testing.T) {
	//gid := ""
	//gid := "123456789"
	//gid := "gid://shopify/asdfasfasf"
	gid := "gid://shopify/Product/123456789"
	id := IdByGraphQL(gid)
	t.Log(id)
}

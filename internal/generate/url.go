package generate

var domainName = "https://zipfy.xyz"

func ShortenURL(id int) string {
	return domainName + UUIDFromInt(id)
}

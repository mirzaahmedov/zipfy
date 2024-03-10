package generate

var domainName = "http://localhost:8080/"

func ShortenURL(id int) string {
	return domainName + UUIDFromInt(id)
}

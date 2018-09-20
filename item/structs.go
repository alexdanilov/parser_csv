package item

import "strings"


type Item struct {
	Id string
	Name string
	Email string
	Phone string
}


// normalize data
func (i *Item) Normalize(code string) {
	replacer := strings.NewReplacer("(", "", ")", "", " ", "")
	i.Phone = replacer.Replace(i.Phone)

	// add country code for local numbers (starts with 0)
	i.Phone = code + i.Phone
}

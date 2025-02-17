package constants

type role struct {
	Owner    string
	Manager  string
	Customer string
}

var Roles = role{
	Owner:    "110",
	Manager:  "010",
	Customer: "001",
}

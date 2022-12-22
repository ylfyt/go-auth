package dtos

type Error struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

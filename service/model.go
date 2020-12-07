package service

type Add struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddAck struct {
	Res int `json:"res"`
}

type Login struct {
	Name string `json:"name"`
	Pass string `json:"password"`
}

type LoginAck struct {
	Token string `json:token`
}
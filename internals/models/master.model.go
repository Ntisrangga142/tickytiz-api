package models

type MasterDirector struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MasterActor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MasterLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MasterTime struct {
	ID   int    `json:"id"`
	Time string `json:"time"`
}

type MasterCinema struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Logo  string `json:"logo"`
	Price int    `json:"price"`
}

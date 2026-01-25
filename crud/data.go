package crud

type Category struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Description string `json:"description"`
}

var Category1 = []Category{
	{ID: 1, Nama: " ", Description: " "},
}

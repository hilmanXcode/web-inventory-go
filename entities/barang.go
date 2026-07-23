package entities

type BarangCreate struct {
	Sku        string `validator:"required" input_name:"sku"`
	NamaBarang string `validator:"required" input_name:"nama"`
	Kategori   string `validator:"required" input_name:"kategori"`
	Stock      string `validator:"required,number" input_name:"stock"`
}

package entities

type User struct {
	ID          int    `validator:"nullable,number"`
	NamaLengkap string `validator:"required,string"`
	Email       string `validator:"required"`
	Password    string `validator:"required"`
	Role        string `validator:"required"`
	Umur        int    `validator:"required,number"`
}

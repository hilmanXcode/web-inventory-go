package entities

type User struct {
	ID          int    `validator:"nullable,number" input_name:"id"`
	NamaLengkap string `validator:"required,string" input_name:"Nama Lengkap"`
	Email       string `validator:"required" input_name:"Email"`
	Password    string `validator:"required" input_name:"Password"`
	Role        string `validator:"required" input_name:"Role"`
}

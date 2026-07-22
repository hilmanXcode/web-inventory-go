package entities

type UserRegister struct {
	ID          int    `validator:"nullable,number" input_name:"id"`
	NamaLengkap string `validator:"required,string" input_name:"Nama Lengkap"`
	Email       string `validator:"required" input_name:"Email"`
	Password    string `validator:"required,match[password.confirm_password]" input_name:"Password"`
	Role        string `validator:"required" input_name:"Role"`
}

type UserLogin struct {
	Email    string `validator:"required" input_name:"Email"`
	Password string `validator:"required" input_name:"Password"`
}

package models

//Rol Model
type Rol struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" faker:"-"`
	Rol         string `gorm:"unique" faker:"name"`
	Alias       string `faker:"first_name"`
	Description string `faker:"paragraph"`
	Users       []User `gorm:"foreignKey:RolID" faker:"-" json:"-"`
}

var (
	Administrator = &Rol{
		Rol:         "admin",
		Alias:       "Administrador",
		Description: "Encargado de la administracion",
	}
	Operator = &Rol{
		Rol:         "operator",
		Alias:       "Operador",
		Description: "Encargado de operacion",
	}
)

package hotel


type Repository interface {}


type repository struct {
	db string
}


func NewRepository(db string) Repository {
	return repository{db};
}





package domain


type User struct {
	ID          int64
	Version     int
	FullName    string 
	PhoneNumber *string
}
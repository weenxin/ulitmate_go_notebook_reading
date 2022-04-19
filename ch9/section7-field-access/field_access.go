package section7_field_access

import "fmt"


type UserAble interface{
	GetName() string
	GetEmail() string
	SetId(int64)
}


type User struct {
	Id int64
	Name string
	Email string
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetEmail() string {
	return  u.Email
}

func (u* User) SetId( id int64) {
	u.Id = id
}


type Customer struct {
	Id int64
	Name string
	Email string
}

func (u Customer) GetName() string {
	return u.Name
}

func (u Customer) GetEmail() string {
	return  u.Email
}

func (u* Customer) SetId( id int64) {
	u.Id = id
}

type Entities interface{
	*User | *Customer
	UserAble
}

//https://github.com/golang/go/issues/48522 目前暂时不支持



func InsertData[T Entities](entity T) (T ,error) {
	fmt.Printf("Insert data: name :%s , email :%s \n" ,entity.GetName(),entity.GetEmail())
	entity.SetId(1000)
	return entity,nil
}








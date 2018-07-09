package def

type Human struct {
	age int
	sex int
	name string
}

func (self *Human) Age() int {
	return self.age
}

func (self *Human) Sex() int {
	return self.sex
}

func (self *Human) Name() string {
	return self.name
}

package accounts

import "errors"

var errDiscount error = errors.New("can't discount")

// account struct
type account struct {
	name string
	balance int
}

func NewAccount(name string) *account {
	return &account{name: name, balance: 10}
}

func (a *account) Deposit(balance int) {
	a.balance += balance
}

func (a *account) Discount(balance int) error {
	if a.balance < balance {
		return errDiscount
	}
	a.balance -= balance
	return nil
}

func (a account) GetName() string {
	return a.name
}

func (a *account) ChangeName(name string) {
	a.name = name
}

func (a account) String() string {
	return "wertyuiop"
}


package models

import "github.com/zebresel-com/mongodm"

type Address struct {
	Street string `json:"street" bson:"street"`
	City   string `json:"city" bson:"city"`
	Zip    int    `json:"zip" bson:"zip"`
}

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	FirstName    string   `json:"firstname"  bson:"firstname" minLen:"2" maxLen:"30" required:"true"`
	LastName     string   `json:"lastname"  bson:"lastname" minLen:"2" maxLen:"30" required:"true"`
	UserName     string   `json:"username"  bson:"username" minLen:"2" maxLen:"15"`
	Email        string   `json:"email" bson:"email" validation:"email" required:"true"`
	PasswordHash string   `json:"-" bson:"passwordHash"`
	Address      *Address `json:"address" bson:"address"`
}

func (self *User) Validate(values ...interface{}) (bool, []error) {

	var valid bool
	var validationErrors []error

	valid, validationErrors = self.DefaultValidate()

	type m map[string]string

	if len(values) > 0 {

		//expect password as first param then validate it with the next rules
		if password, ok := values[0].(string); ok {

			if len(password) < 8 {

				self.AppendError(&validationErrors, mongodm.L("validation.field_minlen", "password", 8))

			} else if len(password) > 50 {

				self.AppendError(&validationErrors, mongodm.L("validation.field_maxlen", "password", 50))
			}

		} else {

			self.AppendError(&validationErrors, mongodm.L("validation.field_required", "password"))
		}
	}

	if len(validationErrors) > 0 {
		valid = false
	}

	return valid, validationErrors
}

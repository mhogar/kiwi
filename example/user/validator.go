package user

import (
	"github.com/mhogar/kiwi/common"
	"github.com/mhogar/kiwi/example/models"
	"github.com/mhogar/kiwi/nodes"
	"github.com/mhogar/kiwi/nodes/validator"
)

type userValidator struct{}

func newUserValidator() userValidator {
	return userValidator{}
}

func (v userValidator) ValidatePasswordComplexity(_ interface{}, val interface{}) (*validator.ValidationErrors, error) {
	password := val.(PasswordField).GetNewPassword()
	return validator.ValidatePassword("password", password, 8, 0, true, true), nil
}

func (userValidator) ValidateUser(_ interface{}, val any) (*validator.ValidationErrors, error) {
	user := val.(*models.User)
	verrs := &validator.ValidationErrors{}

	verrs.Merge(validator.ValidateLength("username", user.Username, 5, 30))
	verrs.Merge(validator.ValidateMin("rank", user.Rank, 0))

	return verrs, nil
}

func (userValidator) ValidateUserUnique(ctx interface{}, val any) (*validator.ValidationErrors, error) {
	user := val.(*models.User)

	verrs, err := validator.ValidateUniqueField(user, ctx.(nodes.Context).GetDataAdapter(), "already in use by another user")
	if err != nil {
		return verrs, common.ChainError("error validating user unique field", err)
	}

	return verrs, nil
}

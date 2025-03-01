package users

import (
	"time"

	"github.com/go-msvc/errors"
)

type AddUserRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date-of-birth" doc:"Date of birth formatted as CCYY-MM-DD"`
}

func (req AddUserRequest) Validate() error {
	if req.Name == "" {
		return errors.Error("missing name")
	}
	if req.DateOfBirth == "" {
		return errors.Error("missing date-of-birth")
	}
	if _, err := time.Parse("2006-01-02", req.DateOfBirth); err != nil {
		return errors.Errorf("date-of-birth:\"%s\" not formatted as CCYY-MM-DD", req.DateOfBirth)
	}
	return nil
}

type UpdateUserRequest struct {
	DateOfBirth *string  `json:"date-of-birth,omitempty" doc:"Date of birth formatted as CCYY-MM-DD"`
	Address     *Address `json:"address,omitempty"`
}

func (req UpdateUserRequest) Validate() error {
	count := 0
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		if _, err := time.Parse("2006-01-02", *req.DateOfBirth); err != nil {
			return errors.Errorf("date-of-birth:\"%s\" not formatted as CCYY-MM-DD", *req.DateOfBirth)
		}
		count++
	}
	if req.Address != nil {
		if err := req.Address.Validate(); err != nil {
			return errors.Wrap(err, "invalid address")
		}
		count++
	}
	if count == 0 {
		return errors.Error("missing both date-of-birth and address")
	}
	return nil
}

type Address struct {
	Street  string
	Country string
}

func (addr Address) Validate() error {
	if addr.Street == "" {
		return errors.Error("missing street")
	}
	if addr.Country == "" {
		return errors.Error("missing country")
	}
	return nil
}

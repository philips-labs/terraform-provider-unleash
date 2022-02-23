package provider

import (
	"errors"
)

// Exported Errors
var (
	ErrStrategyParametersRequired = errors.New("all the strategy parameters values need to be informed")
	ErrPercentageConvertion       = errors.New("the parameter of type percentage could not be converted, please make sure its a number in string format without %")
	ErrNumberConvertion           = errors.New("the parameter of type number could not be converted, please make sure its a number in string format")
	ErrBooleanConvertion          = errors.New("the parameter of type boolean could not be converted, please make sure its true or false in string format")
)

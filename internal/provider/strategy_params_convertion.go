package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func newParameterInputConvertion() map[string]func(string) interface{} {
	inputTypesMap := make(map[string]func(string) interface{})
	inputTypesMap["string"] = func(value string) interface{} {
		return value
	}
	inputTypesMap["percentage"] = func(value string) interface{} {
		if perc, err := strconv.ParseFloat(value, 64); err == nil {
			return perc
		}
		return diag.FromErr(ErrPercentageConvertion)
	}
	inputTypesMap["number"] = func(value string) interface{} {
		if number, err := strconv.Atoi(value); err == nil {
			return number
		}
		return diag.FromErr(ErrNumberConvertion)
	}
	inputTypesMap["boolean"] = func(value string) interface{} {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		return diag.FromErr(ErrBooleanConvertion)
	}
	inputTypesMap["list"] = func(value string) interface{} {
		return strings.Split(strings.TrimSpace(value), ",")
	}
	return inputTypesMap
}

func newParameterOutputConvertion() map[string]func(interface{}) string {
	outputTypesMap := make(map[string]func(interface{}) string)
	outputTypesMap["string"] = func(value interface{}) string {
		return value.(string)
	}
	outputTypesMap["percentage"] = func(value interface{}) string {
		return fmt.Sprintf("%g", value.(float64))
	}
	outputTypesMap["number"] = func(value interface{}) string {
		return strconv.Itoa(value.(int))
	}
	outputTypesMap["boolean"] = func(value interface{}) string {
		return strconv.FormatBool(value.(bool))
	}
	outputTypesMap["list"] = func(value interface{}) string {
		return strings.Join(value.([]string), ",")
	}
	return outputTypesMap
}

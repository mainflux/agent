/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package models

import (
	"reflect"
)

const (
	ValidateTag = "validate"
)

// Validator provides an interface for struct types to implement validation of their internal state. They can also
// indicate to a caller whether their validation has already been completed.
//
// NOTE: This cannot be applied to struct types that are simply aliased to a primitive.
type Validator interface {
	// Validate performs integrity checks on the internal state of the model. It returns a boolean indicating whether
	// the validation passed or not, and the associated error if validation was unsuccessful.
	Validate() (bool, error)
}

func validate(t interface{}) error {
	val := reflect.ValueOf(t)
	typ := reflect.TypeOf(t)
	fields := val.NumField()
	for f := 0; f < fields; f++ {
		field := val.Field(f)
		typfield := typ.Field(f)
		if field.Type().NumMethod() > 0 && field.CanInterface() && typfield.Tag.Get(ValidateTag) != "-" {
			if v, ok := field.Interface().(Validator); ok {
				cast := v.(Validator)
				_, err := cast.Validate()
				if err != nil {
					return NewErrContractInvalid(err.Error())
				}
			}
		}
	}
	return nil
}

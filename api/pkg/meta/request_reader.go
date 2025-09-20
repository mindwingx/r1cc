package meta

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"microservice/pkg/meta/status"
	"microservice/pkg/validator"
	"reflect"
	"strings"
)

type IRequestConvertible[D any] interface {
	ToDomain() D
}

// ReqBodyToDomain binds the request body and evaluates that by `GO Validator`
func ReqBodyToDomain[R IRequestConvertible[D], D any](c echo.Context) (entity D, err error) {
	var body R
	if err = c.Bind(&body); err != nil {
		err = DtoBindErr
		return
	}

	if err = validator.ValidateRequestDto(c.Request().Context(), body); err != nil {
		err = ServiceErr(status.Validate, err)
		return
	}

	entity = body.ToDomain()

	return
}

// ReqRouteParamsToDomain binds the request body and evaluates that by `GO Validator`
func ReqRouteParamsToDomain[R IRequestConvertible[D], D any](c echo.Context) (entity D, err error) {
	var params R

	// Create a new instance of the params struct
	params = reflect.New(reflect.TypeOf(params).Elem()).Interface().(R)
	bindRouteParams(c, params)

	// Validate the parameters
	if err = validator.ValidateRequestDto(c.Request().Context(), params); err != nil {
		err = ServiceErr(status.Validate, err)
		return
	}

	// Convert to domain entity
	entity = params.ToDomain()
	return
}

func ReqQryParamToDomain[R IRequestConvertible[D], D any](c echo.Context) (entity D, err error) {
	var qry R
	qry = reflect.New(reflect.TypeOf(qry).Elem()).Interface().(R)

	if err = c.Bind(qry); err != nil {
		err = DtoBindErr
		return
	}

	if err = validator.ValidateRequestDto(c.Request().Context(), qry); err != nil {
		err = ServiceErr(status.Validate, err)
		return
	}

	entity = qry.ToDomain()
	return
}

func ReqHeaderToDomain[R IRequestConvertible[D], D any](c echo.Context) (entity D, err error) {
	var headers R

	err = mapHeadersToStruct(c, &headers)
	if err != nil {
		return
	}

	if err = validator.ValidateRequestDto(c.Request().Context(), headers); err != nil {
		err = ServiceErr(status.Validate, err)
		return
	}

	entity = headers.ToDomain()

	return
}

// HELPERS

// bindRouteParams to bind route parameters to a struct
func bindRouteParams(c echo.Context, target interface{}) {
	paramNames := c.ParamNames()

	// Use reflection to set the struct fields
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Check for echo tag or use lowercase field name as default
		paramName := field.Tag.Get("param")
		if paramName == "" {
			paramName = strings.ToLower(field.Name)
		}

		for _, name := range paramNames {
			if strings.EqualFold(name, paramName) {
				paramValue := c.Param(name)

				if len(paramValue) > 0 {
					fieldValue.SetString(paramValue)
				}

				break
			}
		}
	}
}

// mapHeadersToStruct: match and map the headers' values to the referenced struct JSON tags
func mapHeadersToStruct(c echo.Context, dest interface{}) error {
	// Get the type and value of the struct
	v := reflect.ValueOf(dest)

	// Check if dest is a pointer
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	// Handle the case where we might have a pointer to a pointer
	// This can happen with generic types like *dto.TenantUuid
	if v.Elem().Kind() == reflect.Ptr {
		// If it's a pointer to a pointer, we need to allocate the inner pointer
		if v.Elem().IsNil() {
			// Create a new instance of the struct
			elemType := v.Elem().Type().Elem()
			newVal := reflect.New(elemType)
			v.Elem().Set(newVal)
		}
		// Now dereference the outer pointer to get the inner pointer
		v = v.Elem()
	}

	// Now dereference the pointer to get the underlying struct
	v = v.Elem()

	// Check if the dereferenced value is a struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < v.NumField(); i++ {
		// Get the field's header tag
		headerTag := t.Field(i).Tag.Get("header")
		if headerTag == "" {
			continue
		}

		// Get the header value
		headerValue := c.Request().Header.Get(headerTag)

		// Handle bearer token scenario
		if strings.Contains(headerValue, "Bearer") || strings.Contains(headerValue, "bearer") {
			values := strings.Split(headerValue, " ")
			if len(values) > 1 {
				headerValue = values[1]
			}
		}

		// Set the value if it's not empty
		if headerValue != "" {
			field := v.Field(i)

			// Check if the field is settable and is a string type
			if field.CanSet() && field.Kind() == reflect.String {
				field.SetString(headerValue)
			}
		}
	}

	return nil
}

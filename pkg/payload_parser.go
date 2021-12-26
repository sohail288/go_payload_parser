package payload_parser

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	validatorFmt = "Validate%s"
	applyFmt     = "Tranform%s"
)

var structTagTypes = []string{
	"query",
	"header",
}

type QueryStructtag struct {
	name          string
	required      bool
	serialization string // SerializationType?
	defaultValue  interface{}
}

func parseRequired(requiredVal string) (bool, error) {
	if requiredVal == "required" {
		return true, nil
	}
	if requiredVal == "-" {
		return false, nil
	}
	return false, fmt.Errorf("invalid value for required")
}

func ParsePayload(shape interface{}, req *http.Request) error {
	v := reflect.ValueOf(shape)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("expected ptr as first argument, try using &%s\n", v.Type())
	}

	elem := v.Elem()
	elemType := elem.Type()

	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("expected root kind to be struct, got %s\n", elem.Kind())
	}

	numFields := elem.NumField()
	queryValues := req.URL.Query()
	headers := req.Header

	for i := 0; i < numFields; i++ {
		f := elem.Field(i)
		log.Printf("i=%d name=%s type=%s value=%s\n", i, elemType.Field(i).Name, f.Type(), f.Interface())

		//  this only works on concrete type? not value
		tField := elemType.Field(i)
		for _, tag := range structTagTypes {
			if tagvalue, ok := tField.Tag.Lookup(tag); ok {
				tagvalues := strings.Split(tagvalue, ",")

				if len(tagvalues) < 2 {
					return fmt.Errorf("tag for: %s, is not complete, requires name and required", tag)
				}

				paramName := tagvalues[0]
				required, err := parseRequired(tagvalues[1])

				if err != nil {
					return fmt.Errorf("tag for %s, contains invalid required value, use required or -", tField.Name)
				}

				var rawVal string
				switch tag {
				case "query":
					// parse struct tag for query or have a parser that returns a struct?
					if vals, ok := queryValues[paramName]; ok {
						rawVal = vals[0]
					} else {
						if required {
							return fmt.Errorf("%s is required but is missing from input", paramName)
						}
					}
				case "header":
					// parse struct tag for header
					if vals, ok := headers[http.CanonicalHeaderKey(paramName)]; ok {
						rawVal = vals[0]
					} else {
						if required {
							return fmt.Errorf("%s is required but is missing from input", paramName)
						}
					}
				}

				switch f.Kind() {
				case reflect.Int:
					val, err := strconv.ParseInt(rawVal, 10, 64)
					if err != nil {
						return fmt.Errorf("Unable to parse integer val for %s: %s", paramName, rawVal)
					}
					f.SetInt(val)
				case reflect.Bool:
					val, err := strconv.ParseBool(rawVal)
					if err != nil {
						return fmt.Errorf("Unable to parse bool val for %s: %s", paramName, rawVal)
					}
					f.SetBool(val)
				case reflect.String:
					f.SetString(rawVal)
				default:
					return fmt.Errorf("Unexpected type found for param %s: %s", paramName, f.Kind())
				}
			}
		}

		validatorName := fmt.Sprintf(validatorFmt, elemType.Field(i).Name)
		//fmt.Printf("Looking for %s\n", validatorName)

		// depending on ptr receiver or value receiver, so would need to check both?
		// ValidateX = value receiver
		// ApplyX - ptr receiver?
		// can we have multiple validators for a single field? order will be tricky
		// would need to iterate over all methods? (ie check if field starts with format)

		validatorMethod := v.MethodByName(validatorName)
		if validatorMethod.IsValid() {
			out := validatorMethod.Call([]reflect.Value{})
			if len(out) == 0 {
				log.Printf("expected non-zero length")
				return fmt.Errorf("expected non-zero length return value from validator")
			}
			err := out[0]
			if !err.IsNil() {
				return fmt.Errorf("validation error: %s\n", err)
			}
		}

	}

	return nil
}

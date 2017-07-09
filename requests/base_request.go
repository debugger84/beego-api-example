package requests

import (
    "strconv"
    "github.com/astaxie/beego/validation"
    "errors"
)

type BaseRequest struct {}

func (this *BaseRequest) exchangeFieldString(key string, values map[string][]string, field *string) error {
    if value, ok := values[key]; ok {
        if len(value) > 0 {
            *field = value[0]
        }
    }

    return nil
}

func (this *BaseRequest) exchangeFieldArrayOfStrings(key string, values map[string][]string, field *[]string) error {
    if value, ok := values[key]; ok {
        if len(value) > 0 {
            *field = value
        }
    }

    return nil
}

func (this *BaseRequest) exchangeFieldInt(key string, values map[string][]string, field *int) error {
    var err error
    if value, ok := values[key]; ok {
        if len(value) > 0 {
            *field, err = strconv.Atoi(value[0])
            if err != nil {
                return err;
            }
        }
    }

    return nil
}

func (this *BaseRequest) returnErrors(valid validation.Validation) []error{
    if valid.HasErrors() {
        errs := make([]error, len(valid.Errors), len(valid.Errors))
        for i, v := range valid.Errors {
            errs[i] = errors.New(v.Key + ": " + v.Message)
        }
        return errs
    }

    return nil
}
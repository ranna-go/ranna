package config

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const envHeader = `
# You can set the following variables if you want to use a custom Docker provider.
# DOCKER_HOST=""
# DOCKER_API_VERSION=""
# DOCKER_CERT_PATH=""
# DOCKER_TLS_VERIFY=""
`

func GenerateDotEnv(w io.Writer) (err error) {
	if _, err = fmt.Fprintln(w, strings.TrimSpace(envHeader)+"\n"); err != nil {
		return err
	}

	extractEnvKV(w, reflect.ValueOf(defaults))
	return nil
}

func extractEnvKV(w io.Writer, val reflect.Value) (err error) {
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return errors.New("value must be of type struct")
	}

	for i := range typ.NumField() {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		if fieldTyp.Type.Kind() == reflect.Struct {
			if err = extractEnvKV(w, val.Field(i)); err != nil {
				return err
			}
			continue
		}

		envKey := fieldTyp.Tag.Get("config")
		if envKey == "" {
			continue
		}

		envVal := fieldVal.Interface()

		_, err = fmt.Fprintf(w, "%s%s=\"%v\"\n", envPrefix, strings.ToUpper(envKey), envVal)
		if err != nil {
			return err
		}
	}

	return nil
}

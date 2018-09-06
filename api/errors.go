package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ResponseError error type for bitbucket 400's
func errorFormatter(err ErrorValue) string {
	return fmt.Sprintf("  context: %s\n   message: %s\n  exception: %s\n", err.Context, err.Message, err.ExceptionName)
}

func authErrors(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	errorJSON := ErrorResponse{}
	err = json.Unmarshal(byt, &errorJSON)
	if err != nil {
		return err
	}
	errString := fmt.Sprintf("\n%s Error\n", resp.Status)
	for _, val := range errorJSON.Errors {
		errString += errorFormatter(val)
	}

	return errors.New(errString)
}

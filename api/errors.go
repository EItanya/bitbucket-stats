package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ResponseError error type for bitbucket 400's
func errorFormatter(err ErrorValue) string {
	if err.Context == nil && err.ExceptionName == nil {
		return fmt.Sprintf("   message: %s\n", err.Message)
	} else if err.Context == nil {
		return fmt.Sprintf("   message: %s\n  exception: %s\n", err.Message, err.ExceptionName)
	} else if err.ExceptionName == nil {
		return fmt.Sprintf("  context: %s\n   message: %s\n", err.Context, err.Message)
	}
	return fmt.Sprintf("  context: %s\n   message: %s\n  exception: %s\n", err.Context, err.Message, err.ExceptionName)
}

func authErrors(resp *http.Response) error {
	if contentTypeArr, ok := resp.Header["Content-Type"]; ok && len(contentTypeArr) > 0 {
		if contentType := contentTypeArr[0]; !strings.Contains(contentType, "application/json") {
			err := []string{
				fmt.Sprintf("%d Error\nURL: %s\n", resp.StatusCode, resp.Request.URL),
				"Data received was not Json therefore not translated \n",
				"This is most likely due to a bad URL",
			}
			return errors.New(strings.Join(err, ""))
		}
		byt, err := ioutil.ReadAll(resp.Body)
		errorJSON := ErrorResponse{}
		err = json.Unmarshal(byt, &errorJSON)
		if err != nil {
			return err
		}
		var errString string
		switch resp.StatusCode {
		case 400:
			errString = "\n400 Error: Bad Request (Check Details Below)\n"
		case 401:
			errString = "\n401 Error: Unauthorized Request (Check Details Below)\n"
		case 403:
			errString = "\n403 Error: Forbidden (Check Details Below)\n"
		case 404:
			errString = "\n404 Error: Not found (Check Details Below)\n"
		default:
			errString = fmt.Sprintf("\n%d Error: Check Details Below\n", resp.StatusCode)
		}
		for _, val := range errorJSON.Errors {
			errString += errorFormatter(val)
		}

		return errors.New(errString)
	}
	panic("Unhandled http error")
}

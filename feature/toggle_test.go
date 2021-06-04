package toggle_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const (
	toggleURL = "http://toggle:8081/v1/toggles"
)

var (
	ctx    = context.Background()
	client = http.DefaultClient

	httpStatus int
	httpBody   []byte
)

type Toggle struct {
	Key         string    `json:"key"`
	IsEnabled   bool      `json:"is_enabled"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetSingleResponse struct {
	Toggle *Toggle `json:"toggle"`
}

type GetAllResponse struct {
	Toggles []*Toggle `json:"toggles"`
}

func TestMain(_ *testing.M) {
	status := godog.TestSuite{
		Name:                "toggle v1alpha1",
		ScenarioInitializer: InitializeScenario,
	}.Run()

	os.Exit(status)
}

func restoreDefaultState(sc *godog.Scenario) {
	err := disableAndDeleteAll()
	checkErr(err)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(restoreDefaultState)

	ctx.Step(`^there are toggles with$`, thereAreTogglesWith)
	ctx.Step(`^the toggle is empty$`, theToggleIsEmpty)
	ctx.Step(`^I create example with body$`, iCreateExampleWithBody)
	ctx.Step(`^I enable toggle with key "([^"]*)"$`, iEnableToggleWithKey)
	ctx.Step(`^I disable toggle with key "([^"]*)"$`, iDisableToggleWithKey)
	ctx.Step(`^I delete toggle with key "([^"]*)"$`, iDeleteToggleWithKey)
	ctx.Step(`^I get all toggles$`, iGetAllToggles)
	ctx.Step(`^I get single toggle with key "([^"]*)"$`, iGetSingleToggleWithKey)
	ctx.Step(`^response status code must be (\d+)$`, responseStatusCodeMustBe)
	ctx.Step(`^response must match json$`, responseMustMatchJSON)
	ctx.Step(`^response single toggle should match$`, responseSingleToggleShouldMatch)
	ctx.Step(`^response toggles should match$`, responseTogglesShouldMatch)
}

func thereAreTogglesWith(requests *godog.Table) error {
	return iCreateExampleWithBody(requests)
}

func theToggleIsEmpty() error {
	err := disableAndDeleteAll()
	return err
}

func iCreateExampleWithBody(requests *godog.Table) error {
	for _, row := range requests.Rows {
		body := strings.NewReader(row.Cells[0].Value)
		if err := callEndpoint(http.MethodPost, toggleURL, body); err != nil {
			return err
		}
	}
	return nil
}

func iEnableToggleWithKey(key string) error {
	return callEndpoint(http.MethodPut, fmt.Sprintf("%s/%s/enable", toggleURL, key), nil)
}

func iDisableToggleWithKey(key string) error {
	return callEndpoint(http.MethodPut, fmt.Sprintf("%s/%s/disable", toggleURL, key), nil)
}

func iDeleteToggleWithKey(key string) error {
	return callEndpoint(http.MethodDelete, fmt.Sprintf("%s/%s", toggleURL, key), nil)
}

func iGetAllToggles() error {
	return callEndpoint(http.MethodGet, toggleURL, nil)
}

func iGetSingleToggleWithKey(key string) error {
	return callEndpoint(http.MethodGet, fmt.Sprintf("%s/%s", toggleURL, key), nil)
}

func responseStatusCodeMustBe(code int) error {
	if httpStatus != code {
		return fmt.Errorf("expected HTTP status code %d, but got %d", code, httpStatus)
	}
	return nil
}

func responseMustMatchJSON(want *godog.DocString) error {
	return deepCompareJSON([]byte(want.Content), httpBody)
}

func responseSingleToggleShouldMatch(want *godog.DocString) error {
	return compareSingleToggleJSON([]byte(want.Content), httpBody)
}

func responseTogglesShouldMatch(want *godog.DocString) error {
	return compareTogglesJSON([]byte(want.Content), httpBody)
}

func getAllToggles() ([]*Toggle, error) {
	if err := callEndpoint(http.MethodGet, toggleURL, nil); err != nil {
		return nil, err
	}

	var resp GetAllResponse
	if err := json.Unmarshal(httpBody, &resp); err != nil {
		return nil, err
	}

	return resp.Toggles, nil
}

func disableAndDeleteAll() error {
	toggles, err := getAllToggles()
	if err != nil {
		return err
	}

	for _, toggle := range toggles {
		if err = callEndpoint(http.MethodPut, fmt.Sprintf("%s/%s/disable", toggleURL, toggle.Key), nil); err != nil {
			return err
		}
		if err = callEndpoint(http.MethodDelete, fmt.Sprintf("%s/%s", toggleURL, toggle.Key), nil); err != nil {
			return err
		}
	}
	return nil
}

func callEndpoint(method, url string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpStatus = resp.StatusCode
	httpBody, err = ioutil.ReadAll(resp.Body)
	return err
}

func deepCompareJSON(want, have []byte) error {
	var expected interface{}
	var actual interface{}

	err := json.Unmarshal(want, &expected)
	if err != nil {
		return err
	}
	err = json.Unmarshal(have, &actual)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func compareSingleToggleJSON(want, have []byte) error {
	var expected GetSingleResponse
	var actual GetSingleResponse

	err := json.Unmarshal(want, &expected)
	if err != nil {
		return err
	}
	err = json.Unmarshal(have, &actual)
	if err != nil {
		return err
	}

	expKey := fmt.Sprintf("%s!!%t!!%s", expected.Toggle.Key, expected.Toggle.IsEnabled, expected.Toggle.Description)
	actKey := fmt.Sprintf("%s!!%t!!%s", actual.Toggle.Key, actual.Toggle.IsEnabled, actual.Toggle.Description)
	if expKey != actKey {
		return fmt.Errorf("expected key: %s, is_enabled: %t, description: %s but not found", expected.Toggle.Key, expected.Toggle.IsEnabled, expected.Toggle.Description)
	}
	return nil
}

func compareTogglesJSON(want, have []byte) error {
	var expected GetAllResponse
	var actual GetAllResponse

	err := json.Unmarshal(want, &expected)
	if err != nil {
		return err
	}
	err = json.Unmarshal(have, &actual)
	if err != nil {
		return err
	}

	if len(expected.Toggles) != len(actual.Toggles) {
		return fmt.Errorf("number of toggle doesn't match, expected %d but got %d", len(expected.Toggles), len(actual.Toggles))
	}

	flag := make(map[string]bool)
	for _, attr := range actual.Toggles {
		keyval := fmt.Sprintf("%s!!%t!!%s", attr.Key, attr.IsEnabled, attr.Description)
		flag[keyval] = true
	}
	for _, attr := range expected.Toggles {
		keyval := fmt.Sprintf("%s!!%t!!%s", attr.Key, attr.IsEnabled, attr.Description)
		if !flag[keyval] {
			return fmt.Errorf("expected key: %s, is_enabled: %t, description: %s but not found", attr.Key, attr.IsEnabled, attr.Description)
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

package extract

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParseJSONForKeyValue(t *testing.T) {
	dataFiles := "./testdata/%s"

	testCases := []struct {
		name       string
		dataFile   string
		root       string
		key        string
		value      string
		wantResult string
	}{
		{
			"very basic JSON (array)",
			"basic-1.json",
			"",
			"name",
			"id",
			`[{"name":"Joe Blogs","value":1},{"name":"Dave Blogs","value":2},{"name":"Trevor Blogs","value":3}]`,
		},
		{
			"basic JSON array with key (object)",
			"basic-2.json",
			"root",
			"name",
			"age",
			`[{"name":"Tom","value":20},{"name":"Dave","value":24},{"name":"Trevor","value":27}]`,
		},
		{
			"morpheus activity response (object)",
			"activity.json",
			"activity",
			"objectType",
			"_id",
			`[{"name":"Instance","value":"2471ed0e-bf19-40af-affa-8d210fe0228c"},{"name":"ComputeServer","value":"ef263cb3-0403-4763-aa65-50fec82ef840"}]`,
		},
		{
			"morpheus users response (object)",
			"users.json",
			"accessTokens",
			"clientId",
			"maskedAccessToken",
			`[{"name":"morph-api","value":"3ae256c-********"},{"name":"morph-cli","value":"10fd4a4-********"}]`,
		},
		{
			"recipes, pick out batters (object)",
			"recipe.json",
			"batter",
			"type",
			"id",
			`[{"name":"Regular","value":"1001"},{"name":"Chocolate","value":"1002"},{"name":"Blueberry","value":"1003"},{"name":"Devil's Food","value":"1004"}]`,
		},
		{
			"recipes, nested, pick out batters (object)",
			"recipe-2.json",
			"batter",
			"type",
			"id",
			`[{"name":"Regular","value":"1001"},{"name":"Chocolate","value":"1002"},{"name":"Blueberry","value":"1003"},{"name":"Devil's Food","value":"1004"}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// make sure test case data is ok
			var checkData interface{}
			data, err := os.ReadFile(fmt.Sprintf(dataFiles, tc.dataFile))
			if err != nil {
				t.Fatalf("problem opening test file: %v", err)
			}

			if err := json.Unmarshal(data, &checkData); err != nil {
				t.Fatalf("the JSON test data appear to be malformed: %v", err)
			}

			if err := json.Unmarshal([]byte(tc.wantResult), &checkData); err != nil {
				t.Fatalf("the JSON expectation appears to be malformed: %v", err)
			}

			// run on the testcase
			gotResult, err := ParseJSONForKeyValue(tc.key, tc.value, data, tc.root)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(gotResult) != tc.wantResult {
				t.Errorf("\ngot `%v`\nwanted `%v`\n", string(gotResult), tc.wantResult)
			}
		})
	}
}

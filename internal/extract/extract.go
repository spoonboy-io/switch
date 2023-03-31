package extract

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// Record is provides the structure for the data we compose
// using a slice of these, in preference to a map, so that we can
// guarantee output order matches order in the original JSON
type Record struct {
	Name  interface{} `json:"name"`
	Value interface{} `json:"value"`
}

// ParseJSONForKeyValue takes a raw JSON document and uses regex and/or
// unmarshalling to extract the provided key and value into a new JSON document
func ParseJSONForKeyValue(key, value string, doc []byte, root ...string) ([]byte, error) {
	var recs []Record
	var data interface{}
	var rootKey string

	// is there a root key?
	if len(root) > 1 {
		return nil, errors.New("cannot pass more than one argument as root")
	}
	if len(root) == 1 {
		rootKey = root[0]
	}

	if rootKey != "" {
		var ok bool
		if doc, ok = obtainRoot(doc, rootKey); !ok {
			return nil, errors.New("cannot locate 'root' key")
		}
	}

	if err := json.Unmarshal(doc, &data); err != nil {
		return nil, err
	}

	recs = walkAtRoot(data, key, value, recs)

	// make json and return
	res, err := json.Marshal(recs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func obtainRoot(data []byte, rootKey string) ([]byte, bool) {
	// helper to extract the array which could be mounted as object key
	// traversing the doc to the key in a unmarshaled map was much more work
	// so we do it in the JSON data with a regex
	regex := fmt.Sprintf(`(?m)"%s":\s*?(\[[[:ascii:]]*\])`, rootKey)
	//regex := fmt.Sprintf(`(?m)"%s":\s*?(\[.*])`, rootKey)
	reg := regexp.MustCompile(regex)
	matches := reg.FindAllStringSubmatch(string(data), 1)

	if len(matches) == 0 {
		// not found
		return []byte{}, false
	}
	return []byte(matches[0][1]), true
}

func walkAtRoot(data interface{}, key, value string, recs []Record) []Record {
	// walk the array, we should able to find map[string]interface{} values
	// we can inspect to extract the desired, key and value
	// we can create them as `name` and `value` fields in a Record and append/return
	switch data.(type) {
	case []interface{}:
		for _, v := range data.([]interface{}) {
			switch i := v.(type) {
			case map[string]interface{}:
				rc := Record{}
				var ok1, ok2 bool
				var k, v interface{}
				if k, ok1 = i[key]; ok1 {
					rc.Name = k
				}
				if v, ok2 = i[value]; ok2 {
					rc.Value = v
				}
				if ok1 && ok2 {
					recs = append(recs, rc)
				}
			}
		}
	}

	return recs
}

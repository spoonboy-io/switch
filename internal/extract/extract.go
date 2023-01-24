package extract

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type Record struct {
	Name  interface{} `json:"name"`
	Value interface{} `json:"value"`
}

// Parses takes raw JSON document and uses regex to extract the
// provided key and value into a new JSON document
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
	regex := fmt.Sprintf(`(?mU)"%s":\s*?(\[[[:ascii:]]*\])`, rootKey)
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
	// we can create them as name value fields in a Record and append/return
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

/*

func SumWithAssertion(payload []byte) (float64, error) {
	var sum float64
	var data interface{}

	err := json.Unmarshal(payload, &data)
	if err != nil {
		return 0, err
	}

	// inspect
	sum += walk(data)

	return math.Ceil(sum*100) / 100, nil
}

// walk traverses the interface and performs the addition, it calls itself recursively for nested aggregate types
func walk(data interface{}) float64 {
	var sum float64
	switch i := data.(type) {
	case []interface{}:
		for _, v := range data.([]interface{}) {
			switch i := v.(type) {
			case int:
				sum += float64(i)
			case float64:
				sum += i
			case []interface{}, map[string]interface{}:
				r := walk(v)
				sum += r
			default:
				fmt.Println("Not numeric skipping:", v, reflect.TypeOf(v))
			}
		}
	case map[string]interface{}:
		for _, v := range data.(map[string]interface{}) {
			switch i := v.(type) {
			case int:
				sum += float64(i)
			case float64:
				sum += i
			case []interface{}, map[string]interface{}:
				r := walk(v)
				sum += r
			default:
				fmt.Println("Not numeric skipping:", v, reflect.TypeOf(v))
			}
		}
	case int:
		sum += float64(i)
	case float64:
		sum += i
	default:
		fmt.Println("Not numeric skipping:", data, reflect.TypeOf(data))
	}
	return sum
}
*/

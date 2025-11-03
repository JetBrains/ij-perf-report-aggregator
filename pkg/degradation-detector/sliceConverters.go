package degradation_detector

import (
	"fmt"
	"strconv"
)

func SliceToSliceInt64(strings []any) ([]int64, error) {
	result := make([]int64, 0, len(strings))
	for _, s := range strings {
		switch v := s.(type) {
		case string:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			result = append(result, i)
		case int64:
			result = append(result, v)
		case int:
			result = append(result, int64(v))
		case float64:
			result = append(result, int64(v))
		default:
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
	}
	return result, nil
}

func SliceToSliceOfInt(strings []any) ([]int, error) {
	result := make([]int, 0, len(strings))
	for _, s := range strings {
		switch v := s.(type) {
		case string:
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			result = append(result, i)
		case int64:
			result = append(result, int(v))
		case int:
			result = append(result, v)
		case float64:
			result = append(result, int(v))
		default:
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
	}
	return result, nil
}

func SliceToSliceOfString(slice []any) ([]string, error) {
	strings := make([]string, 0, len(slice))
	for _, elem := range slice {
		switch v := elem.(type) {
		case string:
			strings = append(strings, v)
		case int:
			strings = append(strings, strconv.Itoa(v))
		case int64:
			strings = append(strings, strconv.FormatInt(v, 10))
		case float64:
			strings = append(strings, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			strings = append(strings, strconv.FormatBool(v))
		case nil:
			strings = append(strings, "")
		// Add other types as needed
		default:
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
	}
	return strings, nil
}

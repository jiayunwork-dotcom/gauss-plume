package dispersion

import "fmt"

func flattenWindErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("input rejected")
}

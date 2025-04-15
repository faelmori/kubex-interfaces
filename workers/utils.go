package workers

import "fmt"

func validateWorkerLimit(v any) error {
	if i, ok := v.(int); !ok || i <= 0 {
		return fmt.Errorf("workerLimit deve ser um inteiro maior que 0")
	}
	return nil
}

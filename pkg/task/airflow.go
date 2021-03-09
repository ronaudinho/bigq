// This is a template Airflow related operations
package task

// Airflow holds Airflow related information
type Airflow struct{}

// NewAirflow creates new instance for Argo
func NewAirflow() *Airflow {
	return &Airflow{}
}

// Process processes tasks Airflow Queueing Logic
func (a *Airflow) Process(payload string) (string, error) {
	// TODO queue logic here
	return payload, nil
}

// Send is a dummy task to satisfy machinery.Task requirement
func (a *Airflow) Send() error {
	return nil
}

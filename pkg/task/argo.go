// This is a template Argo related operations
package task

// Argo holds Argo related information
type Argo struct{}

// NewArgo creates new instance of Argo
func NewArgo() *Argo {
	return &Argo{}
}

// Process processes tasks Argo Queueing Logic
func (a *Argo) Process(payload string) (string, error) {
	// TODO queue logic here
	return payload, nil
}

// Send is a dummy task to satisfy machinery.Task requirement
func (a *Argo) Send() error {
	return nil
}

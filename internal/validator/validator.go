package validator

type ValidationError struct {
	Code     string
	Message  string
	Severity string
	Path     []string
}

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	for _, err := range r.Errors {
		if err.Severity == "error" {
			return true
		}
	}
	return false
}

func (r *ValidationResult) AddError(err ValidationError) {
	r.Errors = append(r.Errors, err)
}

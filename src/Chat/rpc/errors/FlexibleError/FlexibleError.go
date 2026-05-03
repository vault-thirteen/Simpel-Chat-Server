package fe

type FlexibleError struct {
	err error
}

func NewFlexibleError(err error) *FlexibleError {
	if err == nil {
		return nil
	}

	return &FlexibleError{err: err}
}

func (fe *FlexibleError) Value() any {
	if fe == nil {
		return nil
	}

	if fe.err == nil {
		return nil
	}

	return fe.err.Error()
}

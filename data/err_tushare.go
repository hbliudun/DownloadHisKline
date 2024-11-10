package data

type ErrTushare struct {
	error
	err string
}

func (e ErrTushare) Error() string {
	return e.err
}

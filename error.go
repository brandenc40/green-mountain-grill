package gmg

const _unreachableErr = "grill is unreachable: "

type GrillUnreachableErr struct {
	Err error
}

func (g GrillUnreachableErr) Error() string {
	if g.Err == nil {
		return _unreachableErr
	}
	return _unreachableErr + g.Err.Error()
}

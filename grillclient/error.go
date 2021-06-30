package grillclient

const _unreachableErr = "grill is unreachable: "

type GrillUnreachableErr struct {
	err error
}

func (g GrillUnreachableErr) Error() string {
	return _unreachableErr + g.err.Error()
}

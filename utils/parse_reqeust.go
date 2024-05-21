package utils

type ParseOpts struct {
	Body   interface{}
	Header *map[string]string
}

type ParseOpt func(*ParseOpts) error

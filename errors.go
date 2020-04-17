package helpers

type NotImplementedError string

func (e NotImplementedError) Error() string {
	return string(e)
}

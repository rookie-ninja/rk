package rk_common

func NewFileOperationError(msg string) error {
	return &FileOperationError{
		msg: msg,
	}
}

type FileOperationError struct {
	msg string
}

func (err *FileOperationError) Error() string {
	return err.msg
}

func NewGithubClientError(msg string) error {
	return &GithubClientError{
		msg: msg,
	}
}

type GithubClientError struct {
	msg string
}

func (err *GithubClientError) Error() string {
	return err.msg
}

func NewInvalidParamError(msg string) error {
	return &InvalidParamError{
		msg: msg,
	}
}

type InvalidParamError struct {
	msg string
}

func (err *InvalidParamError) Error() string {
	return err.msg
}

func NewMarshalError(msg string) error {
	return &MarshalError{
		msg: msg,
	}
}

type MarshalError struct {
	msg string
}

func (err *MarshalError) Error() string {
	return err.msg
}

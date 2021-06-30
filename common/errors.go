// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package common

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

func NewCommandError(msg string) error {
	return &CommandError{
		msg: msg,
	}
}

type CommandError struct {
	msg string
}

func (err *CommandError) Error() string {
	return err.msg
}

func NewScriptError(msg string) error {
	return &ScriptError{
		msg: msg,
	}
}

type ScriptError struct {
	msg string
}

func (err *ScriptError) Error() string {
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

func NewUnMarshalError(msg string) error {
	return &UnMarshalError{
		msg: msg,
	}
}

type UnMarshalError struct {
	msg string
}

func (err *UnMarshalError) Error() string {
	return err.msg
}

func NewInvalidEnvError(msg string) error {
	return &GithubClientError{
		msg: msg,
	}
}

type InvalidEnvError struct {
	msg string
}

func (err *InvalidEnvError) Error() string {
	return err.msg
}

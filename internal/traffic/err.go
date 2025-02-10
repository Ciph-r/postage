package traffic

import "errors"

var (
	// ErrNotFound is when the specified resource does not exist.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is when the specified resource already exist.
	ErrAlreadyExists = errors.New("already exists")

	// ErrDisconnected is when the specified resource disconnected.
	ErrDisconnected = errors.New("disconnected")
)

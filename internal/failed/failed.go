package failed

import "github.com/yanun0323/errors"

var ErrFailed = errors.New("failed")

type Failed struct{}

func (f Failed) Error() error {
	return errors.Wrap(f.ErrorWithFormatWrap()).With(
		"code", 400,
		"message", "failed",
		"struct", struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    400,
			Message: "failed",
		},
	)
}

func (Failed) ErrorWithWrap() error {
	return errors.Wrap(ErrFailed)
}

func (f Failed) ErrorDelegate() error {
	return f.ErrorWithWrap()
}

func (f Failed) ErrorWithFormatWrap() error {
	// return errors.Wrap(f.ErrorDelegate(), "wrap")
	return errors.Errorf("wrap: %w", f.ErrorDelegate())
}

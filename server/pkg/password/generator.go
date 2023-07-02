package password

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Errors
var (
	ZeroOptionsErr = errors.New("unable to produce the password with no options for characters")
	LengthErr      = errors.New("unable to produce the password with number of options bigger than the length")
)

const (
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits       = "0123456789"
	Symbols      = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

type Options struct {
	Length       int
	LowerLetters bool
	UpperLetters bool
	Digits       bool
	Symbols      bool
}

type optionSizes struct {
	LowerLetters int
	UpperLetters int
	Digits       int
	Symbols      int
}

func Generate(opts Options) (string, error) {
	var s = ""

	err := validate(opts)
	if err != nil {
		return "", err
	}

	sizes, err := computeOptionLengths(opts)
	if err != nil {
		return "", err
	}

	for i := 0; i < int(sizes.LowerLetters); i++ {
		c, err := getNewCharacter(LowerLetters)
		if err != nil {
			return "", err
		}

		i, err := getNewIndex(len(s) + 1)
		if err != nil {
			return "", err
		}
		s = s[:i] + c + s[i:]
	}

	for i := 0; i < int(sizes.UpperLetters); i++ {
		c, err := getNewCharacter(UpperLetters)
		if err != nil {
			return "", err
		}

		i, err := getNewIndex(len(s) + 1)
		if err != nil {
			return "", err
		}
		s = s[:i] + c + s[i:]
	}

	for i := 0; i < int(sizes.Digits); i++ {
		c, err := getNewCharacter(Digits)
		if err != nil {
			return "", err
		}

		i, err := getNewIndex(len(s) + 1)
		if err != nil {
			return "", err
		}
		s = s[:i] + c + s[i:]
	}

	for i := 0; i < int(sizes.Symbols); i++ {
		c, err := getNewCharacter(Symbols)
		if err != nil {
			return "", err
		}

		i, err := getNewIndex(len(s) + 1)
		if err != nil {
			return "", err
		}
		s = s[:i] + c + s[i:]
	}

	return s, nil
}

func getNewIndex(length int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), err
}

func getNewCharacter(charactersSet string) (string, error) {
	length := len(charactersSet)
	n, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
	if err != nil {
		return "", err
	}
	var i int = int(n.Int64())
	return string(charactersSet[i]), err
}

// TODO: refactor this whole function
func computeOptionLengths(opts Options) (optionSizes, error) {
	rest := opts.Length
	optsNum := 0

	if opts.LowerLetters {
		optsNum++
		rest--
	}
	if opts.UpperLetters {
		optsNum++
		rest--
	}
	if opts.Digits {
		optsNum++
		rest--
	}
	if opts.Symbols {
		optsNum++
		rest--
	}

	distributedPool := make([]int, optsNum)
	for i := range distributedPool {
		distributedPool[i] = 1
	}

	for i := 0; i < rest; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(distributedPool))))
		if err != nil {
			return optionSizes{}, err
		}
		distributedPool[int(randomIndex.Int64())] += 1
	}

	sizes := optionSizes{}
	i := 0

	if opts.LowerLetters {
		sizes.LowerLetters = distributedPool[i]
		i++
	}
	if opts.UpperLetters {
		sizes.UpperLetters = distributedPool[i]
		i++
	}
	if opts.Digits {
		sizes.Digits = distributedPool[i]
		i++
	}
	if opts.Symbols {
		sizes.Symbols = distributedPool[i]
		i++
	}

	return sizes, nil
}

func validate(opts Options) error {
	reqOpts := 0
	if opts.LowerLetters {
		reqOpts++
	}
	if opts.UpperLetters {
		reqOpts++
	}
	if opts.Digits {
		reqOpts++
	}
	if opts.Symbols {
		reqOpts++
	}
	if reqOpts == 0 {
		return ZeroOptionsErr
	}
	if reqOpts > opts.Length {
		return LengthErr
	}
	return nil
}

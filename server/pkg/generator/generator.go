package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"

	"github.com/ChristianSassine/password-manager/server/pkg/utils"
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
	Length       int  `json:"length"`
	LowerLetters bool `json:"lowerLetters"`
	UpperLetters bool `json:"upperLetters"`
	Digits       bool `json:"digits"`
	Symbols      bool `json:"symbols"`
}

type optionSizes struct {
	LowerLetters int
	UpperLetters int
	Digits       int
	Symbols      int
}

func Generate(opts Options) (string, error) {
	type poolOption struct {
		count      int
		characters string
	}

	var s strings.Builder

	err := validate(opts)
	if err != nil {
		return "", err
	}

	sizes, err := computeOptionLengths(opts)
	if err != nil {
		return "", err
	}

	var poolOpts []poolOption = []poolOption{}

	if sizes.LowerLetters > 0 {
		poolOpts = append(poolOpts, poolOption{count: sizes.LowerLetters, characters: LowerLetters})
	}

	if sizes.UpperLetters > 0 {
		poolOpts = append(poolOpts, poolOption{count: sizes.UpperLetters, characters: UpperLetters})
	}

	if sizes.Digits > 0 {
		poolOpts = append(poolOpts, poolOption{count: sizes.Digits, characters: Digits})
	}

	if sizes.Symbols > 0 {
		poolOpts = append(poolOpts, poolOption{count: sizes.Symbols, characters: Symbols})
	}

	for len(poolOpts) > 0 {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(poolOpts))))
		if err != nil {
			return "", err
		}
		i := int(r.Int64())

		var characters = poolOpts[i].characters
		r, err = rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return "", err
		}
		j := int(r.Int64())
		s.WriteByte(characters[j])
		poolOpts[i].count--
		if poolOpts[i].count == 0 {
			poolOpts = utils.DeleteIndex(poolOpts, i) // We'll do this at max 4 times. It can be considered constant.
		}
	}

	return s.String(), nil
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

func computeOptionLengths(opts Options) (optionSizes, error) {
	remaining := opts.Length
	optsNum := 0

	if opts.LowerLetters {
		optsNum++
		remaining--
	}
	if opts.UpperLetters {
		optsNum++
		remaining--
	}
	if opts.Digits {
		optsNum++
		remaining--
	}
	if opts.Symbols {
		optsNum++
		remaining--
	}

	distributedPool := make([]int, optsNum)
	for i := range distributedPool {
		distributedPool[i] = 1
	}

	for i := 0; i < remaining; i++ {
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
	optsNum := 0
	if opts.LowerLetters {
		optsNum++
	}
	if opts.UpperLetters {
		optsNum++
	}
	if opts.Digits {
		optsNum++
	}
	if opts.Symbols {
		optsNum++
	}
	if optsNum == 0 {
		return ZeroOptionsErr
	}
	if optsNum > opts.Length {
		return LengthErr
	}
	return nil
}

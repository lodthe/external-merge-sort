package main

import (
	"bufio"
	"math/rand"
)

type generator struct {
	cfg *config
	rnd *rand.Rand
}

func newGenerator(cfg *config, rnd *rand.Rand) *generator {
	return &generator{
		cfg: cfg,
		rnd: rnd,
	}
}

func (g *generator) generate() error {
	w := bufio.NewWriter(g.cfg.output)

	equalPref := g.token(g.cfg.equalPrefixLength)

	for i := 0; i < g.cfg.count; i++ {
		_, err := w.Write(equalPref)
		if err != nil {
			return err
		}

		length := g.cfg.minLength + g.rnd.Intn(g.cfg.maxLength-g.cfg.minLength+1)
		suffix := g.token(length - g.cfg.equalPrefixLength)
		_, err = w.Write(suffix)
		if err != nil {
			return err
		}

		err = w.WriteByte(g.cfg.delimiter)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

func (g *generator) token(n int) []byte {
	token := make([]byte, n)
	for i := range token {
		token[i] = g.byte()
	}

	return token
}

func (g *generator) byte() byte {
	return g.cfg.alphabet[g.rnd.Intn(len(g.cfg.alphabet))]
}

package generator

import (
	rp "github.com/vault-thirteen/auxie/rpofs"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Generator struct {
	ridg *rp.Generator
	vcg  *rp.Generator
	tg   *rp.Generator
}

func NewGenerator(stn *settings.OtherChatSettings) (g *Generator, err error) {
	g = new(Generator)

	g.ridg, err = rp.NewGenerator(int(stn.RequestIdLength), helper.ListRandomSymbolsForIds())
	if err != nil {
		return nil, err
	}

	g.vcg, err = rp.NewGenerator(int(stn.VerificationCodeLength), helper.ListRandomSymbolsForIds())
	if err != nil {
		return nil, err
	}

	g.tg, err = rp.NewGenerator(int(stn.TokenLength), helper.ListRandomSymbolsForTokens())
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Generator) RIDG() *rp.Generator {
	return g.ridg
}
func (g *Generator) VCG() *rp.Generator {
	return g.vcg
}
func (g *Generator) TG() *rp.Generator {
	return g.tg
}

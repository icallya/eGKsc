package smartcard

import (
	"crypto/x509"
	"fmt"
	"log/slog"

	"main/src/brainpool"
)

type CardType string

const (
	CardTypeUnknown CardType = "Unknown"
	CardTypeEGK     CardType = "egk"
	CardTypeHBA     CardType = "hba"
	CardTypeSMCB    CardType = "smc-b"
	CardTypeGSMCK   CardType = "gsmc-k"
	CardTypeGSMCKT  CardType = "gsmc-kt"
)

func (c *Card) SelectMF() error {
	_, err := c.TransmitAPDU(CommandAPDU{0x00, 0xA4, 0x00, 0x00, 0x02, 0x3F, 0x00})
	return err
}

const (
	ShortFileIdentifierEFDIR byte = 0x1E
)

type DedicatedFile struct {
	Name                  string
	ApplicationIdentifier []byte
}

type MasterFile struct {
	Name                  string
	ApplicationIdentifier []byte
	CardType              CardType
}

type ElementaryFile struct {
	FileIdentifier  [2]byte
	ShortIdentifier byte
}

var DF_HCA = DedicatedFile{"DF.HCA", []byte{0xD2, 0x76, 0x00, 0x00, 0x01, 0x02}}
var DF_ESIGN = DedicatedFile{"DF.ESIGN", []byte{0xA0, 0x00, 0x00, 0x01, 0x67, 0x45, 0x53, 0x49, 0x47, 0x4E}}
var DF_QES = DedicatedFile{"DF.QES", []byte{0xD2, 0x76, 0x00, 0x00, 0x66, 0x01}}
var DF_NFD = DedicatedFile{"DF.NFD", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x07}}
var DF_DPE = DedicatedFile{"DF.DPE", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x08}}
var DF_GDD = DedicatedFile{"DF.GDD", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x0A}}
var DF_OSE = DedicatedFile{"DF.OSE", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x0B}}
var DF_AMTS = DedicatedFile{"DF.AMTS", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x0C}}
var DF_HPA = DedicatedFile{"DF.HPA", []byte{0xD2, 0x76, 0x00, 0x01, 0x46, 0x02}}
var DF_CIA_QES = DedicatedFile{"DF.CIA.QES", []byte{0xD2, 0x76, 0x00, 0x00, 0x66, 0x01}}
var DF_AUTO = DedicatedFile{"DF.AUTO", []byte{0xD2, 0x76, 0x00, 0x01, 0x46, 0x03}}
var DF_KT = DedicatedFile{"DF.KT", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x00}}

var KnownMFs []MasterFile
var KnownDFs []DedicatedFile

func init() {
	KnownMFs = []MasterFile{
		{"MF", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x80, 0x00}, CardTypeEGK},
		{"MF", []byte{0xD2, 0x76, 0x00, 0x01, 0x46, 0x01}, CardTypeHBA},
		{"MF", []byte{0xD2, 0x76, 0x00, 0x01, 0x46, 0x06}, CardTypeSMCB},
		{"MF", []byte{0xD2, 0x76, 0x00, 0x01, 0x44, 0x80, 0x03}, CardTypeGSMCKT},
	}
	KnownDFs = []DedicatedFile{
		DF_HCA,
		DF_ESIGN,
		DF_QES,
		DF_NFD,
		DF_DPE,
		DF_GDD,
		DF_OSE,
		DF_AMTS,
		DF_HPA,
		DF_CIA_QES,
		DF_AUTO,
		DF_KT,
	}
	for _, mf := range KnownMFs {
		KnownDFs = append(KnownDFs, DedicatedFile{mf.Name, mf.ApplicationIdentifier})
	}
}

var EF_C_CH_AUT_E256 = ElementaryFile{[2]byte{0xC5, 0x04}, 0x04}

func (c *Card) SelectDF(df DedicatedFile) error {
	apdu := append(CommandAPDU{0x00, 0xA4, 0x04, 0x0C, 0x0A}, df.ApplicationIdentifier...)
	response, err := c.TransmitAPDU(apdu)
	slog.Info("SelectDF", "response", response)
	return err
}

func (c *Card) ReadTransparentEF(ef ElementaryFile) ([]byte, error) {
	apdu := CommandAPDU{0x00, 0xB0, 0x80 + ef.ShortIdentifier, 0x00, 0x00, 0x00, 0x00}
	response, err := c.TransmitAPDU(apdu)
	if err != nil {
		return nil, err
	}
	if response.SW1() != 0x90 || response.SW2() != 0x00 {
		return nil, fmt.Errorf("unexpected response: %02x%02x", response.SW1(), response.SW2())
	}
	return response.Data(), nil

}

func (c *Card) ReadCertificate(ef ElementaryFile) (*x509.Certificate, error) {
	data, err := c.ReadTransparentEF(ef)
	if err != nil {
		return nil, err
	}
	return brainpool.ParseCertificate(data)
}

func DataSequence(data []byte) func() (byte, []byte) {
	pos := 0
	return func() (byte, []byte) {
		if pos >= len(data) {
			return 0, nil
		}
		tag := data[pos]
		pos++
		length := int(data[pos])
		pos++
		value := data[pos : pos+length]
		pos += length
		return tag, value
	}
}

package smartcard

import (
	"bytes"
	"fmt"
	"log/slog"
)

func resolveDF(aid []byte) (*DedicatedFile, bool) {
	for _, df := range KnownDFs {
		if bytes.Equal(aid, df.ApplicationIdentifier) {
			return &df, true
		}
	}
	return nil, false
}

type EFDIR struct {
	CardType     CardType
	Applications []DedicatedFile
}

func (c *Card) EFDIR() (*EFDIR, error) {
	apdu := CommandAPDU{0x00, 0xb2, 0x0, (ShortFileIdentifierEFDIR << 3) + 4, 0x00}
	var applications []DedicatedFile = make([]DedicatedFile, 0)
	var cardType CardType = CardTypeUnknown
	for {
		apdu[2]++
		response, err := c.TransmitAPDU(apdu)
		if err != nil {
			return nil, err
		}

		if response.SW1() == 0x6A && response.SW2() == 0x83 {
			break
		} else if response.SW1() != 0x90 || response.SW2() != 0x00 {
			return nil, fmt.Errorf("unexpected response: %02x%02x", response.SW1(), response.SW2())
		}

		recSeq := DataSequence(response.Data())
		recTag, recVal := recSeq()
		if recTag != 0x61 {
			return nil, fmt.Errorf("unexpected record tag: %02x", recTag)
		}

		aidSeq := DataSequence(recVal)
		aidTag, aidVal := aidSeq()
		if aidTag != 0x4F {
			return nil, fmt.Errorf("unexpected AID tag: %02x", aidTag)
		}
		if df, ok := resolveDF(aidVal); ok {
			for _, mf := range KnownMFs {
				if bytes.Equal(aidVal, mf.ApplicationIdentifier) {
					cardType = mf.CardType
				}
			}
			applications = append(applications, *df)
		} else {
			slog.Error("unknown AID", "aid", fmt.Sprintf("%x", aidVal))
		}
	}

	return &EFDIR{
		Applications: applications,
		CardType:     cardType,
	}, nil
}

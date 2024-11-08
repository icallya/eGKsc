package smartcard

import (
	"fmt"
	"log"
)

type ObjSysProductID struct {
	Vendor  string
	Product string
	Version string
}

type EFVersion2 struct {
	FormatVersion                      string
	ObjSysVersion                      string
	ObjSysProductID                    ObjSysProductID
	EFGOFormatVersion                  string
	EFATRFormatVersion                 string
	EFKeyInfoFormatVersion             string
	EFEnvironmentSettingsFormatVersion string
	EFLoggingFormatVersion             string
}

func bytesToVersion(data []byte) string {
	return fmt.Sprintf("%d.%d.%d", data[0], data[1], data[2])
}

func (c *Card) EFVersion2() (*EFVersion2, error) {
	response, err := c.TransmitAPDU(CommandAPDU{0x00, 0xB0, 0x80 + 0x11, 0x00, 0x00})
	if err != nil {
		log.Fatal(err)
	}

	v := &EFVersion2{}

	data := response.Data()

	rootSeq := DataSequence(data)
	eftag, efval := rootSeq()
	if err != nil {
		return nil, err
	}
	if eftag != 0xEF {
		return nil, fmt.Errorf("expected 0xEF, got %02x", eftag)
	}

	efSeq := DataSequence(efval)
	for tag, val := efSeq(); tag != 0; tag, val = efSeq() {
		switch tag {
		case 0xC0:
			v.FormatVersion = bytesToVersion(val)
		case 0xC1:
			v.ObjSysVersion = bytesToVersion(val)
		case 0xC2:
			v.ObjSysProductID = ObjSysProductID{
				Vendor:  string(val[0:5]),
				Product: string(val[5:13]),
				Version: bytesToVersion(val[13:16]),
			}
		case 0xC4:
			v.EFGOFormatVersion = bytesToVersion(val)
		case 0xC5:
			v.EFATRFormatVersion = bytesToVersion(val)
		case 0xC6:
			v.EFKeyInfoFormatVersion = bytesToVersion(val)
		case 0xC3:
			v.EFEnvironmentSettingsFormatVersion = bytesToVersion(val)
		case 0xC7:
			v.EFLoggingFormatVersion = bytesToVersion(val)
		}
	}

	return v, nil
}

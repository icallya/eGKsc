/*
Package smartcard implements a portable high-level API for communicating with smart cards.

Example:

	ctx, err := smartcard.EstablishContext()
	// handle error, if any
	defer ctx.Release()

	reader, err := ctx.WaitForCardPresent()
	// handle error, if any

	card, err := reader.Connect()
	// handle error, if any
	defer card.Disconnect()

	fmt.Printf("Card ATR: %s\n", card.ATR())
	command := SelectCommand(0xa0, 0x00, 0x00, 0x00, 0x62, 0x03, 0x01, 0xc, 0x01, 0x01)
	response, err := card.TransmitAPDU(command)
	// handle error, if any
	fmt.Printf("Response: %s\n", response)
*/
package smartcard

import (
	"bytes"
	"fmt"

	"main/src/smartcard/pcsc"
)

const (
	// Scope
	SCOPE_USER     = pcsc.CARD_SCOPE_USER
	SCOPE_TERMINAL = pcsc.CARD_SCOPE_TERMINAL
	SCOPE_SYSTEM   = pcsc.CARD_SCOPE_SYSTEM
)

type ATR []byte

// Return string form of ATR.
func (atr ATR) String() string {
	var buffer bytes.Buffer
	for _, b := range atr {
		buffer.WriteString(fmt.Sprintf("%02x", b))
	}
	return buffer.String()
}

// Transmit command APDU to the card and return response.
func (c *Card) TransmitAPDU(cmd CommandAPDU) (ResponseAPDU, error) {
	bytes, err := c.Transmit(cmd)
	if err != nil {
		return nil, err
	}
	r, err := Response(bytes)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ISO7816-4 command APDU.
type CommandAPDU []byte

// Return string form of APDU.
func (cmd CommandAPDU) String() string {
	apdu := ([]byte)(cmd)
	buffer := new(bytes.Buffer)
	buffer.WriteString(fmt.Sprintf("%02X %02X %02X %02X", apdu[0], apdu[1],
		apdu[2], apdu[3]))
	if len(apdu) >= 5 {
		buffer.WriteString(fmt.Sprintf(" %02X", apdu[4]))
		if len(apdu) >= 6 {
			if len(apdu) == int(apdu[4]+5) {
				buffer.WriteString(fmt.Sprintf(" %X", apdu[5:]))
			} else {
				buffer.WriteString(fmt.Sprintf(" %X %02X", apdu[5:len(apdu)-1],
					apdu[len(apdu)-1]))
			}
		}
	}
	return buffer.String()
}

// ISO7816-4 response APDU.
type ResponseAPDU []byte

func Response(bytes []byte) (ResponseAPDU, error) {
	if len(bytes) < 2 {
		return nil, fmt.Errorf("Invalid response apdu size: %d", len(bytes))
	}
	return ResponseAPDU(bytes), nil
}

// Return 16-bit status word.
func (r ResponseAPDU) SW() uint16 {
	return uint16(r[len(r)-2])<<8 | uint16(r[len(r)-1])
}

// Return SW1
func (r ResponseAPDU) SW1() uint8 {
	return r[len(r)-2]
}

// Return SW2
func (r ResponseAPDU) SW2() uint8 {
	return r[len(r)-1]
}

// Return data part of response
func (r ResponseAPDU) Data() []byte {
	if len(r) <= 2 {
		return nil
	}
	return r[:len(r)-2]
}

// Return string form of APDU.
func (r ResponseAPDU) String() string {
	var bytes []byte = r
	if len(r) <= 2 {
		return fmt.Sprintf("%X", bytes)
	}
	return fmt.Sprintf("%X %X", bytes[:len(bytes)-2], bytes[len(bytes)-2:])
}

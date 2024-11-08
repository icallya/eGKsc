package smartcard

import (
	"fmt"
	"testing"
)

func TestInfo(t *testing.T) {
	fmt.Println("\n===================")
	fmt.Println("High Level API Test")
	fmt.Printf("===================\n\n")
}

func TestEstablishReleaseUserContext(t *testing.T) {
	fmt.Println("------------------------------")
	fmt.Println("Test establish/release User Context")
	fmt.Printf("------------------------------\n\n")
	ctx, err := EstablishContext(SCOPE_USER)
	if err != nil {
		t.Error(err)
		return
	}
	err = ctx.Release()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("OK\n\n")
}

func TestEstablishReleaseSystemContext(t *testing.T) {
	fmt.Println("------------------------------")
	fmt.Println("Test establish/release System Context")
	fmt.Printf("------------------------------\n\n")
	ctx, err := EstablishContext(SCOPE_SYSTEM)
	if err != nil {
		t.Error(err)
		return
	}
	err = ctx.Release()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("OK\n\n")
}

func TestListReaders(t *testing.T) {
	fmt.Println("-----------------")
	fmt.Println("Test list readers")
	fmt.Printf("-----------------\n\n")
	ctx, err := EstablishContext()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Release()
	readers, err := ctx.ListReaders()
	if err != nil {
		t.Error(err)
		return
	}
	for _, reader := range readers {
		fmt.Println(reader.Name())
		fmt.Printf("- Card present: %t\n\n", reader.IsCardPresent())
	}
}

func TestListReadersWithCard(t *testing.T) {
	fmt.Println("---------------------------")
	fmt.Println("Test list readers with card")
	fmt.Printf("---------------------------\n\n")
	ctx, err := EstablishContext()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Release()
	readers, err := ctx.ListReadersWithCard()
	if err != nil {
		t.Error(err)
		return
	}
	for _, reader := range readers {
		fmt.Println(reader.Name())
		fmt.Printf("- Card present: %t\n\n", reader.IsCardPresent())
	}
}

func TestWaitForCardPresentRemoved(t *testing.T) {
	fmt.Println("----------------------------------------------------")
	fmt.Println("Test wait for card present / wait until card removed")
	fmt.Printf("----------------------------------------------------\n\n")
	ctx, err := EstablishContext()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Release()
	fmt.Printf("Insert card now...")
	reader, err := ctx.WaitForCardPresent()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("\n\n%s\n", reader.Name())
	fmt.Printf("- Card present: %t\n\n", reader.IsCardPresent())
	fmt.Printf("Remove card now...")
	reader.WaitUntilCardRemoved()
	fmt.Printf("\n\nCard was removed\n\n")
}

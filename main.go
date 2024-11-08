package main

import (
  "fmt"
  "log"
  "log/slog"
  "main/src/smartcard"
)

func main() {
  ctx, err := smartcard.EstablishContext()
  if (err == nil) {
    defer ctx.Release()
    reader, err := ctx.WaitForCardPresent()
    if (err == nil) {
      slog.Info("Reader created", "reader", reader.Name())
      card, err := reader.Connect()
      if (err == nil) {
        defer card.Disconnect()
        slog.Info("Card connected", "atr", card.ATR().String())
        fmt.Printf("\n\nATR: %s\n\n", card.ATR())
        efv, err := card.EFVersion2()
        if (err == nil) {
          slog.Info("EF Version 2", "efv", efv)
          efdir, err := card.EFDIR()
          if (err == nil) {
            slog.Info("EF.DIR", "EF.DIR", efdir)
            err = card.SelectDF(smartcard.DF_ESIGN)
            if (err == nil) {
              cert, err := card.ReadCertificate(smartcard.EF_C_CH_AUT_E256)
              if (err == nil) {
                slog.Info("EF.C.CH.AUT.E256", "subject", cert.Subject.CommonName, "notAfter", cert.NotAfter)
              } else {
                log.Fatal(err)
              }
            } else {
              log.Fatal(err)
            }
          } else {
            log.Fatal(err)
          }
        } else {
          log.Fatal(err)
        }
      } else {
        log.Fatal(err)
      }
    } else {
      log.Fatal(err)
    }
  } else {
    log.Fatal(err)
  }
}

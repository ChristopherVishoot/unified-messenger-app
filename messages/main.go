package main

import (
    "context"
    "fmt"
    "os"

    "github.com/mdp/qrterminal/v3"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    waLog "go.mau.fi/whatsmeow/util/log"
    _ "github.com/lib/pq"
)

func main() {
    // Connect to Postgres for session storage
    dbLog := waLog.Stdout("Database", "DEBUG", true)
    container, err := sqlstore.New("postgres", "postgres://user:pass@localhost:5432/bridge_db?sslmode=disable", dbLog)
    if err != nil {
        panic(err)
    }

    deviceStore, err := container.GetFirstDevice()
    if err != nil {
        panic(err)
    }

    clientLog := waLog.Stdout("Client", "DEBUG", true)
    client := whatsmeow.NewClient(deviceStore, clientLog)

    // Handle incoming messages
    client.AddEventHandler(func(evt interface{}) {
        // Process messages here → save to Postgres
        fmt.Printf("Event: %+v\n", evt)
    })

    if client.Store.ID == nil {
        // No session yet — show QR code
        qrChan, _ := client.GetQRChannel(context.Background())
        err = client.Connect()
        if err != nil {
            panic(err)
        }
        for evt := range qrChan {
            if evt.Event == "code" {
                // Print QR to terminal
                qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
            } else {
                fmt.Println("Login event:", evt.Event)
            }
        }
    } else {
        // Already logged in — just connect
        err = client.Connect()
        if err != nil {
            panic(err)
        }
    }

    // Keep running
    select {}
}
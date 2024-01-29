package payment

import (
	"RajaBot/database"
	"RajaBot/tools"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

func handleCallback(w http.ResponseWriter, r *http.Request) {
	decoder := schema.NewDecoder()
	cbData := &CallbackParams{}

	err := decoder.Decode(cbData, r.URL.Query())
	if err != nil {
		log.Printf("Payment -> Error in GET parameters: %v\n", err)
	}
	if cbData.TrackId == 0 || cbData.OrderId == "" {
		log.Println("Payment -> Error in callback data")
		io.WriteString(w, BadCbMsg)
		return
	}
	// if cbData.Success == 0 {
	// 	log.Println("Payment -> Success == 0")
	// 	io.WriteString(w, OrderFailed)
	// 	return
	// }

	paym := database.GetPayment(cbData.OrderId)
	if paym == nil || paym.IsDone {
		io.WriteString(w, OrderNotFound)
		return
	}

	// verify
	verifyData, ok, err := verifyTransaction(fmt.Sprint(cbData.TrackId))
	if !ok {
		if err != nil {
			log.Printf("Payment -> verify error: %v\n", err)
		}
		paym.IsDone = true
		paym.StatusCode = verifyData.Status
		database.UpdatePayment(paym)
		io.WriteString(w, OrderFailed)
		return
	}

	// transaction successful
	paym.IsDone = true
	paym.StatusCode = verifyData.Status
	paym.ShaparakRefId = fmt.Sprint(verifyData.RefNumber)
	paym.CardNumber = verifyData.CardNumber
	paym.OrderDate = verifyData.PaidAt
	database.UpdatePayment(paym)

	// bot callback
	_, err = tools.AddDaysSub(paym.UserID, 30)
	if err != nil {
		io.WriteString(w, AddSubFailed)
		return
	}
	io.WriteString(w, OrderSuccessful)
	Bot.SendMessage(
		paym.UserID,
		AddSubSuccessful,
		nil,
	)
}

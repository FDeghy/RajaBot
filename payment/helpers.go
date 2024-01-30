package payment

import (
	"RajaBot/config"
	"RajaBot/database"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartPaymentServer() error {
	errChan := make(chan error)
	mux := http.NewServeMux()
	mux.HandleFunc("/verify", handleCallback)
	go func() {
		err := http.ListenAndServe(config.Cfg.Payment.AddressBind, mux)
		errChan <- err
	}()
	select {
	case err := <-errChan:
		return err
	case <-time.After(3 * time.Second):
		break
	}
	log.Println("Payment -> web callback server started.")
	return nil
}

func postToZibal(url string, jsonData []byte) ([]byte, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func NewTransaction(user *database.TgUser, price uint) (*database.Payment, error) {
	uncomp := database.GetUncompletedPayment(user.UserID)
	if len(uncomp) > 0 {
		return nil, ErrUncompletedTransactionFound
	}

	orderId := fmt.Sprintf("%v_%v", user.UserID, time.Now().Unix())

	data := &RequestNewParams{
		Merchant:    config.Cfg.Payment.ApiKey,
		Amount:      price * 10, // toman to rial
		CallbackUrl: config.Cfg.Payment.CallbackDomain + "/verify",
		OrderId:     orderId,
	}

	jsonData, _ := json.Marshal(data)
	body, err := postToZibal(CreateTokenURL, jsonData)
	if err != nil {
		return nil, err
	}
	resp := &ResponseNewParams{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%#v\n", jsonData)
	if resp.Result != 100 {
		return nil, ErrBadCode
	}
	paym := &database.Payment{
		UserID:  user.UserID,
		OrderID: orderId,
		Price:   price,
		TransID: fmt.Sprint(resp.TrackId),
		IsDone:  false,
	}
	database.SavePayment(paym)

	return paym, err
}

func CreateBankLink(transId string) string {
	return fmt.Sprintf(BankURL, transId)
}

func verifyTransaction(transId string) (*ResponseVerifyParams, bool, error) {
	resp := &ResponseVerifyParams{}
	trackId, _ := strconv.ParseInt(transId, 10, 64)
	data := &RequestVerifyParams{
		Merchant: config.Cfg.Payment.ApiKey,
		TrackId:  trackId,
	}
	jsonData, _ := json.Marshal(data)

	body, err := postToZibal(VerifyURL, jsonData)
	if err != nil {
		return resp, false, err
	}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return resp, false, err
	}

	if resp.Result == 100 && (resp.Status == 1 || resp.Status == 2) {
		return resp, true, nil
	}
	return resp, false, nil
}

func CancelTransaction(trainsId string) {
	paym := database.GetPaymentByTransId(trainsId)
	if paym != nil {
		database.DeletePayment(paym)
	}
}

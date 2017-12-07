package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"os/exec"
	"strings"

	"github.com/logpacker/PayPal-Go-SDK"
)

var C, _ = paypalsdk.NewClient("AV2aAlW78rnvuU8EV92wLVsQTnusENXSJJCLSCxo6kUd0nU84ZWdjOoAkt1JNuPP7bk5t3jQ-Wky3on-",
	"EMqdZcKB21emXae8R97HuKzOOISIrnppX06ILsZetgTllfO9hakMr7MvLE548LkqPsKq-Khv7c1UFRMz", paypalsdk.APIBaseSandBox)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello naowal!") // send data to client side
}

func success(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	t, _ := template.ParseFiles("success.gtpl")
	t.Execute(w, nil)

	//print in client side
	fmt.Fprintln(w, "PayerID= "+r.FormValue("PayerID"))
	fmt.Fprintln(w, "paymentId= "+r.FormValue("paymentId"))

	PayerID := r.FormValue("PayerID")
	paymentId := r.FormValue("paymentId")

	//getpayment from payment Id
	payment, _ := C.GetPayment(paymentId)

	//print in console
	fmt.Println("payment createtime : ", payment.CreateTime)
	fmt.Println("payment transactions : ", payment.Transactions)
	fmt.Println("payment  ExperienceProfileID: ", payment.ExperienceProfileID)
	fmt.Println("payment ID: ", payment.ID)
	fmt.Println("payment Intent: ", payment.Intent)
	fmt.Println("payment State: ", payment.State)
	fmt.Println("payment UpdateTime: ", payment.UpdateTime)
	fmt.Println("payment payer FundingInstruments : ", payment.Payer.FundingInstruments)
	fmt.Println("payment payer payerInfo : ", payment.Payer.PayerInfo)
	fmt.Println("payment payer Status : ", payment.Payer.Status)
	fmt.Println("payment payer PaymentMethod : ", payment.Payer.PaymentMethod)

	//Approved payment process
	executeResult, _ := C.ExecuteApprovedPayment(paymentId, PayerID)

	//print Approved Payment from console
	fmt.Println("executeResult : ", executeResult)
	fmt.Println("executeResult ID : ", executeResult.ID)
	fmt.Println("executeResult Links : ", executeResult.Links)
	fmt.Println("executeResult State : ", executeResult.State)
	fmt.Println("executeResult Transactions: ", executeResult.Transactions)

	//print to client side
	fmt.Fprintf(w, "Payment success!") // send data to client side
}

func deny(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Payment Deny!") // send data to client side
}

func payment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "POST" {
		r.ParseForm()
		// logic part of log in
		fmt.Println("Amount:", r.Form["amount"])

		OpenPayment(w, r, strings.Join(r.Form["amount"], "")) //Call OpenPayment function

	} else {
		t, _ := template.ParseFiles("payment.gtpl")
		t.Execute(w, nil)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "POST" {
		r.ParseForm()
		// logic part of log in
		fmt.Println("r Amount:", r.Form["amount"])

		OpenPayment(w, r, strings.Join(r.Form["amount"], "")) //Call OpenPayment function
	} else {
		t, _ := template.ParseFiles("redirect.gtpl")
		t.Execute(w, nil)
	}
}

func OpenPayment(w http.ResponseWriter, r *http.Request, balance string) {

	// In sandbox , Add my Own clientID and secretID
	//C, err := paypalsdk.NewClient("AV2aAlW78rnvuU8EV92wLVsQTnusENXSJJCLSCxo6kUd0nU84ZWdjOoAkt1JNuPP7bk5t3jQ-Wky3on-",
	//	"EMqdZcKB21emXae8R97HuKzOOISIrnppX06ILsZetgTllfO9hakMr7MvLE548LkqPsKq-Khv7c1UFRMz", paypalsdk.APIBaseSandBox)
	//if err != nil {
	//	panic(err)
	//}

	// Retrieve access token
	_, err := C.GetAccessToken()
	if err != nil {
		panic(err)
	}

	// Try to set DirectPaypalPayment
	amount := paypalsdk.Amount{
		Total:    balance, //parse form amount field
		Currency: "USD",
	}
	redirectURI := "http://localhost:8080/success"
	cancelURI := "http://localhost:8080/deny"
	description := "Leaptips following payment"
	paymentResult, err := C.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)
	if err != nil {
		panic(err)
	}
	// find approval_url
	for _, l := range paymentResult.Links {
		if l.Rel == "approval_url" {
			http.Redirect(w, r, l.Href, http.StatusFound)
			return
		}
	}
	return

	payment, err := C.GetPayment(paymentResult.ID)
	//fmt.Println("payment : " + payment)
	fmt.Println("paymentResult.ID: " + paymentResult.ID)
	//fmt.Println("payment.Payer: " + payment.Payer)
	paymentID := paymentResult.ID
	payerID := payment.Payer.PayerInfo.PayerID
	fmt.Println("paymentID : " + paymentID)
	fmt.Println("payerID : " + payerID)
	executeResult, err := C.ExecuteApprovedPayment(paymentID, payerID)
	fmt.Println(executeResult)

}

func redirecttopaypal(w http.ResponseWriter, r *http.Request, url string) {

	http.Redirect(w, r, url, 301)
}

func main() {
	http.HandleFunc("/", sayhelloName) // set router
	http.HandleFunc("/payment", payment)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/success", success)
	http.HandleFunc("/deny", deny)
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/logpacker/PayPal-Go-SDK"
)

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

	//payment, err := C.GetPayment(paymentResult.ID)
	//fmt.Println(payment)
	//fmt.Println(paymentResult.ID)
	//fmt.Println(payment.Payer)
	//paymentID := paymentResult.ID
	//payerID := payment.Payer.PayerInfo.PayerID
	//executeResult, err := C.ExecuteApprovedPayment(paymentID, payerID)
	//fmt.Println(executeResult)
	t, _ := template.ParseFiles("success.gtpl")
	t.Execute(w, nil)

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
	C, err := paypalsdk.NewClient("AV2aAlW78rnvuU8EV92wLVsQTnusENXSJJCLSCxo6kUd0nU84ZWdjOoAkt1JNuPP7bk5t3jQ-Wky3on-",
		"EMqdZcKB21emXae8R97HuKzOOISIrnppX06ILsZetgTllfO9hakMr7MvLE548LkqPsKq-Khv7c1UFRMz", paypalsdk.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	accessToken, err := C.GetAccessToken()

	if err != nil {
		panic(err)
	}

	// Try to set DirectPaypalPayment
	amount := paypalsdk.Amount{
		Total:    balance, //parse form amount field
		Currency: "USD",
	}
	redirectURI := "www.google.com"
	cancelURI := "www.facebook.com"
	description := "Leaptips following payment"
	paymentResult, err := C.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)

	// Just debug in console by printing
	fmt.Println()
	fmt.Println("Token: ", accessToken.Token)
	fmt.Println(paymentResult.Links[0].Rel)
	fmt.Println(paymentResult.Links[0].Href)
	fmt.Println(paymentResult.Links[0].Method)
	fmt.Println(paymentResult.Links[0].Enctype)
	fmt.Println()
	fmt.Println(paymentResult.Links[1].Rel)
	fmt.Println(paymentResult.Links[1].Href)
	fmt.Println(paymentResult.Links[1].Method)
	fmt.Println(paymentResult.Links[1].Enctype)
	fmt.Println()
	fmt.Println(paymentResult.Links[2].Rel)
	fmt.Println(paymentResult.Links[2].Href)
	fmt.Println(paymentResult.Links[2].Method)
	fmt.Println(paymentResult.Links[2].Enctype)

	// open approvel url -> paypal for payment
	exec.Command("xdg-open", paymentResult.Links[1].Href).Run()
	//resp, err := http.Post(paymentResult.Links[2].Href)

	url := paymentResult.Links[0].Href
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	fmt.Fprint(w, responseString)

	//redirecttopaypal(w, r, paymentResult.Links[1].Href)
	// So, Next ?? How am I will do. with accessToken and paymentResult

	payment, err := C.GetPayment(paymentResult.ID)
	fmt.Println(payment)
	fmt.Println(paymentResult.ID)
	fmt.Println(payment.Payer)
	paymentID := paymentResult.ID
	payerID := payment.Payer.PayerInfo.PayerID
	fmt.Println(paymentID)
	fmt.Println(payerID)
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

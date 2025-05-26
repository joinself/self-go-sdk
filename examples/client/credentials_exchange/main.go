package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	// Create a new Self client
	selfClient, err := client.NewClient(client.Config{
		StorageKey:  make([]byte, 32), // In production, use a secure key
		StoragePath: "./credentials_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer selfClient.Close()

	fmt.Printf("My DID: %s\n", selfClient.DID())

	// Set up credential request handlers
	selfClient.Credentials().OnPresentationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("\n📋 Received credential presentation request from: %s\n", req.From())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Type: %v\n", req.Type())
		fmt.Printf("   Expires: %s\n", req.Expires().Format("15:04:05"))

		// Show details of what's being requested
		for i, detail := range req.Details() {
			fmt.Printf("   Detail %d - Credential Type: %v\n", i+1, detail.CredentialType)
			for j, param := range detail.Parameters {
				fmt.Printf("     Parameter %d: %s %s %s\n", j+1, param.Field, operatorToString(param.Operator), param.Value)
			}
		}

		// For demo purposes, we'll reject the request
		// In a real app, you'd look up credentials and respond appropriately
		fmt.Println("   ❌ Rejecting request (demo)")
		err := req.Reject()
		if err != nil {
			fmt.Printf("   Failed to reject request: %v\n", err)
		}
	})

	selfClient.Credentials().OnVerificationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("\n🔍 Received credential verification request from: %s\n", req.From())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Type: %v\n", req.Type())
		fmt.Printf("   Expires: %s\n", req.Expires().Format("15:04:05"))

		// For demo purposes, we'll reject the request
		fmt.Println("   ❌ Rejecting request (demo)")
		err := req.Reject()
		if err != nil {
			fmt.Printf("   Failed to reject request: %v\n", err)
		}
	})

	// Set up credential response handlers
	selfClient.Credentials().OnPresentationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("\n📨 Received credential presentation response from: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", responseStatusToString(resp.Status()))
		fmt.Printf("   Presentations: %d\n", len(resp.Presentations()))

		// Process presentations
		for i, presentation := range resp.Presentations() {
			fmt.Printf("   Presentation %d:\n", i+1)
			fmt.Printf("     Type: %v\n", presentation.PresentationType())
			fmt.Printf("     Holder: %s\n", presentation.Holder().String())
			fmt.Printf("     Credentials: %d\n", len(presentation.Credentials()))
		}
	})

	selfClient.Credentials().OnVerificationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("\n🔐 Received credential verification response from: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", responseStatusToString(resp.Status()))
		fmt.Printf("   Credentials: %d\n", len(resp.Credentials()))

		// Process credentials
		for i, cred := range resp.Credentials() {
			fmt.Printf("   Credential %d:\n", i+1)
			fmt.Printf("     Type: %v\n", cred.CredentialType())
			fmt.Printf("     Subject: %s\n", cred.CredentialSubject().String())
			fmt.Printf("     Issuer: %s\n", cred.Issuer().String())
			fmt.Printf("     Valid From: %s\n", cred.ValidFrom().Format("2006-01-02 15:04:05"))
		}
	})

	// Generate QR code for discovery
	fmt.Println("\n📱 Generating QR code for discovery...")
	qr, err := selfClient.Discovery().GenerateQR()
	if err != nil {
		log.Fatal("Failed to generate QR code:", err)
	}

	fmt.Println("Scan this QR code with another Self client to connect:")
	qrCode, err := qr.Unicode()
	if err != nil {
		log.Fatal("Failed to get QR code:", err)
	}
	fmt.Println(qrCode)

	// Wait for someone to scan the QR code
	fmt.Println("\n⏳ Waiting for connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("No one connected within 2 minutes. Exiting.")
			return
		}
		log.Fatal("Error waiting for connection:", err)
	}

	fmt.Printf("✅ Connected to: %s\n", peer.DID())

	// Wait a moment for the connection to be fully established
	time.Sleep(2 * time.Second)

	// Example 1: Request credential presentations
	fmt.Println("\n📋 Requesting credential presentations...")

	details := []*client.CredentialDetail{
		{
			CredentialType: credential.CredentialTypeEmail,
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "emailAddress",
					Value:    "",
				},
			},
		},
		{
			CredentialType: credential.CredentialTypeLiveness,
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorNotEquals,
					Field:    "sourceImageHash",
					Value:    "",
				},
			},
		},
	}

	presentationReq, err := selfClient.Credentials().RequestPresentationWithTimeout(peer.DID(), details, 30*time.Second)
	if err != nil {
		log.Printf("Failed to request presentations: %v", err)
	} else {
		fmt.Printf("   Request ID: %s\n", presentationReq.RequestID())

		// Wait for response
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		defer cancel()

		resp, err := presentationReq.WaitForResponse(ctx)
		if err != nil {
			fmt.Printf("   ❌ Request failed or timed out: %v\n", err)
		} else {
			fmt.Printf("   ✅ Received response with status: %s\n", responseStatusToString(resp.Status()))
		}
	}

	// Example 2: Request credential verification
	fmt.Println("\n🔍 Requesting credential verification...")

	verificationReq, err := selfClient.Credentials().RequestVerificationWithTimeout(
		peer.DID(),
		credential.CredentialTypeEmail,
		30*time.Second,
	)
	if err != nil {
		log.Printf("Failed to request verification: %v", err)
	} else {
		fmt.Printf("   Request ID: %s\n", verificationReq.RequestID())

		// Wait for response
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		defer cancel()

		resp, err := verificationReq.WaitForResponse(ctx)
		if err != nil {
			fmt.Printf("   ❌ Request failed or timed out: %v\n", err)
		} else {
			fmt.Printf("   ✅ Received response with status: %s\n", responseStatusToString(resp.Status()))
		}
	}

	fmt.Println("\n🎉 Credential exchange demo completed!")
	fmt.Println("The client will continue listening for incoming requests...")
	fmt.Println("Press Ctrl+C to exit.")

	// Keep the program running to receive credential requests
	select {}
}

func operatorToString(op message.ComparisonOperator) string {
	switch op {
	case message.OperatorEquals:
		return "=="
	case message.OperatorNotEquals:
		return "!="
	case message.OperatorGreaterThan:
		return ">"
	case message.OperatorLessThan:
		return "<"
	case message.OperatorGreaterThanOrEquals:
		return ">="
	case message.OperatorLessThanOrEquals:
		return "<="
	default:
		return "unknown"
	}
}

func responseStatusToString(status message.ResponseStatus) string {
	switch status {
	case message.ResponseStatusAccepted:
		return "Accepted"
	case message.ResponseStatusOk:
		return "OK"
	case message.ResponseStatusCreated:
		return "Created"
	case message.ResponseStatusBadRequest:
		return "Bad Request"
	case message.ResponseStatusUnauthorized:
		return "Unauthorized"
	case message.ResponseStatusForbidden:
		return "Forbidden"
	case message.ResponseStatusNotFound:
		return "Not Found"
	case message.ResponseStatusNotAcceptable:
		return "Not Acceptable"
	case message.ResponseStatusConflict:
		return "Conflict"
	default:
		return "Unknown"
	}
}

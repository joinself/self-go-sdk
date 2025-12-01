// Package main demonstrates advanced credential exchange using the Self SDK.
//
// This is the ADVANCED level of credential exchange examples.
// Complete basic_exchange.go and multi_credential_exchange.go first.
//
// This example shows:
// - Complex parameter configurations and filtering
// - Credential verification requests vs presentation requests
// - Advanced response processing and validation
// - Error handling and edge cases
//
// 🎯 What you'll learn:
// • Complex credential filtering with operators
// • Difference between presentation and verification requests
// • Advanced response processing patterns
// • Production-ready error handling
//
// 📚 Prerequisites: basic_exchange.go, multi_credential_exchange.go
// 📚 Next steps: discovery_exchange.go - QR code discovery integration
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/examples/utils"
	"github.com/joinself/self-go-sdk/message"
)

func main() {
	fmt.Println("🔄 Advanced Credential Exchange Demo")
	fmt.Println("====================================")
	fmt.Println("This demo shows advanced credential exchange patterns.")
	fmt.Println("📚 This is the ADVANCED level - complete previous examples first.")
	fmt.Println()

	// Step 1: Create clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create credentials with complex data for filtering
	createAdvancedCredentials(issuer, holder)

	// Step 3: Set up advanced handlers
	setupAdvancedHandlers(issuer, holder)

	// Step 4: Demonstrate advanced exchange patterns
	demonstrateComplexFiltering(issuer, holder)
	demonstrateVerificationRequest(issuer, holder)
	demonstrateAdvancedResponseProcessing(issuer, holder)

	fmt.Println("✅ Advanced demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the final level?")
	fmt.Println("   • Run discovery_exchange.go for QR code discovery integration")
	fmt.Println()
	fmt.Println("The clients will keep running. Press Ctrl+C to exit.")

	select {}
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	issuer, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("advanced_issuer"),
		StoragePath: "./advanced_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	holder, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("advanced_holder"),
		StoragePath: "./advanced_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// createAdvancedCredentials creates credentials with complex data for filtering demonstrations
func createAdvancedCredentials(issuer, holder *client.Client) {
	fmt.Println("📝 Creating credentials with complex data...")

	// Organization credential with numeric and date fields for complex filtering
	fmt.Println("🏢 Creating organization credential with complex claims...")
	_, err := issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeOrganisation).
		Subject(holder.DID()).
		Issuer(issuer.DID()).
		Claims(map[string]interface{}{
			"organizationName": "TechCorp Inc.",
			"employeeId":       "EMP-2024-001",
			"position": map[string]interface{}{
				"title":     "Senior Software Engineer",
				"level":     5,
				"startDate": "2024-01-15",
				"salary":    95000,
			},
			"permissions": []string{
				"read:repositories",
				"write:code",
				"deploy:staging",
			},
			"performance": map[string]interface{}{
				"rating":        4.8,
				"lastReview":    "2024-06-15",
				"nextReview":    "2024-12-15",
				"bonusEligible": true,
			},
		}).
		ValidFrom(time.Now()).
		SignWith(issuer.DID(), time.Now()).
		Issue(issuer)

	if err != nil {
		log.Printf("Failed to create organization credential: %v", err)
	} else {
		fmt.Println("   ✅ Organization credential created with complex claims")
	}

	fmt.Println("✅ Advanced credentials created successfully")
	fmt.Println()
}

// setupAdvancedHandlers configures handlers for advanced exchange patterns
func setupAdvancedHandlers(issuer, holder *client.Client) {
	fmt.Println("🔧 Setting up advanced handlers...")

	// Advanced presentation request handler
	holder.Credentials().OnPresentationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("📨 Holder received PRESENTATION request from: %s\n", req.From())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Expires: %s\n", req.Expires().Format("15:04:05"))

		// Process complex parameters
		fmt.Println("   🔍 Analyzing complex parameters:")
		for i, detail := range req.Details() {
			fmt.Printf("     Detail %d - Type: %v\n", i+1, detail.CredentialType)
			for j, param := range detail.Parameters {
				fmt.Printf("       Parameter %d: %s %s %v\n", j+1, param.Field, operatorToString(param.Operator), param.Value)
			}
		}

		fmt.Println("   ❌ Rejecting presentation request (demo)")
		req.Reject()
	})

	// Advanced verification request handler
	holder.Credentials().OnVerificationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("📨 Holder received VERIFICATION request from: %s\n", req.From())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Requested types: %v\n", req.Type())
		fmt.Printf("   Expires: %s\n", req.Expires().Format("15:04:05"))

		fmt.Println("   🔍 Verification request processing...")
		fmt.Println("      Verification requests validate credential authenticity")
		fmt.Println("      Different from presentation requests which share credential data")

		fmt.Println("   ❌ Rejecting verification request (demo)")
		req.Reject()
	})

	// Advanced response handlers
	issuer.Credentials().OnPresentationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("📨 Issuer received PRESENTATION response from: %s\n", resp.From())
		processAdvancedResponse("PRESENTATION", resp)
	})

	issuer.Credentials().OnVerificationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("📨 Issuer received VERIFICATION response from: %s\n", resp.From())
		processAdvancedResponse("VERIFICATION", resp)
	})

	fmt.Println("✅ Advanced handlers configured")
	fmt.Println()
}

// demonstrateComplexFiltering shows advanced parameter configurations
func demonstrateComplexFiltering(issuer, holder *client.Client) {
	fmt.Println("1️⃣ COMPLEX PARAMETER FILTERING")
	fmt.Println("===============================")
	fmt.Println("🔍 Demonstrating complex credential filtering...")

	// Create request with complex filtering parameters
	details := []*client.CredentialDetail{
		{
			CredentialType: credential.CredentialTypeOrganisation,
			Parameters: []*client.CredentialParameter{
				{
					Operator: message.OperatorEquals,
					Field:    "organizationName",
					Value:    "TechCorp Inc.",
				},
				{
					Operator: message.OperatorGreaterThanOrEquals,
					Field:    "position.level",
					Value:    "5",
				},
				{
					Operator: message.OperatorGreaterThan,
					Field:    "position.salary",
					Value:    "90000",
				},
				{
					Operator: message.OperatorEquals,
					Field:    "performance.bonusEligible",
					Value:    "true",
				},
			},
		},
	}

	fmt.Println("📤 Requesting organization credential with complex filters:")
	fmt.Println("   🏢 Organization = 'TechCorp Inc.'")
	fmt.Println("   📊 Position level >= 5")
	fmt.Println("   💰 Salary > $90,000")
	fmt.Println("   🎯 Bonus eligible = true")

	req, err := issuer.Credentials().RequestPresentationWithTimeout(
		holder.DID(),
		details,
		15*time.Second,
	)
	if err != nil {
		log.Printf("Failed to send complex filtering request: %v", err)
		return
	}

	fmt.Printf("   Request sent with ID: %s\n", req.RequestID())

	// Wait for response
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := req.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("   ⏰ Complex filtering request timed out (expected)")
		} else {
			fmt.Printf("   ❌ Request failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Response received: %s\n", utils.ResponseStatusToString(resp.Status()))
	}

	fmt.Println()
}

// demonstrateVerificationRequest shows the difference between presentation and verification
func demonstrateVerificationRequest(issuer, holder *client.Client) {
	fmt.Println("2️⃣ CREDENTIAL VERIFICATION REQUEST")
	fmt.Println("===================================")
	fmt.Println("🔐 Demonstrating verification vs presentation...")

	fmt.Println("📤 Sending verification request (validates authenticity)...")
	fmt.Println("   🔍 Verification checks if credentials are valid and authentic")
	fmt.Println("   📋 Different from presentation which shares credential data")

	// Send verification request
	verificationReq, err := issuer.Credentials().RequestVerificationWithTimeout(
		holder.DID(),
		credential.CredentialTypeOrganisation,
		15*time.Second,
	)
	if err != nil {
		log.Printf("Failed to send verification request: %v", err)
		return
	}

	fmt.Printf("   Verification request sent with ID: %s\n", verificationReq.RequestID())

	// Wait for verification response
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := verificationReq.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("   ⏰ Verification request timed out (expected)")
		} else {
			fmt.Printf("   ❌ Verification failed: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Verification response: %s\n", utils.ResponseStatusToString(resp.Status()))
	}

	fmt.Println()
}

// demonstrateAdvancedResponseProcessing shows sophisticated response handling
func demonstrateAdvancedResponseProcessing(issuer, holder *client.Client) {
	fmt.Println("3️⃣ ADVANCED RESPONSE PROCESSING")
	fmt.Println("================================")
	fmt.Println("⚙️ Demonstrating sophisticated response handling...")

	// This would normally process actual responses, but for demo we'll show the pattern
	fmt.Println("📊 Advanced response processing includes:")
	fmt.Println("   • Credential validation and verification")
	fmt.Println("   • Complex claim extraction and analysis")
	fmt.Println("   • Business rule validation")
	fmt.Println("   • Error handling and retry logic")
	fmt.Println("   • Audit logging and compliance tracking")
	fmt.Println()

	fmt.Println("🎓 Key differences in advanced processing:")
	fmt.Println("   📋 Presentation responses: contain actual credential data")
	fmt.Println("   🔐 Verification responses: contain validation results")
	fmt.Println("   ⚡ Complex filtering: enables precise credential selection")
	fmt.Println("   🛡️ Error handling: ensures robust production workflows")
	fmt.Println()
}

// processAdvancedResponse demonstrates sophisticated response processing
func processAdvancedResponse(requestType string, resp *client.CredentialResponse) {
	fmt.Printf("   Status: %s\n", utils.ResponseStatusToString(resp.Status()))
	fmt.Printf("   Type: %s response\n", requestType)

	if requestType == "PRESENTATION" {
		fmt.Printf("   Presentations: %d\n", len(resp.Presentations()))
		for i, presentation := range resp.Presentations() {
			fmt.Printf("     Presentation %d: %v\n", i+1, presentation.PresentationType())
			fmt.Printf("       Credentials: %d\n", len(presentation.Credentials()))
		}
	} else if requestType == "VERIFICATION" {
		fmt.Printf("   Credentials verified: %d\n", len(resp.Credentials()))
		for i, cred := range resp.Credentials() {
			fmt.Printf("     Credential %d: %v\n", i+1, cred.CredentialType())
			fmt.Printf("       Valid from: %s\n", cred.ValidFrom().Format("2006-01-02"))
		}
	}
}

// operatorToString converts comparison operators to readable strings
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

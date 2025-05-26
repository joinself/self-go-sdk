// Package main demonstrates comprehensive credential issuance capabilities using the Self SDK client facade.
//
// This example shows how to:
// - Initialize Self clients for issuer and holder roles
// - Create various types of verifiable credentials using the builder pattern
// - Attach evidence/files to credentials for enhanced verification
// - Handle complex nested claims and data structures
// - Create verifiable presentations from credentials
// - Set up credential request/response handlers
// - Manage asset uploads and downloads for evidence
//
// The Self SDK provides decentralized identity and verifiable credential capabilities,
// allowing entities to issue, hold, and verify credentials without requiring
// centralized authorities while maintaining cryptographic integrity and privacy.
//
// 🎯 CREDENTIAL CAPABILITIES DEMONSTRATED:
// • Basic credential creation (Email verification)
// • Multi-claim credentials (Profile information)
// • Custom credentials with file evidence (Certifications)
// • Complex nested data structures (Organization credentials)
// • Credential builder pattern usage
// • Asset/evidence management (file uploads)
// • Verifiable presentation creation
// • Request/response handling workflows
//
// 🔧 KEY SDK COMPONENTS SHOWCASED:
// • client.NewClient() - Client initialization and configuration
// • NewCredentialBuilder() - Fluent API for credential construction
// • CreateAsset() - Evidence and file attachment management
// • CreatePresentation() - Verifiable presentation creation
// • OnVerificationRequest/Response() - Event-driven credential workflows
//
// 📚 EDUCATIONAL PROGRESSION:
// The examples progress from simple to complex, building understanding:
// 1. Basic Email Credential - Simplest form with minimal claims
// 2. Profile Credential - Multiple claims in a single credential
// 3. Custom Credential with Evidence - File attachments and presentations
// 4. Organization Credential - Complex nested data and arrays
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/examples/utils"
)

const (
	// Configuration constants for demo setup
	issuerStorageDir = "./issuer_storage"
	holderStorageDir = "./holder_storage"
)

func main() {
	fmt.Println("🎓 Self SDK Credential Issuance Demo")
	fmt.Println("=====================================")
	fmt.Println("📚 This demo showcases comprehensive credential issuance capabilities:")
	fmt.Println("   • Creating various types of verifiable credentials")
	fmt.Println("   • Using the credential builder pattern")
	fmt.Println("   • Attaching evidence and files to credentials")
	fmt.Println("   • Managing complex nested claims")
	fmt.Println("   • Creating verifiable presentations")
	fmt.Println("   • Handling credential request/response workflows")
	fmt.Println()

	// 🏗️ STEP 1: CLIENT SETUP - Initialize issuer and holder clients
	// The issuer creates and signs credentials, while the holder receives and stores them
	issuerClient, holderClient := setupClients()
	defer issuerClient.Close()
	defer holderClient.Close()

	// 🆔 IDENTITY DISPLAY: Show the unique DIDs for both parties
	// DIDs (Decentralized Identifiers) are cryptographically verifiable identities
	fmt.Printf("🏢 Issuer DID: %s\n", issuerClient.DID())
	fmt.Printf("   This is the credential issuer's unique decentralized identity\n")
	fmt.Printf("👤 Holder DID: %s\n", holderClient.DID())
	fmt.Printf("   This is the credential holder's unique decentralized identity\n")
	fmt.Println()

	// 🔧 STEP 2: HANDLER SETUP - Configure credential request/response handlers
	// These handlers demonstrate how to process incoming credential requests
	setupCredentialHandlers(issuerClient, holderClient)

	// 📚 STEP 3: CREDENTIAL ISSUANCE EXAMPLES
	// Progressive examples from simple to complex credential types
	fmt.Println("📚 CREDENTIAL ISSUANCE EXAMPLES")
	fmt.Println("================================")
	fmt.Println("🎯 The following examples demonstrate progressive complexity:")
	fmt.Println("   Each example builds upon concepts from the previous ones")
	fmt.Println()

	// 📧 EXAMPLE 1: Basic Email Credential - Foundation concepts
	demonstrateBasicCredential(issuerClient, holderClient)

	// 👤 EXAMPLE 2: Profile Credential - Multiple claims
	demonstrateProfileCredential(issuerClient, holderClient)

	// 🎓 EXAMPLE 3: Custom Credential - Evidence and presentations
	demonstrateCustomCredentialWithEvidence(issuerClient, holderClient)

	// 🏢 EXAMPLE 4: Organization Credential - Complex data structures
	demonstrateOrganizationCredential(issuerClient, holderClient)

	// 🔗 STEP 4: OPTIONAL DISCOVERY DEMO
	// Discovery workflow is separated to maintain focus on credential issuance
	fmt.Println("\n🔗 DISCOVERY & CONNECTION (Optional)")
	fmt.Println("====================================")
	fmt.Println("📱 The discovery workflow demonstrates peer-to-peer connections")
	fmt.Println("   For credential issuance focus, this section is optional")
	fmt.Println("   Uncomment runDiscoveryDemo() below to enable QR code discovery")
	fmt.Println()

	// Uncomment the line below to run the discovery demo
	// runDiscoveryDemo(issuerClient, holderClient)

	// 🎉 STEP 5: EDUCATIONAL SUMMARY
	// Comprehensive summary of demonstrated features and next steps
	printSummary()
}

// setupClients demonstrates client initialization and configuration
// This function showcases how to:
// - Create Self SDK clients with proper configuration
// - Set up storage paths for different client roles
// - Configure environment and logging settings
// - Handle client lifecycle management
func setupClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 SETTING UP SELF SDK CLIENTS")
	fmt.Println("===============================")
	fmt.Println("🏗️ Initializing issuer and holder clients...")
	fmt.Println("   The issuer creates and signs credentials")
	fmt.Println("   The holder receives and manages credentials")
	fmt.Println()

	// 🏢 ISSUER CLIENT: Creates and signs verifiable credentials
	// The issuer client has the authority to create credentials for subjects
	fmt.Println("🏢 Creating issuer client...")
	issuerClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("issuer"), // Unique key for issuer storage encryption
		StoragePath: issuerStorageDir,                   // Dedicated storage directory for issuer
		Environment: client.Sandbox,                     // Use Sandbox environment for development
		LogLevel:    client.LogInfo,                     // Show informational log messages
	})
	if err != nil {
		log.Fatal("❌ Failed to create issuer client:", err)
	}

	// 👤 HOLDER CLIENT: Receives and stores verifiable credentials
	// The holder client manages credentials issued by various issuers
	fmt.Println("👤 Creating holder client...")
	holderClient, err := client.NewClient(client.Config{
		StorageKey:  utils.GenerateStorageKey("holder"), // Unique key for holder storage encryption
		StoragePath: holderStorageDir,                   // Dedicated storage directory for holder
		Environment: client.Sandbox,                     // Use Sandbox environment for development
		LogLevel:    client.LogInfo,                     // Show informational log messages
	})
	if err != nil {
		log.Fatal("❌ Failed to create holder client:", err)
	}

	fmt.Println("✅ Clients created successfully")
	fmt.Println("   🔐 Both clients use encrypted local storage")
	fmt.Println("   🌐 Connected to Self Sandbox environment")
	fmt.Println()

	return issuerClient, holderClient
}

// setupCredentialHandlers demonstrates credential request/response handling
// This function showcases how to:
// - Register handlers for incoming verification requests
// - Process credential requests with proper responses
// - Handle verification responses from peers
// - Implement event-driven credential workflows
func setupCredentialHandlers(issuerClient, holderClient *client.Client) {
	fmt.Println("🔧 SETTING UP CREDENTIAL HANDLERS")
	fmt.Println("==================================")
	fmt.Println("📨 Configuring request/response handlers...")
	fmt.Println("   These handlers process incoming credential requests")
	fmt.Println("   In production, implement business logic for credential validation")
	fmt.Println()

	// 📨 VERIFICATION REQUEST HANDLER: Process incoming credential verification requests
	// This handler runs when someone requests credential verification from the holder
	holderClient.Credentials().OnVerificationRequest(func(req *client.IncomingCredentialRequest) {
		fmt.Printf("📨 VERIFICATION REQUEST RECEIVED\n")
		fmt.Printf("   From: %s\n", req.From())
		fmt.Printf("   Requested types: %v\n", req.Type())
		fmt.Printf("   Request ID: %s\n", req.RequestID())
		fmt.Printf("   Evidence items: %d\n", len(req.Evidence()))
		fmt.Printf("   Proof presentations: %d\n", len(req.Proof()))

		// 🔄 DEMO RESPONSE: For demonstration, we reject requests
		// In production, implement logic to:
		// - Validate the request against business rules
		// - Check if holder has requested credentials
		// - Respond with appropriate credentials or rejection
		fmt.Println("   ❌ Rejecting request (demo - no credentials to share)")
		fmt.Println("      In production: implement credential lookup and validation")
		err := req.Reject()
		if err != nil {
			fmt.Printf("   ❌ Failed to reject request: %v\n", err)
		} else {
			fmt.Printf("   ✅ Request rejected successfully\n")
		}
		fmt.Println()
	})

	// 📨 VERIFICATION RESPONSE HANDLER: Process credential verification responses
	// This handler runs when the issuer receives responses to verification requests
	issuerClient.Credentials().OnVerificationResponse(func(resp *client.CredentialResponse) {
		fmt.Printf("📨 VERIFICATION RESPONSE RECEIVED\n")
		fmt.Printf("   From: %s\n", resp.From())
		fmt.Printf("   Status: %s\n", utils.ResponseStatusToString(resp.Status()))
		fmt.Printf("   Credentials received: %d\n", len(resp.Credentials()))

		// 🔍 CREDENTIAL PROCESSING: In production, validate and process received credentials
		for i, cred := range resp.Credentials() {
			fmt.Printf("   Credential %d: %v\n", i+1, cred.CredentialType())
		}
		fmt.Println()
	})

	fmt.Println("✅ Handlers configured successfully")
	fmt.Println("   📨 Ready to process credential requests and responses")
	fmt.Println("   🔄 Event-driven workflow established")
	fmt.Println()
}

// demonstrateBasicCredential showcases the simplest form of credential issuance
// This example demonstrates:
// - Basic credential builder usage
// - Simple claim addition
// - Credential signing and issuance
// - Foundation concepts for all credential types
func demonstrateBasicCredential(issuerClient, holderClient *client.Client) {
	fmt.Println("1️⃣ BASIC EMAIL CREDENTIAL")
	fmt.Println("==========================")
	fmt.Println("📧 Creating a simple email verification credential...")
	fmt.Println("   This demonstrates the foundation of credential issuance")
	fmt.Println("   Key concepts: builder pattern, claims, signing, issuance")
	fmt.Println()

	// 🏗️ CREDENTIAL BUILDER: Use the fluent builder pattern for credential creation
	// The builder provides a clean, readable API for constructing credentials
	fmt.Println("🏗️ Using credential builder pattern...")
	emailCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeEmail).                       // Set credential type to email
		Subject(holderClient.DID()).                                // Specify who the credential is about
		Issuer(issuerClient.DID()).                                 // Specify who is issuing the credential
		Claim("emailAddress", "john.doe@example.com").              // Add email address claim
		Claim("verified", true).                                    // Add verification status claim
		Claim("verificationDate", time.Now().Format("2006-01-02")). // Add verification date claim
		ValidFrom(time.Now()).                                      // Set when credential becomes valid
		SignWith(issuerClient.DID(), time.Now()).                   // Sign with issuer's key
		Issue(issuerClient)                                         // Issue the credential

	if err != nil {
		log.Printf("❌ Failed to create email credential: %v", err)
		return
	}

	// ✅ SUCCESS REPORTING: Display credential creation results
	fmt.Printf("   ✅ Email credential created successfully\n")
	fmt.Printf("   📧 Email: john.doe@example.com\n")
	fmt.Printf("   ✔️  Verified: true\n")
	fmt.Printf("   📅 Verification Date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("   🔒 Credential Type: %v\n", emailCredential.CredentialType())
	fmt.Printf("   🆔 Subject: %s\n", emailCredential.CredentialSubject().String())
	fmt.Printf("   🏢 Issuer: %s\n", emailCredential.Issuer().String())
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Credentials contain claims about a subject")
	fmt.Println("   • Builder pattern provides clean, readable construction")
	fmt.Println("   • Cryptographic signatures ensure integrity")
	fmt.Println("   • Timestamps establish validity periods")
	fmt.Println()
}

// demonstrateProfileCredential showcases credentials with multiple claims
// This example demonstrates:
// - Adding multiple claims to a single credential
// - Different data types in claims
// - Organizing related information in one credential
// - Building upon basic credential concepts
func demonstrateProfileCredential(issuerClient, holderClient *client.Client) {
	fmt.Println("2️⃣ PROFILE CREDENTIAL WITH MULTIPLE CLAIMS")
	fmt.Println("===========================================")
	fmt.Println("👤 Creating a profile credential with multiple claims...")
	fmt.Println("   This demonstrates how to include multiple pieces of information")
	fmt.Println("   in a single credential for related identity attributes")
	fmt.Println()

	// 🏗️ MULTI-CLAIM BUILDER: Demonstrate adding multiple related claims
	fmt.Println("🏗️ Building credential with multiple claims...")
	profileCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeProfileName).                 // Profile credential type
		Subject(holderClient.DID()).                                // Subject of the credential
		Issuer(issuerClient.DID()).                                 // Credential issuer
		Claim("firstName", "John").                                 // First name claim
		Claim("lastName", "Doe").                                   // Last name claim
		Claim("displayName", "John Doe").                           // Display name claim
		Claim("profileLevel", "verified").                          // Verification level claim
		Claim("country", "United States").                          // Country claim
		Claim("registrationDate", time.Now().Format("2006-01-02")). // Registration date
		ValidFrom(time.Now()).                                      // Validity start time
		SignWith(issuerClient.DID(), time.Now()).                   // Cryptographic signature
		Issue(issuerClient)                                         // Issue the credential

	if err != nil {
		log.Printf("❌ Failed to create profile credential: %v", err)
		return
	}

	// ✅ SUCCESS REPORTING: Display comprehensive credential information
	fmt.Printf("   ✅ Profile credential created successfully\n")
	fmt.Printf("   👤 Name: John Doe\n")
	fmt.Printf("   🌍 Country: United States\n")
	fmt.Printf("   ⭐ Profile Level: verified\n")
	fmt.Printf("   📅 Registration: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("   🔒 Credential Type: %v\n", profileCredential.CredentialType())
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Multiple related claims can be grouped in one credential")
	fmt.Println("   • Claims can contain different data types (strings, booleans, dates)")
	fmt.Println("   • Grouping related information improves efficiency")
	fmt.Println("   • Each claim is cryptographically protected")
	fmt.Println()
}

// demonstrateCustomCredentialWithEvidence showcases advanced credential features
// This example demonstrates:
// - Creating custom credential types
// - Attaching file evidence to credentials
// - Asset management and upload functionality
// - Creating verifiable presentations
// - Linking evidence to credential claims
func demonstrateCustomCredentialWithEvidence(issuerClient, holderClient *client.Client) {
	fmt.Println("3️⃣ CUSTOM CREDENTIAL WITH EVIDENCE")
	fmt.Println("===================================")
	fmt.Println("🎓 Creating a certification credential with file evidence...")
	fmt.Println("   This demonstrates advanced features: custom types, evidence, presentations")
	fmt.Println("   Evidence provides additional proof supporting credential claims")
	fmt.Println()

	// 📄 EVIDENCE CREATION: Create and upload supporting documentation
	fmt.Println("📄 Creating evidence asset...")
	fmt.Println("   Evidence can be any file type: PDFs, images, documents, etc.")
	certificateData := []byte("This is a mock certificate document for demonstration purposes.\n" +
		"Certificate of Completion\n" +
		"Advanced Go Programming Course\n" +
		"Student: John Doe\n" +
		"Grade: A+\n" +
		"Date: " + time.Now().Format("2006-01-02"))

	evidence, err := issuerClient.Credentials().CreateAsset("certificate.pdf", "application/pdf", certificateData)
	if err != nil {
		log.Printf("❌ Failed to create evidence: %v", err)
		return
	}

	fmt.Printf("   📄 Evidence created: %s\n", evidence.Name)
	fmt.Printf("   🔗 Asset ID: %x\n", evidence.ID())
	fmt.Printf("   🔐 Content Hash: %x\n", evidence.Hash())
	fmt.Println("   ✅ Evidence uploaded to secure storage")
	fmt.Println()

	// 🏗️ CUSTOM CREDENTIAL: Create credential with evidence reference
	fmt.Println("🏗️ Building custom certification credential...")
	customCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type([]string{"VerifiableCredential", "CertificationCredential"}). // Custom credential type
		Subject(holderClient.DID()).                                       // Credential subject
		Issuer(issuerClient.DID()).                                        // Credential issuer
		Claim("certificationType", "Professional Development").            // Type of certification
		Claim("courseName", "Advanced Go Programming").                    // Course name
		Claim("completionDate", time.Now().Format("2006-01-02")).          // Completion date
		Claim("certificateHash", fmt.Sprintf("%x", evidence.Hash())).      // Link to evidence
		Claim("grade", "A+").                                              // Achievement grade
		Claim("institution", "Self SDK Academy").                          // Issuing institution
		Claim("courseHours", 40).                                          // Course duration
		ValidFrom(time.Now()).                                             // Validity period
		SignWith(issuerClient.DID(), time.Now()).                          // Cryptographic signature
		Issue(issuerClient)                                                // Issue credential

	if err != nil {
		log.Printf("❌ Failed to create custom credential: %v", err)
		return
	}

	// 📋 PRESENTATION CREATION: Create verifiable presentation from credential
	fmt.Println("📋 Creating verifiable presentation...")
	fmt.Println("   Presentations package credentials for sharing with verifiers")
	presentation, err := createPresentation(issuerClient, customCredential)
	if err != nil {
		log.Printf("❌ Failed to create presentation: %v", err)
		return
	}

	// ✅ SUCCESS REPORTING: Display comprehensive results
	fmt.Printf("   ✅ Certification credential created successfully\n")
	fmt.Printf("   🎓 Course: Advanced Go Programming\n")
	fmt.Printf("   📅 Completed: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("   🏆 Grade: A+\n")
	fmt.Printf("   🏫 Institution: Self SDK Academy\n")
	fmt.Printf("   ⏱️  Duration: 40 hours\n")
	fmt.Printf("   🔒 Credential Type: %v\n", customCredential.CredentialType())
	fmt.Printf("   📋 Presentation Type: %v\n", presentation.PresentationType())
	fmt.Printf("   🔗 Evidence Hash: %x\n", evidence.Hash())
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Custom credential types support specific use cases")
	fmt.Println("   • Evidence provides additional verification material")
	fmt.Println("   • Asset management handles secure file storage")
	fmt.Println("   • Presentations package credentials for sharing")
	fmt.Println("   • Hash references link credentials to evidence")
	fmt.Println()
}

// demonstrateOrganizationCredential showcases complex data structures in credentials
// This example demonstrates:
// - Complex nested objects in claims
// - Arrays and collections in credentials
// - Hierarchical data organization
// - Real-world organizational data modeling
// - Advanced claim structuring techniques
func demonstrateOrganizationCredential(issuerClient, holderClient *client.Client) {
	fmt.Println("4️⃣ ORGANIZATION CREDENTIAL WITH COMPLEX CLAIMS")
	fmt.Println("===============================================")
	fmt.Println("🏢 Creating an organization credential with complex nested data...")
	fmt.Println("   This demonstrates advanced data structures: nested objects, arrays")
	fmt.Println("   Real-world credentials often contain hierarchical information")
	fmt.Println()

	// 🏗️ COMPLEX CLAIMS: Demonstrate nested objects and arrays
	fmt.Println("🏗️ Building credential with complex nested claims...")
	orgCredential, err := issuerClient.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeOrganisation). // Organization credential type
		Subject(holderClient.DID()).                 // Employee subject
		Issuer(issuerClient.DID()).                  // Organization issuer
		Claims(map[string]interface{}{               // Complex claims structure
			"organizationName": "TechCorp Inc.", // Company name
			"employeeId":       "EMP-2024-001",  // Employee identifier
			"position": map[string]interface{}{ // Nested position object
				"title":      "Senior Software Engineer", // Job title
				"department": "Engineering",              // Department
				"level":      "L5",                       // Career level
				"startDate":  "2024-01-15",               // Start date
				"manager":    "jane.smith@techcorp.com",  // Manager reference
			},
			"permissions": []string{ // Array of permissions
				"read:repositories",    // Repository access
				"write:code",           // Code modification
				"deploy:staging",       // Staging deployment
				"review:pull-requests", // Code review
				"admin:team-resources", // Team administration
			},
			"contact": map[string]interface{}{ // Contact information
				"email":    "john.doe@techcorp.com",        // Work email
				"phone":    "+1-555-0123",                  // Work phone
				"office":   "Building A, Floor 3, Desk 42", // Office location
				"timezone": "America/New_York",             // Timezone
			},
			"benefits": map[string]interface{}{ // Benefits package
				"healthInsurance": true, // Health coverage
				"retirement401k":  true, // Retirement plan
				"paidTimeOff":     25,   // PTO days
				"stockOptions":    1000, // Stock options
				"remoteWork":      true, // Remote work eligibility
			},
			"certifications": []map[string]interface{}{ // Array of certifications
				{
					"name":       "AWS Solutions Architect", // Certification name
					"level":      "Professional",            // Certification level
					"issueDate":  "2023-06-15",              // Issue date
					"expiryDate": "2026-06-15",              // Expiry date
					"verified":   true,                      // Verification status
				},
				{
					"name":       "Kubernetes Administrator", // Second certification
					"level":      "Certified",                // Certification level
					"issueDate":  "2023-09-20",               // Issue date
					"expiryDate": "2026-09-20",               // Expiry date
					"verified":   true,                       // Verification status
				},
			},
		}).
		ValidFrom(time.Now()).                    // Validity start
		SignWith(issuerClient.DID(), time.Now()). // Cryptographic signature
		Issue(issuerClient)                       // Issue credential

	if err != nil {
		log.Printf("❌ Failed to create organization credential: %v", err)
		return
	}

	// ✅ SUCCESS REPORTING: Display comprehensive organizational information
	fmt.Printf("   ✅ Organization credential created successfully\n")
	fmt.Printf("   🏢 Company: TechCorp Inc.\n")
	fmt.Printf("   💼 Position: Senior Software Engineer (L5)\n")
	fmt.Printf("   🏬 Department: Engineering\n")
	fmt.Printf("   🆔 Employee ID: EMP-2024-001\n")
	fmt.Printf("   📧 Email: john.doe@techcorp.com\n")
	fmt.Printf("   📍 Office: Building A, Floor 3, Desk 42\n")
	fmt.Printf("   🔑 Permissions: 5 access levels\n")
	fmt.Printf("   🎯 Benefits: Health, 401k, 25 PTO days, Stock options\n")
	fmt.Printf("   🏆 Certifications: 2 professional certifications\n")
	fmt.Printf("   🔒 Credential Type: %v\n", orgCredential.CredentialType())
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Credentials can contain complex nested data structures")
	fmt.Println("   • Arrays enable multiple values for single claim types")
	fmt.Println("   • Hierarchical organization mirrors real-world data")
	fmt.Println("   • Complex claims maintain cryptographic integrity")
	fmt.Println("   • Structured data enables precise verification queries")
	fmt.Println()
}

// runDiscoveryDemo demonstrates the QR code-based peer discovery workflow
// This function showcases how to:
// - Generate QR codes for peer discovery
// - Handle peer connections and responses
// - Integrate discovery with credential workflows
// - Manage connection timeouts and error handling
func runDiscoveryDemo(issuerClient, holderClient *client.Client) {
	fmt.Println("🔗 PEER DISCOVERY DEMONSTRATION")
	fmt.Println("===============================")
	fmt.Println("📱 Generating QR code for peer discovery...")
	fmt.Println("   Discovery enables secure peer-to-peer connections")
	fmt.Println("   QR codes contain cryptographic material for secure handshake")
	fmt.Println()

	// 🔑 QR GENERATION: Create discovery QR code with embedded crypto material
	qr, err := issuerClient.Discovery().GenerateQR()
	if err != nil {
		log.Printf("❌ Failed to generate QR code: %v", err)
		return
	}

	qrCode, err := qr.Unicode()
	if err != nil {
		log.Printf("❌ Failed to get QR code: %v", err)
		return
	}

	fmt.Println("📱 QR CODE FOR CONNECTION:")
	fmt.Println("   Scan this with another Self client to establish connection")
	fmt.Println(qrCode)
	fmt.Println("   🔐 Contains cryptographic keys for secure connection")
	fmt.Println("   📱 Compatible with Self mobile apps and other SDK clients")
	fmt.Println()

	// ⏳ CONNECTION WAITING: Wait for peer to scan QR code
	fmt.Println("⏳ Waiting for peer connection (10 seconds)...")
	fmt.Println("   In production, use longer timeouts for user convenience")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	peer, err := qr.WaitForResponse(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("⏰ No connection received (this is normal for a demo)")
			fmt.Println("   In real usage, peers would scan the QR code to connect")
		} else {
			log.Printf("❌ Connection error: %v", err)
		}
		return
	}

	// ✅ CONNECTION SUCCESS: Handle successful peer connection
	fmt.Printf("✅ Connected to peer: %s\n", peer.DID())
	fmt.Println("   🔐 Secure encrypted channel established")
	fmt.Println("   💬 Ready for credential exchange workflows")
	fmt.Println()
}

// printSummary displays a comprehensive educational summary
// This function provides:
// - Summary of demonstrated features
// - Key SDK components used
// - Educational takeaways
// - Next steps for developers
// - Additional learning resources
func printSummary() {
	fmt.Println("🎉 CREDENTIAL ISSUANCE DEMO COMPLETED!")
	fmt.Println("======================================")
	fmt.Println("✅ FEATURES SUCCESSFULLY DEMONSTRATED:")
	fmt.Println("   • Basic credential creation (Email verification)")
	fmt.Println("   • Multi-claim credentials (Profile information)")
	fmt.Println("   • Custom credentials with evidence (Certification with PDF)")
	fmt.Println("   • Complex nested claims (Organization with hierarchical data)")
	fmt.Println("   • Credential builder pattern usage")
	fmt.Println("   • Asset/evidence management and secure storage")
	fmt.Println("   • Verifiable presentation creation")
	fmt.Println("   • Request/response handler configuration")
	fmt.Println()
	fmt.Println("🔧 KEY SDK COMPONENTS UTILIZED:")
	fmt.Println("   • client.NewClient() - Client initialization and configuration")
	fmt.Println("   • NewCredentialBuilder() - Fluent API for credential construction")
	fmt.Println("   • CreateAsset() - Evidence and file attachment management")
	fmt.Println("   • CreatePresentation() - Verifiable presentation packaging")
	fmt.Println("   • OnVerificationRequest/Response() - Event-driven workflows")
	fmt.Println("   • Cryptographic signing and verification (automatic)")
	fmt.Println()
	fmt.Println("📚 EDUCATIONAL TAKEAWAYS:")
	fmt.Println("   • Credentials are cryptographically signed attestations")
	fmt.Println("   • Builder pattern provides clean, readable construction")
	fmt.Println("   • Evidence enhances credential trustworthiness")
	fmt.Println("   • Complex data structures enable rich information modeling")
	fmt.Println("   • Presentations package credentials for selective disclosure")
	fmt.Println("   • Event handlers enable reactive credential workflows")
	fmt.Println()
	fmt.Println("🚀 NEXT STEPS FOR DEVELOPMENT:")
	fmt.Println("   1. Explore credential verification workflows")
	fmt.Println("   2. Implement real peer-to-peer connections")
	fmt.Println("   3. Design custom credential schemas for your use case")
	fmt.Println("   4. Integrate credential workflows into your application")
	fmt.Println("   5. Add business logic for credential validation")
	fmt.Println("   6. Implement selective disclosure and zero-knowledge proofs")
	fmt.Println()
	fmt.Println("📖 ADDITIONAL LEARNING RESOURCES:")
	fmt.Println("   • Self SDK Documentation: https://docs.joinself.com")
	fmt.Println("   • W3C Verifiable Credentials: https://w3.org/TR/vc-data-model/")
	fmt.Println("   • Decentralized Identity: https://identity.foundation")
	fmt.Println("   • Example Applications: /examples directory")
	fmt.Println()
}

// createPresentation demonstrates verifiable presentation creation
// This helper function shows how to:
// - Package credentials into presentations
// - Set presentation types and metadata
// - Prepare credentials for sharing with verifiers
func createPresentation(client *client.Client, cred *credential.VerifiableCredential) (*credential.VerifiablePresentation, error) {
	// 📋 PRESENTATION CREATION: Package credential for sharing
	// Presentations allow selective disclosure of credential information
	return client.Credentials().CreatePresentation(
		[]string{"VerifiablePresentation", "DemoPresentation"}, // Presentation type
		[]*credential.VerifiableCredential{cred},               // Credentials to include
	)
}

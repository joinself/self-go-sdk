// Package main demonstrates complex credential issuance using the Self SDK.
//
// This is the COMPLEX level of credential issuance examples.
// Prerequisites: Complete basic/main.go, multi_claim/main.go, and evidence/main.go first.
//
// This example shows:
// - Complex nested objects in claims
// - Arrays and collections in credentials
// - Hierarchical data organization
// - Real-world organizational data modeling
// - Advanced claim structuring techniques
//
// 🎯 What you'll learn:
// • How to structure complex nested data in credentials
// • Arrays and collections in claims
// • Hierarchical data organization
// • Real-world data modeling patterns
// • Advanced claim structuring
//
// 📚 Next steps:
// • advanced/main.go - All features combined
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joinself/self-go-sdk/client"
	"github.com/joinself/self-go-sdk/credential"
	"github.com/joinself/self-go-sdk/examples/utils"
)

func main() {
	fmt.Println("🎓 Complex Credential Issuance Demo")
	fmt.Println("====================================")
	fmt.Println("This demo shows how to create credentials with complex nested data.")
	fmt.Println("📚 This is the COMPLEX level - advanced data structures.")
	fmt.Println()

	// Step 1: Create issuer and holder clients
	issuer, holder := createClients()
	defer issuer.Close()
	defer holder.Close()

	fmt.Printf("🏢 Issuer: %s\n", issuer.DID())
	fmt.Printf("👤 Holder: %s\n", holder.DID())
	fmt.Println()

	// Step 2: Create credentials with complex data structures
	createOrganizationCredential(issuer, holder)

	fmt.Println("✅ Complex demo completed!")
	fmt.Println()
	fmt.Println("📚 Ready for the next level?")
	fmt.Println("   • Run ../advanced/main.go for all features combined")
	fmt.Println()
}

// createClients sets up the issuer and holder clients
func createClients() (*client.Client, *client.Client) {
	fmt.Println("🔧 Setting up clients...")

	// Create issuer client
	issuer, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("complex_issuer"),
		StoragePath: "./complex_issuer_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create issuer:", err)
	}

	// Create holder client
	holder, err := client.New(client.Config{
		StorageKey:  utils.GenerateStorageKey("complex_holder"),
		StoragePath: "./complex_holder_storage",
		Environment: client.Sandbox,
		LogLevel:    client.LogInfo,
	})
	if err != nil {
		log.Fatal("Failed to create holder:", err)
	}

	fmt.Println("✅ Clients created successfully")
	return issuer, holder
}

// createOrganizationCredential creates an organization credential with complex nested data
func createOrganizationCredential(issuer, holder *client.Client) {
	fmt.Println("🏢 Creating organization credential with complex nested data...")
	fmt.Println("   This demonstrates advanced data structures: nested objects, arrays")
	fmt.Println("   Real-world credentials often contain hierarchical information")
	fmt.Println()

	// Create complex claims structure with nested objects and arrays
	orgCredential, err := issuer.Credentials().NewCredentialBuilder().
		Type(credential.CredentialTypeOrganisation). // Organization credential type
		Subject(holder.DID()).                       // Employee subject
		Issuer(issuer.DID()).                        // Organization issuer
		Claims(map[string]interface{}{               // Complex claims structure
			"organizationName": "TechCorp Inc.", // Company name
			"employeeId":       "EMP-2024-001",  // Employee identifier
			"position": map[string]interface{}{ // Nested position object
				"title":      "Senior Software Engineer", // Job title
				"department": "Engineering",              // Department
				"level":      "L5",                       // Career level
				"startDate":  "2024-01-15",               // Start date
				"manager":    "jane.smith@techcorp.com",  // Manager reference
				"team":       "Backend Infrastructure",   // Team assignment
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
				"address": map[string]interface{}{ // Nested address
					"street":  "123 Tech Street",
					"city":    "San Francisco",
					"state":   "CA",
					"zipCode": "94105",
					"country": "United States",
				},
			},
			"benefits": map[string]interface{}{ // Benefits package
				"healthInsurance": true, // Health coverage
				"retirement401k":  true, // Retirement plan
				"paidTimeOff":     25,   // PTO days
				"stockOptions":    1000, // Stock options
				"remoteWork":      true, // Remote work eligibility
				"wellness": map[string]interface{}{ // Nested wellness benefits
					"gymMembership":    true,
					"mentalHealth":     true,
					"annualWellness":   "$1000",
					"flexibleSchedule": true,
				},
			},
			"certifications": []map[string]interface{}{ // Array of certifications
				{
					"name":       "AWS Solutions Architect", // Certification name
					"level":      "Professional",            // Certification level
					"issueDate":  "2023-06-15",              // Issue date
					"expiryDate": "2026-06-15",              // Expiry date
					"verified":   true,                      // Verification status
					"provider":   "Amazon Web Services",     // Certification provider
				},
				{
					"name":       "Kubernetes Administrator", // Second certification
					"level":      "Certified",                // Certification level
					"issueDate":  "2023-09-20",               // Issue date
					"expiryDate": "2026-09-20",               // Expiry date
					"verified":   true,                       // Verification status
					"provider":   "Cloud Native Computing Foundation",
				},
			},
			"projects": []map[string]interface{}{ // Array of projects
				{
					"name":         "Payment Gateway Redesign",
					"role":         "Lead Developer",
					"startDate":    "2023-01-01",
					"endDate":      "2023-06-30",
					"status":       "Completed",
					"technologies": []string{"Go", "PostgreSQL", "Redis", "Docker"},
				},
				{
					"name":         "Microservices Migration",
					"role":         "Senior Engineer",
					"startDate":    "2023-07-01",
					"endDate":      "2024-01-31",
					"status":       "Completed",
					"technologies": []string{"Go", "Kubernetes", "gRPC", "Prometheus"},
				},
			},
		}).
		ValidFrom(time.Now()).              // Validity start
		SignWith(issuer.DID(), time.Now()). // Cryptographic signature
		Issue(issuer)                       // Issue credential

	if err != nil {
		log.Printf("Failed to create organization credential: %v", err)
		return
	}

	// Display comprehensive organizational information
	fmt.Printf("   ✅ Organization credential created successfully\n")
	fmt.Printf("   🏢 Company: TechCorp Inc.\n")
	fmt.Printf("   💼 Position: Senior Software Engineer (L5)\n")
	fmt.Printf("   🏬 Department: Engineering - Backend Infrastructure\n")
	fmt.Printf("   🆔 Employee ID: EMP-2024-001\n")
	fmt.Printf("   📧 Email: john.doe@techcorp.com\n")
	fmt.Printf("   📍 Office: Building A, Floor 3, Desk 42\n")
	fmt.Printf("   🏠 Address: 123 Tech Street, San Francisco, CA 94105\n")
	fmt.Printf("   🔑 Permissions: 5 access levels\n")
	fmt.Printf("   🎯 Benefits: Health, 401k, 25 PTO days, Stock options, Wellness\n")
	fmt.Printf("   🏆 Certifications: 2 professional certifications\n")
	fmt.Printf("   🚀 Projects: 2 completed projects\n")
	fmt.Printf("   🔒 Type: %v\n", orgCredential.CredentialType())
	fmt.Println()
	fmt.Println("🎓 What happened:")
	fmt.Println("   1. Created credential with deeply nested data structures")
	fmt.Println("   2. Used arrays for multiple values (permissions, certifications, projects)")
	fmt.Println("   3. Nested objects for hierarchical organization (position, contact, benefits)")
	fmt.Println("   4. Mixed data types throughout the structure")
	fmt.Println("   5. Maintained cryptographic integrity for all nested data")
	fmt.Println()
	fmt.Println("📚 Key Learning Points:")
	fmt.Println("   • Credentials can contain complex nested data structures")
	fmt.Println("   • Arrays enable multiple values for single claim types")
	fmt.Println("   • Hierarchical organization mirrors real-world data")
	fmt.Println("   • Complex claims maintain cryptographic integrity")
	fmt.Println("   • Structured data enables precise verification queries")
	fmt.Println("   • Nested objects can contain other nested objects")
	fmt.Println("   • Arrays can contain objects with different structures")
	fmt.Println()
}

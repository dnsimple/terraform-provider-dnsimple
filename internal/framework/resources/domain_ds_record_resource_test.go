package resources_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

// dsRecordData holds the generated DNSSEC DS record data for testing
type dsRecordData struct {
	Algorithm  int
	DigestType int
	KeyTag     int
	Digest     string
	PublicKey  string
}

// generateDsRecordData generates valid DNSSEC DS record data using a randomly generated ECDSA key.
// This ensures each test run uses unique values, avoiding conflicts with leftover data from previous runs.
//
// Go doesn't have built-in support for DNSSEC, we would need to use miekg/dns to simplify the code below.
func generateDsRecordData(domainName string) (*dsRecordData, error) {
	// Generate an ECDSA P-256 key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}

	algorithm := 13
	flags := uint16(257)
	protocol := uint8(3)

	// Encode the ECDSA public key in DNSKEY format (RFC 6605)
	// For P-256: 32 bytes X coordinate + 32 bytes Y coordinate = 64 bytes total
	pubKeyData := make([]byte, 64)
	privateKey.PublicKey.X.FillBytes(pubKeyData[0:32])
	privateKey.PublicKey.Y.FillBytes(pubKeyData[32:64])
	publicKeyBase64 := base64.StdEncoding.EncodeToString(pubKeyData)

	// Build DNSKEY RDATA for key tag and digest calculation
	dnskeyRdata := make([]byte, 4+len(pubKeyData))
	binary.BigEndian.PutUint16(dnskeyRdata[0:2], flags)
	dnskeyRdata[2] = protocol
	dnskeyRdata[3] = byte(algorithm)
	copy(dnskeyRdata[4:], pubKeyData)

	// Calculate key tag (RFC 4034, Appendix B)
	keyTag := calculateKeyTag(dnskeyRdata)

	// Calculate DS digest (SHA-256, digest type 2)
	// DS digest = SHA-256(owner name in wire format | DNSKEY RDATA)
	ownerWireFormat := domainToWireFormat(domainName)
	digestInput := append(ownerWireFormat, dnskeyRdata...)
	digestBytes := sha256.Sum256(digestInput)
	digest := strings.ToUpper(fmt.Sprintf("%x", digestBytes))

	return &dsRecordData{
		Algorithm:  algorithm,
		DigestType: 2, // SHA-256
		KeyTag:     keyTag,
		Digest:     digest,
		PublicKey:  publicKeyBase64,
	}, nil
}

// calculateKeyTag computes the key tag per RFC 4034 Appendix B
func calculateKeyTag(dnskeyRdata []byte) int {
	var ac uint32
	for i, b := range dnskeyRdata {
		if i&1 == 1 {
			ac += uint32(b)
		} else {
			ac += uint32(b) << 8
		}
	}
	ac += ac >> 16 & 0xFFFF
	return int(ac & 0xFFFF)
}

// domainToWireFormat converts a domain name to DNS wire format (RFC 1035)
func domainToWireFormat(domain string) []byte {
	var result []byte
	// Remove trailing dot if present
	domain = strings.TrimSuffix(domain, ".")
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		result = append(result, byte(len(label)))
		result = append(result, []byte(strings.ToLower(label))...)
	}
	result = append(result, 0) // Root label
	return result
}

func TestAccDomainDsRecordResource(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_ds_record.test"

	// Generate unique DNSSEC data for this test run
	dsData, err := generateDsRecordData(domainName)
	if err != nil {
		t.Fatalf("failed to generate DS record data: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainDsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainDsRecordResourceConfig(domainName, dsData),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", strconv.Itoa(dsData.Algorithm)),
					resource.TestCheckResourceAttr(resourceName, "digest", dsData.Digest),
					resource.TestCheckResourceAttr(resourceName, "digest_type", strconv.Itoa(dsData.DigestType)),
					resource.TestCheckResourceAttr(resourceName, "key_tag", strconv.Itoa(dsData.KeyTag)),
					resource.TestCheckResourceAttr(resourceName, "public_key", dsData.PublicKey),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDomainDsRecordImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Updates are a no-op
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDomainDsRecordImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", errors.New("No resource ID set")
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["domain"], rs.Primary.ID), nil
	}
}

func testAccCheckDomainDsRecordResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_ds_record" {
			continue
		}

		dsIdRaw := rs.Primary.Attributes["id"]
		domain := rs.Primary.Attributes["domain"]

		dsId, err := strconv.ParseInt(dsIdRaw, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert domain delegation signer record ID to int: %s", err)
		}

		_, err = dnsimpleClient.Domains.GetDelegationSignerRecord(context.Background(), testAccAccount, domain, dsId)

		if err == nil {
			return fmt.Errorf("domain delegation signer record still exists")
		}
	}
	return nil
}

func testAccDomainDsRecordResourceConfig(domainName string, dsData *dsRecordData) string {
	return fmt.Sprintf(`
resource "dnsimple_ds_record" "test" {
	domain = %[1]q
	algorithm = %[2]q
	digest = %[3]q
	digest_type = %[4]q
	key_tag = %[5]q
	public_key = %[6]q
}`, domainName, strconv.Itoa(dsData.Algorithm), dsData.Digest, strconv.Itoa(dsData.DigestType), strconv.Itoa(dsData.KeyTag), dsData.PublicKey)
}

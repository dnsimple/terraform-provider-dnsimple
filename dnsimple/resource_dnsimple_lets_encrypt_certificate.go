package dnsimple

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSimpleLetsEncryptCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSimpleLetsEncryptCertificateCreate,
		ReadContext:   resourceDNSimpleLetsEncryptCertificateRead,
		UpdateContext: resourceDNSimpleLetsEncryptCertificateUpdate,
		DeleteContext: resourceDNSimpleLetsEncryptCertificateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_id": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "contact_id is deprecated and has no effect. The attribute will be removed in the next major version.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"years": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authority_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceDNSimpleLetsEncryptCertificateDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSimpleLetsEncryptCertificateUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSimpleLetsEncryptCertificateRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = resource.RetryContext(ctx, time.Minute, func() *resource.RetryError {
		provider := meta.(*Client)

		domainID := data.Get("domain_id").(string)
		certificateID, _ := strconv.ParseInt(data.Id(), 10, 64)

		response, err := provider.client.Certificates.GetCertificate(ctx, provider.config.Account, domainID, certificateID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				log.Printf("DNSimple Certificate Not Found - Refreshing from State")
				data.SetId("")
				return resource.NonRetryableError(fmt.Errorf("failed to download DNSimple Certificate: %s", err))
			}
		}

		certificate := response.Data

		if certificate.State != "completed" {
			return resource.RetryableError(fmt.Errorf("certificate order is still %s", errors.New(certificate.State)))
		}

		populateCertificateData(data, certificate)

		return nil
	})
	return nil
}

func resourceDNSimpleLetsEncryptCertificateCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	domainID := data.Get("domain_id").(string)

	certificateAttributes := dnsimple.LetsencryptCertificateAttributes{
		AutoRenew:          data.Get("auto_renew").(bool),
		Name:               data.Get("name").(string),
		SignatureAlgorithm: data.Get("signature_algorithm").(string),
	}

	response, err := provider.client.Certificates.PurchaseLetsencryptCertificate(ctx, provider.config.Account, domainID, certificateAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			return attributeErrorsToDiagnostics(errorResponse)
		}

		return diag.Errorf("Failed to purchase Let's Encrypt Certificate: %s", err)
	}

	certificateID := response.Data.CertificateID

	issueResponse, issueErr := provider.client.Certificates.IssueLetsencryptCertificate(ctx, provider.config.Account, domainID, certificateID)

	if issueErr != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(issueErr, &errorResponse) {
			return attributeErrorsToDiagnostics(errorResponse)
		}

		return diag.Errorf("Failed to issue Let's Encrypt Certificate: %s", issueErr)
	}

	certificate := issueResponse.Data
	populateCertificateData(data, certificate)

	return nil
}

func populateCertificateData(data *schema.ResourceData, certificate *dnsimple.Certificate) {
	data.SetId(strconv.Itoa(int(certificate.ID)))
	data.Set("id", strconv.Itoa(int(certificate.ID)))
	data.Set("domain_id", strconv.Itoa(int(certificate.DomainID)))
	data.Set("contact_id", certificate.ContactID)
	data.Set("years", certificate.Years)
	data.Set("state", certificate.State)
	data.Set("authority_identifier", certificate.AuthorityIdentifier)
	data.Set("auto_renew", certificate.AutoRenew)
	data.Set("created_at", certificate.CreatedAt)
	data.Set("updated_at", certificate.UpdatedAt)
	data.Set("csr", certificate.CertificateRequest)
}

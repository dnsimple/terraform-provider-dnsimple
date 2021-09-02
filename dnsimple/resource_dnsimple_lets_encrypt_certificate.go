package dnsimple

import (
	"context"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSimpleLetsEncryptCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSimpleLetsEncryptCertificateCreate,
		ReadContext: resourceDNSimpleLetsEncryptCertificateRead,
		UpdateContext: resourceDNSimpleLetsEncryptCertificateUpdate,
		DeleteContext: resourceDNSimpleLetsEncryptCertificateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type: schema.TypeString,
				Optional: true,
			},
			"contact_id": {
				Type: schema.TypeInt,
				Required: true,
			},
			"name": {
				Type: schema.TypeString,
				Required: true,
			},
			"alternate_names": {
				Type: schema.TypeList,
				Elem: schema.TypeString,
				Required: true,
			},
			"years": {
				Type: schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type: schema.TypeBool,
				Computed: true,
			},
			"authority_identifier": {
				Type: schema.TypeString,
				Computed: true,
			},
			"auto_renew": {
				Type: schema.TypeBool,
				Required: true,
			},
			"created_at": {
				Type: schema.TypeString,
				Computed:  true,
			},
			"updated_at": {
				Type: schema.TypeString,
				Computed:  true,
			},
			"expires_on": {
				Type: schema.TypeString,
				Computed:  true,
			},
			"csr": {
				Type: schema.TypeString,
				Computed:  true,
			},

		},
	}
}

func resourceDNSimpleLetsEncryptCertificateDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSimpleLetsEncryptCertificateUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

func resourceDNSimpleLetsEncryptCertificateRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	domainID := data.Id()
	certificateID := data.Get("certificate_id")

	response, err := provider.client.Certificates.DownloadCertificate(context.Background(), provider.config.Account, domainID, int64(certificateID.(int)))

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Certificate Not Found - Refreshing from State")
			data.SetId("")
			return nil
		}
		return diag.Errorf("Failed to download DNSimple Certificate: %s", err)
	}

	certificate := response.Data
	data.Set("server", certificate.ServerCertificate)
	data.Set("root", certificate.RootCertificate)
	data.Set("chain", certificate.IntermediateCertificates)

	return nil
}

func resourceDNSimpleLetsEncryptCertificateCreate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	domainID := data.Get("domain").(string)

	certificateAttributes := dnsimple.LetsencryptCertificateAttributes{
		ContactID: int64(data.Get("contact_id").(int)),
		AutoRenew: data.Get("auto_renew").(bool),
		Name: data.Get("name").(string),
		AlternateNames: data.Get("alternative_names").([]string),
	}

	response, err := provider.client.Certificates.PurchaseLetsencryptCertificate(context.Background(), provider.config.Account, domainID, certificateAttributes)

	if err != nil {
		return diag.Errorf("Failed to purchase Let's Encrypt Certificate: %s", err)
	}

	certificateID := response.Data.CertificateID

	issueResponse, issueErr := provider.client.Certificates.IssueLetsencryptCertificate(context.Background(), provider.config.Account, domainID, certificateID)

	if issueErr != nil {
		return diag.Errorf("Failed to issue Let's Encrypt Certificate: %s", issueErr)
	}

	certificate := issueResponse.Data
	data.Set("id", certificate.ID)
	data.Set("domain_id", certificate.DomainID)
	data.Set("contact_id", certificate.ContactID)
	data.Set("alternative_names", append(certificate.AlternateNames, certificate.CommonName))
	data.Set("years", certificate.Years)
	data.Set("state", certificate.State)
	data.Set("authority_identifier", certificate.AuthorityIdentifier)
	data.Set("auto_renew", certificate.AutoRenew)
	data.Set("created_at", certificate.CreatedAt)
	data.Set("updated_at", certificate.UpdatedAt)
	data.Set("expires_on", certificate.ExpiresOn)
	data.Set("csr", certificate.CertificateRequest)

	return nil
}

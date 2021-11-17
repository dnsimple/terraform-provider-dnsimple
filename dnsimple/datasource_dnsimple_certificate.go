package dnsimple

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNSimpleCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSimpleCertificateRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"server_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"root_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_chain": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSimpleCertificateRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	certificateId, err := strconv.Atoi(data.Get("certificate_id").(string))
	if err != nil {
		return diag.Errorf("Error converting Certificate ID: %s", err)
	}

	cert, err := provider.client.Certificates.DownloadCertificate(context.Background(), provider.config.Account, data.Get("domain").(string), int64(certificateId))

	if err != nil {
		return diag.Errorf("Couldn't find DNSimple SSL Certificate: %s", err)
	}

	certificate := cert.Data
	data.Set("server_certificate", certificate.ServerCertificate)
	data.Set("root_certificate", certificate.RootCertificate)
	data.Set("certificate_chain", certificate.IntermediateCertificates)

	key, err := provider.client.Certificates.GetCertificatePrivateKey(context.Background(), provider.config.Account, data.Get("domain").(string), int64(certificateId))

	data.Set("private_key", key.Data.PrivateKey)
	data.SetId(time.Now().UTC().String())

	return nil
}

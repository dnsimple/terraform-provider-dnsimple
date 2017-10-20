package dnsimple

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDNSimpleCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSimpleCertificateRead,
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

func dataSourceDNSimpleCertificateRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	certificate_id, err := strconv.Atoi(d.Get("certificate_id").(string))
	if err != nil {
		return fmt.Errorf("Error converting Certificate ID: %s", err)
	}

	cert, err := provider.client.Certificates.DownloadCertificate(provider.config.Account, d.Get("domain").(string), certificate_id)

	if err != nil {
		return fmt.Errorf("Couldn't find DNSimple SSL Certificate: %s", err)
	}

	certificate := cert.Data
	d.Set("server_certificate", certificate.ServerCertificate)
	d.Set("root_certificate", certificate.RootCertificate)
	d.Set("certificate_chain", certificate.IntermediateCertificates)

	key, err := provider.client.Certificates.GetCertificatePrivateKey(provider.config.Account, d.Get("domain").(string), certificate_id)

	d.Set("private_key", key.Data.PrivateKey)
	d.SetId(time.Now().UTC().String())

	return nil
}

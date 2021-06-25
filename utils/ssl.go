package utils

import (
	"github.com/mcoo/requests"
	"time"
)

type SSLResult struct {
	Code int `json:"code"`
	Data struct {
		Version   string   `json:"version"`
		Hosts     []string `json:"hosts"`
		CheckHost string   `json:"check_host"`
		Status    struct {
			Suggests []struct {
				Code  int    `json:"code"`
				Tip   string `json:"tip"`
				Link  string `json:"link"`
				Level int    `json:"level"`
			} `json:"suggests"`
			ProtocolDetail []struct {
				Name    string `json:"name"`
				Support bool   `json:"support"`
			} `json:"protocol_detail"`
			Protocols []struct {
				Name    string `json:"name"`
				Support bool   `json:"support"`
			} `json:"protocols"`
			Certs struct {
				Rsas []struct {
					LeafCertInfo struct {
						Hash               string    `json:"hash"`
						CertStatus         int       `json:"cert_status"`
						CertStatusText     string    `json:"cert_status_text"`
						CommonName         string    `json:"common_name"`
						PublickeyAlgorithm string    `json:"publickey_algorithm"`
						Issuer             string    `json:"issuer"`
						SignatureAlgorithm string    `json:"signature_algorithm"`
						Organization       string    `json:"organization"`
						OrganizationUnit   string    `json:"organization_unit"`
						Sans               []string  `json:"sans"`
						Transparency       string    `json:"transparency"`
						IsCtQualified      int       `json:"is_ct_qualified"`
						CertType           string    `json:"cert_type"`
						CertDomainType     int       `json:"cert_domain_type"`
						BrandName          string    `json:"brand_name"`
						ValidFrom          time.Time `json:"valid_from"`
						ValidTo            time.Time `json:"valid_to"`
						IsSni              bool      `json:"is_sni"`
						OcspMustStaple     bool      `json:"ocsp_must_staple"`
						OcspURL            []string  `json:"ocsp_url"`
						OcspStatus         int       `json:"ocsp_status"`
						Country            string    `json:"country"`
						Locality           string    `json:"locality"`
						StreetAddress      string    `json:"street_address"`
						PostalCode         string    `json:"postal_code"`
					} `json:"leaf_cert_info"`
					CertsFormServer struct {
						ProvidedNumber int `json:"provided_number"`
						Certs          []struct {
							Order            int       `json:"order"`
							CommonName       string    `json:"common_name"`
							Hash             string    `json:"hash"`
							Pin              string    `json:"pin"`
							ValidTo          time.Time `json:"valid_to"`
							IsExpired        bool      `json:"is_expired"`
							KeyAlgo          string    `json:"key_algo"`
							SignAlgo         string    `json:"sign_algo"`
							IssuerCommonName string    `json:"issuer_common_name"`
						} `json:"certs"`
					} `json:"certs_form_server"`
					Chain struct {
						Certs []struct {
							Type               string    `json:"type"`
							Missing            bool      `json:"missing"`
							FromServer         bool      `json:"from_server"`
							CommonName         string    `json:"common_name"`
							PublickeyAlgorithm string    `json:"publickey_algorithm"`
							SignatureAlgorithm string    `json:"signature_algorithm"`
							Sha1               string    `json:"sha1"`
							Pin                string    `json:"pin"`
							ExpiresIn          int       `json:"expires_in"`
							Issuer             string    `json:"issuer"`
							BeginTime          time.Time `json:"begin_time"`
							EndTime            time.Time `json:"end_time"`
							Order              int       `json:"order,omitempty"`
							IsCa               bool      `json:"is_ca"`
						} `json:"certs"`
						HasRoot    bool   `json:"has_root"`
						HasOtherCa bool   `json:"has_other_ca"`
						MissCa     bool   `json:"miss_ca"`
						ID         string `json:"id"`
					} `json:"chain"`
				} `json:"rsas"`
				Eccs []interface{} `json:"eccs"`
				Sigs []interface{} `json:"sigs"`
				Encs []interface{} `json:"encs"`
			} `json:"certs"`
			Ciphers interface{} `json:"ciphers"`
		} `json:"status"`
	} `json:"data"`
}

func SSLStatus(host string) (r SSLResult, e error) {
	res, e := requests.Get("https://myssl.com/api/v1/ssl_status?port=443&c=0&domain=" + host)
	if e != nil {
		return
	}
	e = res.Json(&r)
	if e != nil {
		return
	}
	return
}

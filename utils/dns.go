package utils

import "github.com/mcoo/requests"

type DNSResult struct {
	Code int `json:"code"`
	Data struct {
		Num86 []struct {
			Answer struct {
				TimeConsume string `json:"time_consume"`
				Records     []struct {
					TTL        int    `json:"ttl"`
					Value      string `json:"value"`
					IPLocation string `json:"ip_location"`
				} `json:"records"`
				Error string `json:"error"`
			} `json:"answer"`
		} `json:"86"`
		Num852 []struct {
			Answer struct {
				TimeConsume string `json:"time_consume"`
				Records     []struct {
					TTL        int    `json:"ttl"`
					Value      string `json:"value"`
					IPLocation string `json:"ip_location"`
				} `json:"records"`
				Error string `json:"error"`
			} `json:"answer"`
		} `json:"852"`
		Num01 []struct {
			Answer struct {
				TimeConsume string `json:"time_consume"`
				Records     []struct {
					TTL        int    `json:"ttl"`
					Value      string `json:"value"`
					IPLocation string `json:"ip_location"`
				} `json:"records"`
				Error string `json:"error"`
			} `json:"answer"`
		} `json:"01"`
	} `json:"data"`
}

func DnsQuery(host string) (r DNSResult, e error) {
	res, e := requests.Get("https://myssl.com/api/v1/tools/dns_query?qtype=1&qmode=-1&host=" + host)
	if e != nil {
		return
	}
	e = res.Json(&r)
	if e != nil {
		return
	}
	return
}

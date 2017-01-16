package rule

import (
	"strings"
	"github.com/miekg/dns"
	"github.com/miekg/dns/dnsutil"
)

const (
	CN_GROUP = "CN"
	SYSTEM_GROUP = "SYSTEM"
	REJECT_GROUP = "REJECT"

	SCHEME_DOMAIN_SUFFIX = "DOMAIN-SUFFIX"
	SCHEME_DOMAIN_MATCH = "DOMAIN"
	SCHEME_DOMAIN_KEYWORD = "DOMAIN-KEYWORD"
)

type DomainRules []DomainRule

func (self DomainRules) FindGroup(domain string) string {

	if dns.IsFqdn(domain) {
		domain = dnsutil.TrimDomainName(domain, ".")
	}

	for _, rule := range self {
		if rule.Match(domain) {
			return rule.Group
		}
	}
	return ""
}

type DomainRule struct {
	MatchType string `toml:"scheme"`
	Group     string `toml:"group"`
	Values    []string `toml:"value"`
}

func (rule DomainRule) Match(input string) bool {

	for _, value := range rule.Values {
		switch rule.MatchType {
		case SCHEME_DOMAIN_MATCH:
			if input == value {
				return true
			}
		case SCHEME_DOMAIN_SUFFIX:
			if strings.HasSuffix(input, value) {
				return true
			}
		case SCHEME_DOMAIN_KEYWORD:
			if strings.Contains(input, value) {
				return true
			}
		default:
			continue
		}
	}

	return false
}




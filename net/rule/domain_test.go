package rule

import (
	"testing"
)

func TestDomainRules(t *testing.T) {
	rule := DomainRule{Group:CN_GROUP,
		MatchType:SCHEME_DOMAIN_SUFFIX,
		Values:[]string{"baidu.com", "xiaomi.com"},
	}

	r := rule.Match("baidu.com")
	if r != true {
		t.Fail()
	}

	r = rule.Match("xiaomi.com")
	if r != true {
		t.Fail()
	}

	r = rule.Match("www.google.com")
	if r != false {
		t.Fail()
	}

	r = rule.Match("google.com")
	if r != false {
		t.Fail()
	}

	rules := DomainRules{rule}

	g := rules.FindGroup("baidu.com")
	if g != CN_GROUP {
		t.Fatal("rules match fail!")
	}

}


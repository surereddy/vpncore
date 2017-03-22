package dns

import (
	"github.com/athom/goset"
	"github.com/miekg/dns"
	"github.com/FTwOoO/go-logger"
	"github.com/FTwOoO/vpncore/net/geoip"
	"github.com/FTwOoO/vpncore/net/rule"
	"github.com/FTwOoO/vpncore/net/addr"
	"time"
	"net"
	"sync"
)

const (
	dnsDefaultPort = 53
	dnsDefaultTtl = 600
	dnsDefaultPacketSize = 4096
	dnsDefaultReadTimeout = 5
	dnsDefaultWriteTimeout = 5
)

type DNSServer struct {
	config    *Config
	geoReader *geoip.Reader
	logger    *logger.Logger
	cache     *DnsCache
	client    *dns.Client
}

func createNormalDnsServer(addr string) (*dns.Server, error) {
	server := &dns.Server{
		Net:          "udp",
		Addr:         addr,
		UDPSize:      dnsDefaultPacketSize,
		ReadTimeout:  time.Duration(dnsDefaultReadTimeout) * time.Second,
		WriteTimeout: time.Duration(dnsDefaultWriteTimeout) * time.Second,
	}

	return server, nil
}

func NewDNSServer(config *Config, dontServ bool) (ds *DNSServer, err error) {
	if config == nil {
		config = DefaultConfig
	} else {
		config.DNSGroups[rule.SYSTEM_GROUP] = DefaultConfig.DNSGroups[rule.SYSTEM_GROUP]
	}

	var cache *DnsCache
	if config.DNSCache.Enable {
		cache = NewDnsCache(config.DNSCache.MaxCount)
	}

	client := &dns.Client{
		Net:          "udp",
		UDPSize:      dnsDefaultPacketSize,
		ReadTimeout:  time.Duration(dnsDefaultReadTimeout) * time.Second,
		WriteTimeout: time.Duration(dnsDefaultWriteTimeout) * time.Second,
	}

	mlogger, err := logger.NewLogger(config.LogConfig.LogFile, config.LogConfig.LogLevel)
	if err != nil {
		return
	}

	var reader *geoip.Reader
	if config.GeoIPValidate.Enable && config.GeoIPValidate.GeoIpDBPath != "" {
		reader, err = geoip.Open(config.GeoIPValidate.GeoIpDBPath)
		if err != nil {
			return nil, err
		}
	}

	ds = &DNSServer{
		config:     config,
		geoReader:  reader,
		logger:     mlogger,
		cache:      cache,
		client:     client,
	}

	if dontServ == false {
		inDnsServ, _ := createNormalDnsServer(config.Addr)
		inDnsServ.Handler = ds
		go func() {
			err = inDnsServ.ListenAndServe()
			panic(err)
		}()
	}

	return
}

func (d *DNSServer) QueryIPv4(host string, ctx QueryContext, h MsgHandler) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), dns.TypeA)
	msg.RecursionDesired = true

	w := &sessionWriter{handler:h, ctx:ctx}
	d.ServeDNS(w, msg)
	return
}

func (d *DNSServer) QueryByDNSMsg(msg *dns.Msg, ctx QueryContext, h Handler) {
	w := &sessionWriter{handler:h, ctx:ctx}
	d.ServeDNS(w, msg)
}

func (d *DNSServer) QueryByData(requestData []byte, ctx QueryContext, h Handler) {

	w := &sessionWriter{handler:h, ctx:ctx}

	req := new(dns.Msg)
	err := req.Unpack(requestData)
	if err != nil {
		x := new(dns.Msg)
		x.SetRcodeFormatError(req)
		w.WriteMsg(x)
		w.Close()
	}

	d.ServeDNS(w, req)
}


// Main callback for miekg/dns. Collects information about the query,
// constructs a response, and returns it to the connector.
func (ds *DNSServer) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {

	if len(req.Question) == 0 {
		dns.HandleFailed(w, req)
		return
	}

	var resp *dns.Msg
	hitCache := false

	question := req.Question[0]
	if ds.cache != nil && question.Qclass == dns.ClassINET &&
		(question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA) {

		if cacheResp, ok := ds.cache.Get(question.Name); ok {
			ds.logger.Debugf("%s hit cache", question.Name)
			// dont change cache object, copy it
			newResp := *cacheResp
			newResp.Id = req.Id
			resp = &newResp
			hitCache = true
		} else if cacheResp, ok := ds.config.Hosts.Get(req); ok {
			ds.logger.Debugf("%s found in hosts file", question.Name)
			newResp := *cacheResp
			newResp.Id = req.Id
			resp = &newResp
		}
	}

	if resp == nil {
		group := ds.config.DomainRules.FindGroup(question.Name)
		if group == rule.REJECT_GROUP {
			ds.logger.Debugf("Reject %s!", question.Name)
			dns.HandleFailed(w, req)
			return
		} else if group != "" {
			resp = ds.sendRequest(req, []string{group})
		} else {
			resp = ds.sendRequest(req, ds.config.DefaultGroups)
		}
	}

	if resp != nil {
		w.WriteMsg(resp)

		if question.Qclass == dns.ClassINET &&
			(question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA) &&
			len(resp.Answer) > 0 {

			if hitCache == false {
				ttl := time.Duration(resp.Answer[0].Header().Ttl)
				ds.logger.Debugf("Insert %s into cache with ttl:%d", question.Name, ttl)
				ds.cache.Add(question.Name, resp, ttl)

			}
		}
	}
}

type DNSResult struct {
	Group    string
	DnsIp    net.IP

	Response *dns.Msg
	Rtt      time.Duration
	Err      error
}

func (ds *DNSServer) sendRequest(req *dns.Msg, dnsgroups []string) (resp *dns.Msg) {

	chanLen := 0
	for _, group := range dnsgroups {
		chanLen += len(ds.config.DNSGroups[group])
	}

	// the chan is big enough to hold all results
	results := make(chan DNSResult, chanLen)
	ds.sendDNSRequestsAsync(req, results, dnsgroups)

	WaitingDNSResponse:
	for result := range results {
		if result.Err != nil {
			ds.logger.Errorf("Error from group[%s] DNS[%s]: ===>\n %v \n<===\n", result.Group, result.DnsIp.String(), result.Err)
			continue
		} else {
			ds.logger.Debugf("Result from group[%s] DNS[%s]: ===>\n %v \n<===\n", result.Group, result.DnsIp.String(), result.Response)
		}

		if result.Response.Rcode == dns.RcodeServerFailure {
			ds.logger.Errorf("Resolve on group [%s:%s] failed: code %d", result.Group, result.DnsIp.String(), result.Response.Rcode)
			continue
		}

		if len(result.Response.Answer) < 1 {
			ds.logger.Debugf("0 answer response from  %s", result.DnsIp.String())
			continue
		}

		isThisResultOk := true

		ScanResponse: for _, as := range result.Response.Answer {
			switch as.(type) {
			case *dns.A:
				aRecord, _ := as.(*dns.A)
				resultIp := aRecord.A

				if !ds.isIpOK(result.Group, resultIp) {
					isThisResultOk = false
					break ScanResponse
				}
			default:
				continue ScanResponse
			}
		}

		if isThisResultOk {
			resp = result.Response
			break WaitingDNSResponse
		}
	}

	ds.logger.Debugf("response for request:%v\n", resp)
	return resp
}

func (ds *DNSServer) isIpOK(dnsGroup string, resultIp net.IP) bool {
	if ds.config.IPBlocker.FindIP(resultIp) {
		ds.logger.Infof("block ip %v", resultIp.String())
		return false
	}

	if ds.config.GeoIPValidate.Enable == true {
		country, err := ds.geoReader.Country(resultIp)

		if err != nil {
			ds.logger.Errorf("cant reconize the localtion of ip:%s", resultIp.String())
			return false
		}

		groupMatched := goset.IsIncluded(ds.config.GeoIPValidate.Groups, dnsGroup)
		countryMatched := string(ds.config.GeoIPValidate.GeoCountry) == country.Country.IsoCode

		if groupMatched  && countryMatched {
			ds.logger.Debugf("DNS result IP[%s] in Geo country[%s] from DNS server[%s] can be trusted!", resultIp, country.Country.IsoCode, dnsGroup)
			return true
		} else if !groupMatched  && !countryMatched {
			ds.logger.Debugf("DNS result IP[%s] in Geo country[%s] from DNS server[%s] can be trusted!", resultIp, country.Country.IsoCode, dnsGroup)
			return true
		}

		return false
	} else {
		return true
	}
}

func (ds *DNSServer) sendDNSRequestsAsync(req *dns.Msg, results chan <- DNSResult, dnsgroups []string) {
	var wg sync.WaitGroup

	for _, group := range dnsgroups {
		dnsL := ds.config.DNSGroups[group]
		wg.Add(len(dnsL))

		for _, dnsAddr := range dnsL {
			go func(group string, dnsAddr addr.DNSAddresss) {
				defer wg.Done()

				//c := &dns.Client{Net: "udp", Timeout:10 * time.Second}
				resp, rtt, err := ds.client.Exchange(req, dnsAddr.String())
				results <- DNSResult{Response:resp, Rtt: rtt, Err: err, Group:group, DnsIp:dnsAddr.Ip}

			}(group, dnsAddr)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

}
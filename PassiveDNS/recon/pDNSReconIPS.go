package recon

import (
	"fmt"
	"github.com/IncSW/geoip2"
	"net"
	"time"
)


func GetGeoInformationIP(ip string) (uint32,string,string,string) {

	reader, err := geoip2.NewCityReaderFromFile("/home/jbcui/go/src/PTwork/recon/GeoLite2-City.mmdb")
	if err != nil {
		fmt.Errorf("some error %v", err)

	}

	record, err := reader.Lookup(net.ParseIP(ip))
	if err != nil {
		fmt.Errorf("some error %v", err)
	}

	readerA, err := geoip2.NewASNReaderFromFile("/home/jbcui/go/src/PTwork/recon/GeoLite2-ASN.mmdb")
	if err != nil {
		fmt.Errorf("some error %v", err)

	}
	recordA,err:=readerA.Lookup(net.ParseIP(ip))
	if err != nil {
		fmt.Errorf("some error %v", err)
	}
	asn:=record.Country.GeoNameID
	asnOrg:=recordA.AutonomousSystemOrganization
	geoIP:=record.Country.ISOCode
	firstSeen:=time.Now().Format("2006-1-2")
	return asn,asnOrg,geoIP,firstSeen
}



package recon

import (
	"PTwork/models"
	"fmt"
	"github.com/IncSW/geoip2"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//Get host tld. Example "com", "asia" etc.
func GetTLD(host string) (tld string) {

	t:=strings.Split(host,".")

	for _,i:=range t[1:]{
		tld = i
	}

	return tld
}

//Get hosts IPs. Example "1.1.1.1", "8.8.8.8" etc.
func GetIP(host string) []string{

	IpArray:=make([]string,0,2)
	ip,err:=net.LookupIP(host)

	if err!=nil{
		fmt.Errorf("some error %v", err)
	}

	for _,j:=range ip{
		IpArray = append(IpArray, j.String())
	}

	return IpArray
}

//Open databases collection information from host to read ASN, GeoIP, ASNOrg.
func OpenGeoLite()(Reader *geoip2.CityReader, ReaderA *geoip2.ASNReader){

	curPath,_:=os.Getwd()
	Reader,err:= geoip2.NewCityReaderFromFile(curPath+"/recon/GeoLite2-City.mmdb")

	if err!=nil{
		fmt.Errorf("Something went wrong when open GeoLite-City %v\n",err)
	}

	ReaderA,err=geoip2.NewASNReaderFromFile(curPath+"/recon/GeoLite2-ASN.mmdb")

	if err!=nil{
		fmt.Errorf("Something went wrong when open GeoLite-ASN %v\n",err)
	}

	return Reader,ReaderA
}


//https://github.com/IncSW/geoip2 docs
//https://dev.maxmind.com/geoip/geoip2/downloadable/
//Get Geo Information about host.
func GetGeoInformation(ip string, reader *geoip2.CityReader,readerA *geoip2.ASNReader) (asn uint32,asnOrg string,geoIP string,firstSeen string, lastSeen string) {

	record, err := reader.Lookup(net.ParseIP(ip))
	if err != nil {
		return 0,"undefined","undefined",time.Now().Format("2006-1-2"),time.Now().Format("2006-1-2")
	}

	recordA,err:=readerA.Lookup(net.ParseIP(ip))
	if err != nil {
		return 0,"undefined","undefined",time.Now().Format("2006-1-2"),time.Now().Format("2006-1-2")
	}

	asn=record.Country.GeoNameID
	asnOrg=recordA.AutonomousSystemOrganization
	geoIP=record.Country.ISOCode
	firstSeen=time.Now().Format("2006-1-2")
	lastSeen = time.Now().Format("2006-1-2")

	return asn,asnOrg,geoIP,firstSeen,lastSeen
}

//Collect WhoIS information about host.
func GetWhoIS(host string) (whoIs models.WhoISInformation){

	defer func() {
		if r:=recover(); r!=nil{
			fmt.Println("Panic happend in Whois block",r)
		}
	}()


	query,err:=whois.Whois(host)
	if err!=nil{
		panic("Could not get query "+ err.Error())
	}

	if !strings.Contains(query,"No match for"){
		result,err:=whoisparser.Parse(query)
		if err!=nil{
			panic("Could not parse whois query "+err.Error())
		}


		if result.Registrant !=nil {
			whoIs.RegistrarId=result.Registrar.ID
			whoIs.RegistrarName=result.Registrar.Name
		}

		if result.Registrant != nil{
			whoIs.RegistrantName = result.Registrant.Name
			whoIs.RegistrantCompany=result.Registrant.Organization
			whoIs.RegistrantAddress=result.Registrant.Street
			whoIs.RegistrantCity=result.Registrant.City
			whoIs.RegistrantState = result.Registrant.Province
			whoIs.RegistrantZip=result.Registrant.PostalCode
			whoIs.RegistrantCountry=result.Registrant.Country
			whoIs.RegistrantEmail=result.Registrant.Email
			whoIs.RegistrantPhone=result.Registrant.Phone
			whoIs.RegistrantFax=result.Registrant.Fax
		}

		if result.Domain !=nil{
			whoIs.NameServers=result.Domain.NameServers
			whoIs.CreateDate=result.Domain.CreatedDate
			whoIs.UpdateDate=result.Domain.UpdatedDate
			whoIs.ExpireDate=result.Domain.ExpirationDate
		}

		whoIs.FirstSeen=time.Now().Format("2006-1-2")
		whoIs.LastSeen=time.Now().Format("2006-1-2")

		return whoIs

	}else{
		whoIs.FirstSeen=time.Now().Format("2006-1-2")
		whoIs.LastSeen=time.Now().Format("2006-1-2")
	}

	time.Sleep(10*time.Millisecond)
	return whoIs
}




func GetRegisterTime() string{

	registerTime:=time.Now().Format("2006-1-2")
	return registerTime

}



//Get information about hosts from domain registrant
func GetInformationAboutHost(hosts <-chan string,dataOut chan<- models.DomainData,wg *sync.WaitGroup, reader *geoip2.CityReader,readerA *geoip2.ASNReader)  {

	defer wg.Done()

	var data  models.DomainData
	for host:=range hosts{

		//Get tld from host
		data.Tld = GetTLD(host)

		//Get ips slice
		ips:=GetIP(host)
		ipv4:=make([]string,0)
		for _,ip:=range ips{
			if net.ParseIP(ip).To4()!=nil {
				ipv4=append(ipv4,ip)
			}
		}

		//Get Resolvers
		resolvers :=make([]models.Resolvers,0)
		for _,ip:=range ipv4{
			asn,asnOrg,geoIP,firstSeen,lastSeen:=GetGeoInformation(ip, reader, readerA)
			resolvers = append(resolvers, models.Resolvers{
				IP:        ip,
				ASN:       asn,
				ASNOrg:    asnOrg,
				GeoIP:     geoIP,
				FirstSeen: firstSeen,
				LastSeen:  lastSeen,
			})
		}
		data.Resolvers= resolvers

		//Get host status
		if len(data.Resolvers) == 0{
			data.Status = "expired"
		}else{
			data.Status = "active"
		}

		data.FirstSeen=GetRegisterTime()
		data.LastSeen=GetRegisterTime()
		data.Host=host

		//Get WhoISInformation
		whoISInformation:=make([]models.WhoISInformation,0)
		whoIs:=GetWhoIS(host)

		whoISInformation=append(whoISInformation,whoIs)
		data.WhoISInformation=whoISInformation

		//Send collect data to channel
		dataOut <- data
	}
	close(dataOut)

}


//Get information about hosts from domain registrant
func GetInformationAboutHostTranco(hosts <-chan string,dataOut chan<- models.TrancoDomainData,wg *sync.WaitGroup, reader *geoip2.CityReader,readerA *geoip2.ASNReader)  {

	defer wg.Done()

	var data  models.TrancoDomainData
	data.Identifier = "tranco_top"

	for host:=range hosts{

		//Get tld for tranco list
		data.Tld = GetTLD(host)

		//Get ips slice
		ips:=GetIP(host)
		ipv4:=make([]string,0)
		for _,ip:=range ips{

			if net.ParseIP(ip).To4()!=nil {

				ipv4=append(ipv4,ip)
			}
		}

		data.FirstSeen=GetRegisterTime()
		data.LastSeen=GetRegisterTime()
		data.Host=host

		//Get WhoISInformation
		whoISInformation:=make([]models.WhoISInformation,0)
		whos:=GetWhoIS(host)

		whoISInformation=append(whoISInformation,whos)
		data.WhoISInformation=whoISInformation

		//Get resolvers slice
		resolvers :=make([]models.Resolvers,0)
		for _,ip:=range ipv4{

			asn,asnOrg,geoIP,firstSeen,lastSeen:=GetGeoInformation(ip, reader, readerA)
			resolvers = append(resolvers, models.Resolvers{

				IP:        ip,
				ASN:       asn,
				ASNOrg:    asnOrg,
				GeoIP:     geoIP,
				FirstSeen: firstSeen,
				LastSeen:  lastSeen,

			})
		}

		data.Resolvers= resolvers

		//Get host status
		if len(data.Resolvers) == 0{

			data.Status = "expired"

		}else{

			data.Status = "active"

		}

		//Send collect data to channel
		dataOut <- data
	}

	close(dataOut)

}


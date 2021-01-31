package models

import (
	"time"
)


//PassiveDNS IPs structure to collect data from gathering ips from hosts:
type IPSData struct {

	IP 		string 		`json:"ip"`
	ASN 	uint32 		`json:"asn"`
	ASNOrg 	string 		`json:"asn_org"`
	GeoIP 	string 		`json:"geo"`
	Domains []Domains 	`json:"domains"`

}

//Substructure for PassiveDNS IPs:
type Domains struct {

	Domain 		string 		`json:"domain"`
	FirstSeen 	time.Time 	`json:"first_seen"`
	LastSeen 	time.Time 	`json:"last_seen"`

}

//Main page data pars structures:
type MainPages 		[]struct {

	DomainAmount 	[]string 	`json:"domain_amount"`
	DomainZones 	string 		`json:"domain_zone"`
	Links 			[]string 	`json:"links"`

}

//Tranco pDNS domain info:
type TrancoDomainData struct {

	Host             string             `json:"domain"`
	Tld              string             `json:"tld"`
	Status           string             `json:"status"`
	Identifier       string             `json:"identifier"`
	Subdomains       []string           `json:"subdomains,omitempty"`
	FirstSeen        string             `json:"first_seen"`
	LastSeen         string             `json:"last_seen"`
	Resolvers        []Resolvers        `json:"resolvers,omitempty"`
	WhoISInformation []WhoISInformation `json:"whois_information,omitempty"`
	Sources          []string           `json:"sources,omitempty"`
	Tags             []string           `json:"tags,omitempty"`

}

//Domain data structure for gathering information from scrapping hosts from domain-status.com:
type DomainData struct {

	Host 				string 				`json:"domain"`
	Tld 				string 				`json:"tld"`
	Status 				string 				`json:"status"`
	Subdomains 			[]string 			`json:"subdomains,omitempty"`
	FirstSeen 			string 				`json:"first_seen"`
	LastSeen 			string 				`json:"last_seen"`
	Resolvers 			[]Resolvers 		`json:"resolvers,omitempty"`
	WhoISInformation 	[]WhoISInformation 	`json:"whois_information,omitempty"`
	Source 				string				`json:"source"`
	Tag					[]string			`json:"tag"`

}

//Substructure from DomainData struct to collect hosts resolvers:
type Resolvers struct{

	IP 			string 		`json:"ip"`
	ASN 		uint32 		`json:"asn,omitempty"`
	ASNOrg 		string 		`json:"asn_org,omitempty"`
	GeoIP 		string 		`json:"geo,omitempty"`
	FirstSeen 	string 		`json:"first_seen"`
	LastSeen 	string 		`json:"last_seen"`

}


//Substructure to collect whois information about scrapping hosts:
type WhoISInformation struct {

	FirstSeen 			string 			`json:"first_seen"`
	LastSeen 			string 			`json:"last_seen"`
	RegistrarId 		string 			`json:"registrar_id,omitempty"`
	RegistrarName 		string 			`json:"registrar_name,omitempty"`
	RegistrantName 		string 			`json:"registrant_name,omitempty"`
	RegistrantCompany 	string 			`json:"registrant_company,omitempty"`
	RegistrantAddress 	string 			`json:"registrant_address,omitempty"`
	RegistrantCity 		string 			`json:"registrant_city,omitempty"`
	RegistrantState 	string 			`json:"registrant_state,omitempty"`
	RegistrantZip 		string 			`json:"registrant_zip,omitempty"`
	RegistrantCountry 	string 			`json:"registrant_country,omitempty"`
	RegistrantEmail 	string 			`json:"registrant_email,omitempty"`
	RegistrantPhone 	string 			`json:"registrant_phone,omitempty"`
	RegistrantFax 		string 			`json:"registrant_fax,omitempty"`
	CreateDate 			string 			`json:"create_date,omitempty"`
	UpdateDate 			string 			`json:"update_date,omitempty"`
	ExpireDate 			string 			`json:"expire_date,omitempty"`
	NameServers 		[]string 		`json:"name_servers,omitempty"`

}

//Kafka message structure to send queue message:
type KafkaMessage struct {

	Source 		string 		`json:"source"`
	Tag 		[]string 	`json:"tag"`

}
package gufodao

import "strings"

func maskemail(email string) string {
	mailarr := strings.Split(email, "@")
	domain := mailarr[1]
	milbody := mailarr[0]
	domainarr := strings.Split(domain, ".")
	domainbody := domainarr[0]
	mailbodymask := milbody[0:1] + "***"
	domainbodymask := domainbody[0:1] + "***" + domainbody[len(domainbody)-1:]
	maskedemail := mailbodymask + "@" + domainbodymask + "." + domainarr[1]
	return maskedemail
}

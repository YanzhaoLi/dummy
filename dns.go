package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"os"
	flags "github.com/jessevdk/go-flags"
)


func isHostInDNS(ctx context.Context, host string, dnsNames []string) (ok bool, err error) {
	fmt.Println("---Try reverse lookup---")
	// reverse lookup
	wildcards, names := []string{}, []string{}
	for _, dns := range dnsNames {
		if strings.HasPrefix(dns, "*.") {
			wildcards = append(wildcards, dns[1:])
		} else {
			names = append(names, dns)
		}
	}
	fmt.Println("INFO: wildcards: ", names)

	lnames, lerr := net.DefaultResolver.LookupAddr(ctx, host)
	fmt.Println("INFO: reverse-dns Result of " + host + " :", lnames)
	for _, name := range lnames {
		// strip trailing '.' from PTR record
		if name[len(name)-1] == '.' {
			name = name[:len(name)-1]
		}
		fmt.Println("INFO: result of dns name: ", name)
		for _, wc := range wildcards {
			fmt.Println("INFO: one wildcard: ", wc)
			if strings.HasSuffix(name, wc) {
				return true, nil
			}
		}
		for _, n := range names {
			if n == name {
				return true, nil
			}
		}
	}
	err = lerr

	fmt.Println("---Try forward lookup---")
	// forward lookup
	for _, dns := range names {
		addrs, lerr := net.DefaultResolver.LookupHost(ctx, dns)
		fmt.Println("INFO: forward-dns result of " + dns + " :", addrs)
		if lerr != nil {
			err = lerr
			fmt.Println("ERROR: ", lerr.Error())
			continue
		}
		for _, addr := range addrs {
			if addr == host {
				return true, nil
			}
		}
	}
	return false, err
}



func main() {
	var opts struct {
		// Example of a required flag
		Host string `long:"host" description:"Host to nslookup" required:"true"`
		// Example of a slice of strings
		DnsNames []string `long:"dns-names" description:"A slice of dns names to match" required:"true"`

	}

	// Parse flags from `args'. Note that here we use flags.ParseArgs for
	// the sake of making a working example. Normally, you would simply use
	// flags.Parse(&opts) which uses os.Args
	// args, err := flags.ParseArgs(&opts, args)
	if _, err := flags.Parse(&opts); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	host := opts.Host
	dnsNames := opts.DnsNames
	ok, err := isHostInDNS(context.Background(), host, dnsNames)
	if ok {
		fmt.Println("Yes Found")
	} else {
		errStr := ""
		if err != nil {
			errStr = " (" + err.Error() + ")"
		}
		fmt.Println("Not Found", errStr)
	}
}


package main

import (
	"flag"
	"log"
	"net"
	"net/netip"
	"os"
	"slices"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	mmdbv2 "github.com/oschwald/maxminddb-golang/v2"
)

var (
	countries     string
	input         string
	output        string
	countryFilter []string
)

func init() {

	const (
		usage = "... -h"
	)
	flag.StringVar(&countries, "countries", "VN", usage)
	flag.StringVar(&input, "input", "", usage)
	flag.StringVar(&output, "output", "", usage)
	flag.Parse()
	countryFilter = []string{}
	for _, c := range strings.Split(countries, ",") {
		countryFilter = append(countryFilter, c)
	}
	if len(countryFilter) == 0 {
		log.Fatal("at least one country is required")
	}
	if input == "" {
		log.Fatal("input file is required")
	}
	if output == "" {
		log.Fatal("output file is required")
	}
}

func prefixToIPNet(p netip.Prefix) *net.IPNet {
	addr := p.Addr()
	ip := net.IP(addr.AsSlice())
	ones := p.Bits()
	bits := addr.BitLen() // 32 for IPv4, 128 for IPv6
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(ones, bits),
	}
}

func main() {
	// 1) Open source database
	reader, err := mmdbv2.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// 2) Create a new writer tree
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType: "GeoLite2-Country",
		RecordSize:   24,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Struct used only for quick filtering
	var filter struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	kept := 0

	// 3) Iterate networks
	totally := 0
	for result := range reader.Networks() {
		if err := result.Err(); err != nil {
			log.Fatal(err)
		}

		// A) decode minimal fields to decide keep/drop
		filter.Country.ISOCode = ""
		if err := result.Decode(&filter); err != nil {
			log.Fatal(err)
		}
		totally++
		if !slices.Contains(countryFilter, filter.Country.ISOCode) {
			continue
		}

		// B) decode full record into mmdbwriter's types
		var record mmdbtype.Map
		if err := result.Decode(&record); err != nil {
			log.Fatal(err)
		}

		// C) insert into new DB
		ipnet := prefixToIPNet(result.Prefix())
		if err := writer.Insert(ipnet, record); err != nil {
			log.Printf("insert failed for %s: %v", result.Prefix(), err)
			continue
		}

		kept++
	}

	// 4) Write output once

	fh, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	if _, err := writer.WriteTo(fh); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully created %s (records kept: %d out of %d)", output, kept, totally)
}

package main

import (
	"testing"
)

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"192.168.1.1", true},            // Valid IPv4
		{"255.255.255.255", true},        // Valid IPv4
		{"0.0.0.0", true},                // Valid IPv4
		{"::1", true},                    // Valid IPv6 (loopback)
		{"2001:db8::ff00:42:8329", true}, // Valid IPv6
		{"invalid-ip", false},            // Invalid IP
		{"256.256.256.256", false},       // Out of range
		{"1234:5678::12345", false},      // Invalid IPv6
		{"", false},                      // Empty input
	}

	for _, test := range tests {
		result := isValidIP(test.ip)
		if result != test.expected {
			t.Errorf("isValidIP(%q) = %v; expected %v", test.ip, result, test.expected)
		}
	}
}

func TestIsValidDNSName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"example.com", true},     // Valid
		{"sub.example.com", true}, // Valid subdomain
		{"example123.com", true},  // Valid with numbers
		{"-example.com", false},   // Invalid: starts with hyphen
		{"example-.com", false},   // Invalid: ends with hyphen
		{"exa_mple.com", false},   // Invalid: underscore not allowed
		{"example", false},        // Invalid: no TLD
		{"123.com", true},         // Valid numeric domain
		{"", false},               // Invalid: empty input
		{"toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong.com", false}, // Invalid: label too long
	}

	for _, test := range tests {
		result := isValidDNSName(test.name)
		if result != test.expected {
			t.Errorf("isValidDNSName(%q) = %v; expected %v", test.name, result, test.expected)
		}
	}
}

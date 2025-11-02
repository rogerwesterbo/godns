package validation

import (
	"testing"
)

func TestValidateDNSAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		// Valid addresses
		{
			name:    "valid :53",
			address: ":53",
			wantErr: false,
		},
		{
			name:    "valid 0.0.0.0:53",
			address: "0.0.0.0:53",
			wantErr: false,
		},
		{
			name:    "valid 192.168.1.1:53",
			address: "192.168.1.1:53",
			wantErr: false,
		},
		{
			name:    "valid 127.0.0.1:5353",
			address: "127.0.0.1:5353",
			wantErr: false,
		},
		{
			name:    "valid IPv6 [::]:53",
			address: "[::]:53",
			wantErr: false,
		},
		{
			name:    "valid IPv6 [2001:db8::1]:53",
			address: "[2001:db8::1]:53",
			wantErr: false,
		},
		{
			name:    "valid hostname with port",
			address: "localhost:53",
			wantErr: false,
		},
		{
			name:    "valid high port :8053",
			address: ":8053",
			wantErr: false,
		},

		// Invalid addresses
		{
			name:    "invalid empty string",
			address: "",
			wantErr: true,
		},
		{
			name:    "invalid port only 53",
			address: "53",
			wantErr: true,
		},
		{
			name:    "invalid IP without port",
			address: "192.168.1.1",
			wantErr: true,
		},
		{
			name:    "invalid colon without port",
			address: ":",
			wantErr: true,
		},
		{
			name:    "invalid port too high",
			address: ":99999",
			wantErr: true,
		},
		{
			name:    "invalid port zero",
			address: ":0",
			wantErr: true,
		},
		{
			name:    "invalid negative port",
			address: ":-1",
			wantErr: true,
		},
		{
			name:    "invalid non-numeric port",
			address: ":abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDNSAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDNSAddress(%q) error = %v, wantErr %v", tt.address, err, tt.wantErr)
			}
		})
	}
}

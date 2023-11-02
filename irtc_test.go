package irtc

import (
	"net"
	"testing"
)

func TestLowerEqual(t *testing.T) {
	type args struct {
		a string
		b string
	}

	testCases := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "lower",
			args: args{
				a: "1.1.1.0",
				b: "1.1.1.1",
			},
			want: true,
		},
		{
			name: "equal",
			args: args{
				a: "1.1.1.1",
				b: "1.1.1.1",
			},
			want: true,
		},
		{
			name: "upper",
			args: args{
				a: "1.1.1.2",
				b: "1.1.1.1",
			},
			want: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			aIP := net.ParseIP(tt.args.a)
			if aIP == nil {
				t.Fatal(&net.ParseError{
					Type: "IP address",
					Text: tt.args.a,
				})
			}
			bIP := net.ParseIP(tt.args.b)
			if bIP == nil {
				t.Fatal(&net.ParseError{
					Type: "IP address",
					Text: tt.args.b,
				})
			}

			got := lowerEqual(aIP, bIP)
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestConvertIPRangeToCIDRs(t *testing.T) {
	type args struct {
		begin string
		end   string
	}

	testCases := []struct {
		name    string
		args    args
		want    []string
		wantErr error
	}{
		{
			name: "happy for IPv4 single result",
			args: args{
				begin: "192.168.1.0",
				end:   "192.168.1.255",
			},
			want: []string{
				"192.168.1.0/24",
			},
			wantErr: nil,
		},
		{
			name: "happy for IPv4 special single result",
			args: args{
				begin: "192.168.1.111",
				end:   "192.168.1.111",
			},
			want: []string{
				"192.168.1.111/32",
			},
			wantErr: nil,
		},
		{
			name: "happy for IPv4 multi result",
			args: args{
				begin: "192.168.1.3",
				end:   "192.168.1.254",
			},
			want: []string{
				"192.168.1.3/32",
				"192.168.1.4/30",
				"192.168.1.8/29",
				"192.168.1.16/28",
				"192.168.1.32/27",
				"192.168.1.64/26",
				"192.168.1.128/26",
				"192.168.1.192/27",
				"192.168.1.224/28",
				"192.168.1.240/29",
				"192.168.1.248/30",
				"192.168.1.252/31",
				"192.168.1.254/32",
			},
			wantErr: nil,
		},
		{
			name: "end parse error",
			args: args{
				begin: "192.255.1.1",
				end:   "192.268.1.289",
			},
			want: nil,
			wantErr: &net.ParseError{
				Type: "IP address",
				Text: "192.268.1.289",
			},
		},
		{
			name: "begin parse error",
			args: args{
				begin: "192.268.1.288",
				end:   "192.268.1.289",
			},
			want: nil,
			wantErr: &net.ParseError{
				Type: "IP address",
				Text: "192.268.1.288",
			},
		},
		{
			name: "end > begin",
			args: args{
				begin: "198.168.1.2",
				end:   "192.168.1.1",
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "happy for IPv6 multi result",
			args: args{
				begin: "2408:874f:2000:100::fff",
				end:   "2408:874f:2000:1ff:ffff:ffff:ffff:ffff",
			},
			want: []string{
				"2408:874f:2000:100::fff/128",
				"2408:874f:2000:100::1000/116",
				"2408:874f:2000:100::2000/115",
				"2408:874f:2000:100::4000/114",
				"2408:874f:2000:100::8000/113",
				"2408:874f:2000:100::1:0/112",
				"2408:874f:2000:100::2:0/111",
				"2408:874f:2000:100::4:0/110",
				"2408:874f:2000:100::8:0/109",
				"2408:874f:2000:100::10:0/108",
				"2408:874f:2000:100::20:0/107",
				"2408:874f:2000:100::40:0/106",
				"2408:874f:2000:100::80:0/105",
				"2408:874f:2000:100::100:0/104",
				"2408:874f:2000:100::200:0/103",
				"2408:874f:2000:100::400:0/102",
				"2408:874f:2000:100::800:0/101",
				"2408:874f:2000:100::1000:0/100",
				"2408:874f:2000:100::2000:0/99",
				"2408:874f:2000:100::4000:0/98",
				"2408:874f:2000:100::8000:0/97",
				"2408:874f:2000:100:0:1::/96",
				"2408:874f:2000:100:0:2::/95",
				"2408:874f:2000:100:0:4::/94",
				"2408:874f:2000:100:0:8::/93",
				"2408:874f:2000:100:0:10::/92",
				"2408:874f:2000:100:0:20::/91",
				"2408:874f:2000:100:0:40::/90",
				"2408:874f:2000:100:0:80::/89",
				"2408:874f:2000:100:0:100::/88",
				"2408:874f:2000:100:0:200::/87",
				"2408:874f:2000:100:0:400::/86",
				"2408:874f:2000:100:0:800::/85",
				"2408:874f:2000:100:0:1000::/84",
				"2408:874f:2000:100:0:2000::/83",
				"2408:874f:2000:100:0:4000::/82",
				"2408:874f:2000:100:0:8000::/81",
				"2408:874f:2000:100:1::/80",
				"2408:874f:2000:100:2::/79",
				"2408:874f:2000:100:4::/78",
				"2408:874f:2000:100:8::/77",
				"2408:874f:2000:100:10::/76",
				"2408:874f:2000:100:20::/75",
				"2408:874f:2000:100:40::/74",
				"2408:874f:2000:100:80::/73",
				"2408:874f:2000:100:100::/72",
				"2408:874f:2000:100:200::/71",
				"2408:874f:2000:100:400::/70",
				"2408:874f:2000:100:800::/69",
				"2408:874f:2000:100:1000::/68",
				"2408:874f:2000:100:2000::/67",
				"2408:874f:2000:100:4000::/66",
				"2408:874f:2000:100:8000::/65",
				"2408:874f:2000:101::/64",
				"2408:874f:2000:102::/63",
				"2408:874f:2000:104::/62",
				"2408:874f:2000:108::/61",
				"2408:874f:2000:110::/60",
				"2408:874f:2000:120::/59",
				"2408:874f:2000:140::/58",
				"2408:874f:2000:180::/57",
			},
		},
	}

	for index, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertIPRangeToCIDRs(tt.args.begin, tt.args.end)
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("want err: %v, got err: %v", tt.wantErr, err)
				}
			}
			for i, cidr := range tt.want {
				if len(got) <= i || cidr != got[i] {
					t.Errorf("run testCases[%d] failed", index)
				}
			}
		})
	}
}

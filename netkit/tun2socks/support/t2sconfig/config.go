package t2sconfig

import (
	"github.com/wsrf16/swiss/sugar/base/ipkit"
	"github.com/wsrf16/swiss/sugar/base/pointkit"
	"runtime"
	"strings"
	"time"
)

type TunConfig struct {
	TunName         *string
	TunAddress      *string
	TunGateway      *string
	TunGWMetric     *int
	TunMask         *string
	TunDNS          *string
	TunPersist      *bool
	TunMode         *string
	MTU             *int
	BlockOutsideDNS *bool
	PublicOnly      *bool
	EnableDNSCache  *bool
	ProxyType       *string
	ProxyServer     *string
	DefaultGateway  *string
	UDPTimeout      *time.Duration
	DNSFallback     *bool
}

const (
	TunModeWinTun     = "wintun"
	TunModeWinSyscall = "syscall"
	TunModeWater      = "water"
)

// var Config *TunConfig

func (config *TunConfig) GetTunDNS() []string {
	dnsServers := strings.Split(*config.TunDNS, ",")
	return dnsServers
}

func (config *TunConfig) GetProxyHost() string {
	proxyHost := strings.Split(*config.ProxyServer, ":")[0]
	return proxyHost
}

func (config *TunConfig) FillWithDefaultValue() {
	if pointkit.IsEmpty(config.TunName) {
		tunName := "MYNIC"
		config.TunName = &tunName
	}
	if pointkit.IsEmpty(config.UDPTimeout) {
		udpTimeout := 1 * time.Minute
		config.UDPTimeout = &udpTimeout
	}
	if pointkit.IsEmpty(config.TunAddress) {
		tunAddr := "10.255.0.4"
		config.TunAddress = &tunAddr
	}
	if pointkit.IsEmpty(config.TunGateway) {
		// tunGw := "10.255.0.1"
		tunGw, _ := ipkit.GetIPv4Gateway(*config.TunAddress)
		config.TunGateway = &tunGw
	}
	if config.TunGWMetric == nil {
		tunGWMetric := 8
		config.TunGWMetric = &tunGWMetric
	}
	if pointkit.IsEmpty(config.TunMask) {
		tunMask := "255.255.255.0"
		config.TunMask = &tunMask
	}
	if pointkit.IsEmpty(config.TunDNS) {
		tunDNS := "8.8.8.8,8.8.4.4"
		config.TunDNS = &tunDNS
	}
	if pointkit.IsEmpty(config.TunMode) {
		if runtime.GOOS == "windows" {
			tunMode := TunModeWinTun
			config.TunMode = &tunMode
		} else {
			tunMode := TunModeWater
			config.TunMode = &tunMode
		}
	}
	if pointkit.IsEmpty(config.MTU) {
		mtu := 1500
		config.MTU = &mtu
	}
	if pointkit.IsEmpty(config.ProxyType) {
		proxyType := "socks"
		config.ProxyType = &proxyType
	}
	if pointkit.IsEmpty(config.DefaultGateway) {
		defaultGateway := "192.168.0.1"
		config.DefaultGateway = &defaultGateway
	}
	if pointkit.IsEmpty(config.TunPersist) {
		tunPersist := false
		config.TunPersist = &tunPersist
	}
	if pointkit.IsEmpty(config.PublicOnly) {
		publicOnly := false
		config.PublicOnly = &publicOnly
	}
	if pointkit.IsEmpty(config.EnableDNSCache) {
		enableDNSCache := false
		config.EnableDNSCache = &enableDNSCache
	}
	if pointkit.IsEmpty(config.BlockOutsideDNS) {
		blockOutsideDNS := false
		config.BlockOutsideDNS = &blockOutsideDNS
	}

}

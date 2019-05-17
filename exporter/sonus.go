package exporter

import (
	"encoding/xml"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/teliax/sonus-metrics-exporter/config"
)

// Exporter is used to store Metrics data and embeds the config struct.
// This is done so that the relevant functions have easy access to the
// user defined runtime configuration when the Collect method is called.
type Exporter struct {
	APIMetrics map[string]*prometheus.Desc
	config.Config
}

type TGcollection struct {
	XMLName                  xml.Name                    `xml:"collection,omitempty" json:"collection,omitempty"`
	AttrXmlnsy               string                      `xml:"xmlns y,attr"  json:",omitempty"`
	TGglobalTrunkGroupStatus []*TGglobalTrunkGroupStatus `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 globalTrunkGroupStatus,omitempty" json:"globalTrunkGroupStatus,omitempty"`
}

type TGaddressContext struct {
	XMLName xml.Name `xml:"addressContext,omitempty" json:"addressContext,omitempty"`
	Value   string   `xml:",chardata" json:",omitempty"`
}

type TGbwAvailable struct {
	XMLName xml.Name `xml:"bwAvailable,omitempty" json:"bwAvailable,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGbwCurrentLimit struct {
	XMLName xml.Name `xml:"bwCurrentLimit,omitempty" json:"bwCurrentLimit,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGbwInboundUsage struct {
	XMLName xml.Name `xml:"bwInboundUsage,omitempty" json:"bwInboundUsage,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGbwOutboundUsage struct {
	XMLName xml.Name `xml:"bwOutboundUsage,omitempty" json:"bwOutboundUsage,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGglobalTrunkGroupStatus struct {
	XMLName                      xml.Name                      `xml:"globalTrunkGroupStatus,omitempty" json:"globalTrunkGroupStatus,omitempty"`
	Attrxmlns                    string                        `xml:"xmlns,attr"  json:",omitempty"`
	TGaddressContext             *TGaddressContext             `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 addressContext,omitempty" json:"addressContext,omitempty"`
	TGbwAvailable                *TGbwAvailable                `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwAvailable,omitempty" json:"bwAvailable,omitempty"`
	TGbwCurrentLimit             *TGbwCurrentLimit             `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwCurrentLimit,omitempty" json:"bwCurrentLimit,omitempty"`
	TGbwInboundUsage             *TGbwInboundUsage             `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwInboundUsage,omitempty" json:"bwInboundUsage,omitempty"`
	TGbwOutboundUsage            *TGbwOutboundUsage            `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 bwOutboundUsage,omitempty" json:"bwOutboundUsage,omitempty"`
	TGinboundCallsUsage          *TGinboundCallsUsage          `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 inboundCallsUsage,omitempty" json:"inboundCallsUsage,omitempty"`
	TGname                       *TGname                       `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 name,omitempty" json:"name,omitempty"`
	TGoutboundCallsUsage         *TGoutboundCallsUsage         `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 outboundCallsUsage,omitempty" json:"outboundCallsUsage,omitempty"`
	TGpacketOutDetectState       *TGpacketOutDetectState       `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 packetOutDetectState,omitempty" json:"packetOutDetectState,omitempty"`
	TGpriorityBwUsage            *TGpriorityBwUsage            `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 priorityBwUsage,omitempty" json:"priorityBwUsage,omitempty"`
	TGpriorityCallUsage          *TGpriorityCallUsage          `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 priorityCallUsage,omitempty" json:"priorityCallUsage,omitempty"`
	TGstate                      *TGstate                      `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 state,omitempty" json:"state,omitempty"`
	TGtotalCallsAvailable        *TGtotalCallsAvailable        `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsAvailable,omitempty" json:"totalCallsAvailable,omitempty"`
	TGtotalCallsConfigured       *TGtotalCallsConfigured       `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsConfigured,omitempty" json:"totalCallsConfigured,omitempty"`
	TGtotalCallsInboundReserved  *TGtotalCallsInboundReserved  `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalCallsInboundReserved,omitempty" json:"totalCallsInboundReserved,omitempty"`
	TGtotalOutboundCallsReserved *TGtotalOutboundCallsReserved `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 totalOutboundCallsReserved,omitempty" json:"totalOutboundCallsReserved,omitempty"`
	TGzone                       *TGzone                       `xml:"http://sonusnet.com/ns/mibs/SONUS-GLOBAL-TRUNKGROUP/1.0 zone,omitempty" json:"zone,omitempty"`
}

type TGinboundCallsUsage struct {
	XMLName xml.Name `xml:"inboundCallsUsage,omitempty" json:"inboundCallsUsage,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGname struct {
	XMLName xml.Name `xml:"name,omitempty" json:"name,omitempty"`
	Value   string   `xml:",chardata" json:",omitempty"`
}

type TGoutboundCallsUsage struct {
	XMLName xml.Name `xml:"outboundCallsUsage,omitempty" json:"outboundCallsUsage,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGpacketOutDetectState struct {
	XMLName xml.Name `xml:"packetOutDetectState,omitempty" json:"packetOutDetectState,omitempty"`
	Value   string   `xml:",chardata" json:",omitempty"`
}

type TGpriorityBwUsage struct {
	XMLName xml.Name `xml:"priorityBwUsage,omitempty" json:"priorityBwUsage,omitempty"`
	Flag    bool     `xml:",chardata" json:",omitempty"`
}

type TGpriorityCallUsage struct {
	XMLName xml.Name `xml:"priorityCallUsage,omitempty" json:"priorityCallUsage,omitempty"`
	Flag    bool     `xml:",chardata" json:",omitempty"`
}

type TGstate struct {
	XMLName xml.Name `xml:"state,omitempty" json:"state,omitempty"`
	Value   string   `xml:",chardata" json:",omitempty"`
}

type TGtotalCallsAvailable struct {
	XMLName xml.Name `xml:"totalCallsAvailable,omitempty" json:"totalCallsAvailable,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGtotalCallsConfigured struct {
	XMLName xml.Name `xml:"totalCallsConfigured,omitempty" json:"totalCallsConfigured,omitempty"`
	Number  float64  `xml:",chardata" json:",omitempty"`
}

type TGtotalCallsInboundReserved struct {
	XMLName xml.Name `xml:"totalCallsInboundReserved,omitempty" json:"totalCallsInboundReserved,omitempty"`
	Flag    bool     `xml:",chardata" json:",omitempty"`
}

type TGtotalOutboundCallsReserved struct {
	XMLName xml.Name `xml:"totalOutboundCallsReserved,omitempty" json:"totalOutboundCallsReserved,omitempty"`
	Flag    bool     `xml:",chardata" json:",omitempty"`
}

type TGzone struct {
	XMLName xml.Name `xml:"zone,omitempty" json:"zone,omitempty"`
	Value   string   `xml:",chardata" json:",omitempty"`
}

// Response struct is used to store http.Response and associated data
type Response struct {
	url      string
	response *http.Response
	body     []byte
	err      error
}

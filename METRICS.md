# Metrics

Below are an example of the metrics as exposed by this exporter.

## DSP Statistics

```
# HELP sonus_dsp_codec_utilization Codec utilization, in percent
# TYPE sonus_dsp_codec_utilization gauge
sonus_dsp_codec_utilization{codec="G.711",system="sbc01"} 1
sonus_dsp_codec_utilization{codec="G.711 V8",system="sbc01"} 1
sonus_dsp_codec_utilization{codec="G.722",system="sbc01"} 0

# HELP sonus_dsp_compression_utilization Compression resource utilization, in percent
# TYPE sonus_dsp_compression_utilization gauge
sonus_dsp_compression_utilization{system="sbc01"} 1

# HELP sonus_dsp_resources_total Total compression resources
# TYPE sonus_dsp_resources_total gauge
sonus_dsp_resources_total{system="sbc01"} 50500

# HELP sonus_dsp_resources_used Usage of DSP resources per slot
# TYPE sonus_dsp_resources_used gauge
sonus_dsp_resources_used{slot="1",system="sbc01"} 120
```

## Fan Status

```
# HELP sonus_fan_speed Current speed of fans, in RPM
# TYPE sonus_fan_speed gauge
sonus_fan_speed{fanID="FAN1/BOT",server="sbc01a"} 5632
sonus_fan_speed{fanID="FAN1/BOT",server="sbc01b"} 5632
```

## IP Interface Statistics

```
# HELP sonus_ipinterface_media_streams Number of media streams currently on ipInterfaceGroup
# TYPE sonus_ipinterface_media_streams gauge
sonus_ipinterface_media_streams{name="IPINT1"} 0
sonus_ipinterface_media_streams{name="IPINT2"} 98

# HELP sonus_ipinterface_rxbandwidth Receive bandwidth in use on interface, in bytes per second
# TYPE sonus_ipinterface_rxbandwidth gauge
sonus_ipinterface_rxbandwidth{name="IPINT1"} 0
sonus_ipinterface_rxbandwidth{name="IPINT2"} 4523

# HELP sonus_ipinterface_rxpackets Number of packets received on ipInterfaceGroup
# TYPE sonus_ipinterface_rxpackets counter
sonus_ipinterface_rxpackets{name="IPINT1"} 1.7744832e+07
sonus_ipinterface_rxpacket2{name="IPINT2"} 3.817481033e+09

# HELP sonus_ipinterface_status Current status of ipInterfaceGroup
# TYPE sonus_ipinterface_status gauge
sonus_ipinterface_status{name="IPINT1",status_text="resAllocated"} 0
sonus_ipinterface_status{name="IPINT2",status_text="resAllocated"} 0

# HELP sonus_ipinterface_txbandwidth Transmit bandwidth in use on interface, in bytes per second
# TYPE sonus_ipinterface_txbandwidth gauge
sonus_ipinterface_txbandwidth{name="IPINT1"} 0
sonus_ipinterface_txbandwidth{name="IPINT2"} 3424

# HELP sonus_ipinterface_txpackets Number of packets transmitted on ipInterfaceGroup
# TYPE sonus_ipinterface_txpackets counter
sonus_ipinterface_txpackets{name="IPINT1"} 2.0886046e+07
sonus_ipinterface_txpackets{name="IPINT2"} 2.300479202e+09
```

## Power Supplies

```
# HELP sonus_powersupply_powerfault Is there a power fault, per supply
# TYPE sonus_powersupply_powerfault gauge
sonus_powersupply_powerfault{powerSupplyID="PSA",server="sbc01a"} 0
sonus_powersupply_powerfault{powerSupplyID="PSA",server="sbc01b"} 0
sonus_powersupply_powerfault{powerSupplyID="PSB",server="sbc01a"} 0
sonus_powersupply_powerfault{powerSupplyID="PSB",server="sbc01b"} 0

# HELP sonus_powersupply_voltagefault Is there a voltage fault, per supply
# TYPE sonus_powersupply_voltagefault gauge
sonus_powersupply_voltagefault{powerSupplyID="PSA",server="sbc01a"} 0
sonus_powersupply_voltagefault{powerSupplyID="PSA",server="sbc01b"} 0
sonus_powersupply_voltagefault{powerSupplyID="PSB",server="sbc01a"} 0
sonus_powersupply_voltagefault{powerSupplyID="PSB",server="sbc01b"} 0
```

## Server Status

```
# HELP sonus_system_redundancy_role Current role of server. 1 = active
# TYPE sonus_system_redundancy_role gauge
sonus_system_redundancy_role{role_name="active",server="sbc01a"} 1
sonus_system_redundancy_role{role_name="standby",server="sbc01b"} 0

# HELP sonus_system_sync_status Current synchronization status. 1 = syncCompleted
# TYPE sonus_system_sync_status gauge
sonus_system_sync_status{server="sbc01a",status_name="syncCompleted"} 1
sonus_system_sync_status{server="sbc01b",status_name="syncCompleted"} 1

# HELP sonus_system_uptime Current uptime of server, in seconds
# TYPE sonus_system_uptime counter
sonus_system_uptime{server="sbc01a",type="application"} 23400
sonus_system_uptime{server="sbc01a",type="os"} 23499
sonus_system_uptime{server="sbc01b",type="application"} 568
sonus_system_uptime{server="sbc01b",type="os"} 954
```

## SIP ARS Status

```
# HELP sonus_sipars_endpoint_status State of a sipArs monitored endpoint
# TYPE sonus_sipars_endpoint_status gauge
sonus_sipars_endpoint_status{endpoint_address="52.11.22.33",endpoint_port="5060",state_name="blacklisted",zone="ZONE1"} 1
```

## Trunk Groups

```
# HELP sonus_TG_bytes Bandwidth in use by current calls
# TYPE sonus_TG_bytes gauge
sonus_TG_bytes{direction="inbound",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_bytes{direction="inbound",name="ZONE1-OUT-TG",zone="ZONE1"} 0
sonus_TG_bytes{direction="outbound",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_bytes{direction="outbound",name="ZONE1-OUT-TG",zone="ZONE1"} 0

# HELP sonus_TG_outbound_state State of outbound calls on the trunkgroup
# TYPE sonus_TG_outbound_state gauge
sonus_TG_outbound_state{name="ZONE1-IN-TG",zone="ZONE1"} 1
sonus_TG_outbound_state{name="ZONE1-OUT-TG",zone="ZONE1"} 1

# HELP sonus_TG_sip_req_recv Number of SIP requests received
# TYPE sonus_TG_sip_req_recv counter
sonus_TG_sip_req_recv{method="BYE",name="ZONE1-IN-TG",zone="ZONE1"} 40
sonus_TG_sip_req_recv{method="BYE",name="ZONE1-OUT-TG",zone="ZONE1"} 30

# HELP sonus_TG_sip_req_sent Number of SIP requests sent
# TYPE sonus_TG_sip_req_sent counter
sonus_TG_sip_req_sent{method="BYE",name="ZONE1-IN-TG",zone="ZONE1"} 30
sonus_TG_sip_req_sent{method="BYE",name="ZONE1-OUT-TG",zone="ZONE1"} 40

# HELP sonus_TG_sip_resp_recv Number of SIP responses received
# TYPE sonus_TG_sip_resp_recv counter
sonus_TG_sip_resp_recv{code="18x",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_sip_resp_recv{code="18x",name="ZONE1-OUT-TG",zone="ZONE1"} 0

# HELP sonus_TG_sip_resp_sent Number of SIP responses sent
# TYPE sonus_TG_sip_resp_sent counter
sonus_TG_sip_resp_sent{code="18x",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_sip_resp_sent{code="18x",name="ZONE1-OUT-TG",zone="ZONE1"} 0

# HELP sonus_TG_state State of the trunkgroup
# TYPE sonus_TG_state gauge
sonus_TG_state{name="ZONE1-IN-TG",zone="ZONE1"} 1
sonus_TG_state{name="ZONE1-OUT-TG",zone="ZONE1"} 1

# HELP sonus_TG_total_channels Number of configured channels
# TYPE sonus_TG_total_channels gauge
sonus_TG_total_channels{name="ZONE1-IN-TG",zone="ZONE1"} 1000
sonus_TG_total_channels{name="ZONE1-OUT-TG",zone="ZONE1"} 15000

# HELP sonus_TG_usage_total Number of active calls
# TYPE sonus_TG_usage_total gauge
sonus_TG_usage_total{direction="inbound",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_usage_total{direction="inbound",name="ZONE1-OUT-TG",zone="ZONE1"} 0
sonus_TG_usage_total{direction="outbound",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_usage_total{direction="outbound",name="ZONE1-OUT-TG",zone="ZONE1"} 0
```

## Exporter Metric Disposition

```
# HELP sonus_exporter_metric_disposition Number of times each metric has succeeded or failed being collected
# TYPE sonus_exporter_metric_disposition counter
sonus_exporter_metric_disposition{name="DSP",successful="true"} 2
sonus_exporter_metric_disposition{name="Fan",successful="true"} 2
sonus_exporter_metric_disposition{name="IPInterface",successful="true"} 4
sonus_exporter_metric_disposition{name="PowerSupply",successful="true"} 2
sonus_exporter_metric_disposition{name="SipStatistic",successful="true"} 10
sonus_exporter_metric_disposition{name="TrunkGroup",successful="true"} 2
```

## Exporter Metric Duration

```
sonus_exporter_metric_duration{name="IPInterface",stage="http",quantile="0.5"} 0.071724695
sonus_exporter_metric_duration{name="IPInterface",stage="http",quantile="0.9"} 0.071964337
sonus_exporter_metric_duration{name="IPInterface",stage="http",quantile="0.99"} 0.071964337
sonus_exporter_metric_duration_sum{name="IPInterface",stage="http"} 0.425327406
sonus_exporter_metric_duration_count{name="IPInterface",stage="http"} 6
sonus_exporter_metric_duration{name="IPInterface",stage="process",quantile="0.5"} 0.000705656
sonus_exporter_metric_duration{name="IPInterface",stage="process",quantile="0.9"} 0.00397619
sonus_exporter_metric_duration{name="IPInterface",stage="process",quantile="0.99"} 0.00397619
```

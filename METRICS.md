# Metrics

Below are an example of the metrics as exposed by this exporter.

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
# HELP sonus_TG_state State of the trunkgroup
# TYPE sonus_TG_state gauge
sonus_TG_state{name="ZONE1-IN-TG",zone="ZONE1"} 1
sonus_TG_state{name="ZONE1-OUT-TG",zone="ZONE1"} 1
# HELP sonus_TG_usage_total Number of active calls
# TYPE sonus_TG_usage_total gauge
sonus_TG_usage_total{direction="inbound",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_usage_total{direction="inbound",name="ZONE1-OUT-TG",zone="ZONE1"} 0
sonus_TG_usage_total{direction="outbound",name="ZONE1-IN-TG",zone="ZONE1"} 0
sonus_TG_usage_total{direction="outbound",name="ZONE1-OUT-TG",zone="ZONE1"} 0```

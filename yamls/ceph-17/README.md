# 每个 Ceph Monitor节点:
kubectl label node <nodename> ceph-mon-17=enabled
# 每个 OSD node节点: DEVx如 vdb sdb
kubectl label node <nodename> ceph-osd-17=enabled ceph-osd-dev-{DEV1}=enabled ceph-osd-dev-{DEV2}=enabled

# 部署osd节点
helm template osd-devices --debug | kubectl apply -f -

# Ceph monitor
 see https://docs.ceph.com/en/latest/mgr/prometheus/

k exec -it deploy/ceph-mgr-deploy bash
ceph mgr module enable prometheus
ceph config set mgr mgr/prometheus/rbd_stats_pools "*"

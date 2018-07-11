# Agent v1.0.944 (2018-06-27)
FROM newrelic/infrastructure:0.0.24
ADD nr-kubernetes-definition.yml /var/db/newrelic-infra/newrelic-integrations/
ADD bin/nr-kubernetes /var/db/newrelic-infra/newrelic-integrations/bin/
# Warning: First, Edit sample file to suit your needs and rename it to
# `nr-kubernetes-config.yml`
ADD nr-kubernetes-config.yml.sample /etc/newrelic-infra/integrations.d/nr-kubernetes-config.yml

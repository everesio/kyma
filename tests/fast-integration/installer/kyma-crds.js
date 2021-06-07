const { exportCluster } = require("@kubernetes/client-node/dist/config_types");

const kymaCrds =
  [
    'authorizationpolicies.security.istio.io',
    'destinationrules.networking.istio.io',
    'envoyfilters.networking.istio.io',
    'gateways.networking.istio.io',
    'istiooperators.install.istio.io',
    'peerauthentications.security.istio.io',
    'requestauthentications.security.istio.io',
    'serviceentries.networking.istio.io',
    'sidecars.networking.istio.io',
    'virtualservices.networking.istio.io',
    'workloadentries.networking.istio.io',
    'workloadgroups.networking.istio.io',
    'addonsconfigurations.addons.kyma-project.io',
    'alertmanagerconfigs.monitoring.coreos.com',
    'alertmanagers.monitoring.coreos.com',
    'apirules.gateway.kyma-project.io',
    'applicationmappings.applicationconnector.kyma-project.io',
    'applications.applicationconnector.kyma-project.io',
    'assetgroups.rafter.kyma-project.io',
    'assets.rafter.kyma-project.io',
    'backendmodules.ui.kyma-project.io',
    'buckets.rafter.kyma-project.io',
    'centralconnections.applicationconnector.kyma-project.io',
    'certificaterequests.applicationconnector.kyma-project.io',
    'clusteraddonsconfigurations.addons.kyma-project.io',
    'clusterassetgroups.rafter.kyma-project.io',
    'clusterassets.rafter.kyma-project.io',
    'clusterbuckets.rafter.kyma-project.io',
    'clustermicrofrontends.ui.kyma-project.io',
    'clusterservicebrokers.servicecatalog.k8s.io',
    'clusterserviceclasses.servicecatalog.k8s.io',
    'clusterserviceplans.servicecatalog.k8s.io',
    'compassconnections.compass.kyma-project.io',
    'clustertestsuites.testing.kyma-project.io',
    'eventactivations.applicationconnector.kyma-project.io',
    'functions.serverless.kyma-project.io',
    'gitrepositories.serverless.kyma-project.io',
    'groups.authentication.kyma-project.io',
    'httpsources.sources.kyma-project.io',
    'jaegers.jaegertracing.io',
    'microfrontends.ui.kyma-project.io',
    'oauth2clients.hydra.ory.sh',
    'podmonitors.monitoring.coreos.com',
    'podpresets.settings.svcat.k8s.io',
    'probes.monitoring.coreos.com',
    'prometheuses.monitoring.coreos.com',
    'prometheusrules.monitoring.coreos.com',
    'rules.oathkeeper.ory.sh',
    'servicebindings.servicecatalog.k8s.io',
    'servicebindingusages.servicecatalog.kyma-project.io',
    'servicebrokers.servicecatalog.k8s.io',
    'serviceclasses.servicecatalog.k8s.io',
    'serviceinstances.servicecatalog.k8s.io',
    'servicemonitors.monitoring.coreos.com',
    'serviceplans.servicecatalog.k8s.io',
    'subscriptions.eventing.kyma-project.io',
    'testdefinitions.testing.kyma-project.io',
    'thanosrulers.monitoring.coreos.com',
    'tokenrequests.applicationconnector.kyma-project.io',
    'usagekinds.servicecatalog.kyma-project.io',
    'authcodes.dex.coreos.com',
    'authrequests.dex.coreos.com',
    'oauth2clients.dex.coreos.com',
    'signingkeies.dex.coreos.com',
    'refreshtokens.dex.coreos.com',
    'passwords.dex.coreos.com',
    'offlinesessionses.dex.coreos.com',
    'connectors.dex.coreos.com',
    'devicerequests.dex.coreos.com',
    'devicetokens.dex.coreos.com',
    'natsclusters.nats.io',
    'natsserviceroles.nats.io',
  ];
module.exports = kymaCrds;

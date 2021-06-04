module github.com/kyma-project/kyma/components/compass-runtime-agent

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/kofalt/go-memoize v0.0.0-20200917044458-9b55a8d73e1c
	github.com/kyma-incubator/compass v0.0.0-20200813093525-96b1a733a11b
	github.com/kyma-incubator/compass/components/director v0.0.0-20200813093525-96b1a733a11b
	github.com/kyma-project/kyma/components/application-operator v0.0.0-20210604090801-9a75cdfd1c7c
	github.com/kyma-project/rafter v0.0.0-20200626063334-5a8dd27d1976
	github.com/machinebox/graphql v0.2.3-0.20181106130121-3a9253180225
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/vrischmann/envconfig v1.3.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/metrics v0.20.1
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.4
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.25+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc92

	go.etcd.io/etcd => go.etcd.io/etcd v3.3.25+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/text => golang.org/x/text v0.3.3

)

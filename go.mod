module github.com/dnephin/dobi

go 1.13

require (
	github.com/Unknwon/com v0.0.0-20170213072014-0db4a625e949 // indirect
	github.com/dnephin/configtf v0.0.0-20161020003418-6b0d1fdf5e68
	github.com/dnephin/go-os-user v0.0.0-20161029070903-44e2994deb1e
	github.com/docker/cli v0.0.0-20200303215952-eb310fca4956
	github.com/docker/docker v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/fsouza/go-dockerclient v1.6.4
	github.com/gogits/git-module v0.0.0-20170608205522-1de103dca47a
	github.com/golang/mock v1.1.1
	github.com/google/go-cmp v0.4.0
	github.com/kballard/go-shellquote v0.0.0-20170619183022-cd60e84ee657
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mcuadros/go-version v0.0.0-20161105183618-257f7b9a7d87 // indirect
	github.com/metakeule/fmtdate v1.1.2-0.20150502212323-427373e7d2f8
	github.com/opencontainers/runc v1.0.0-rc3.0.20170716065720-825b5c020ace // indirect
	github.com/pkg/errors v0.8.1
	github.com/renstrom/dedent v1.0.1-0.20150819195903-020d11c3b9c0
	github.com/sirupsen/logrus v1.4.1
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cobra v0.0.2-0.20171109065643-2da4a54c5cee
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v0.0.0-20170224212429-dcecefd839c4
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.4
	gotest.tools/v3 v3.0.2
)

replace github.com/spf13/cobra => github.com/dnephin/cobra v1.5.2-0.20170125185912-5d13e8c9d917

replace github.com/Nvveen/Gotty => github.com/ijc/Gotty v0.0.0-20170406111628-a8b993ba6abd

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20190830141801-acfa387b8d69

%if ! 0%{?gobuild:1}
%define gobuild(o:) go build -ldflags "${LDFLAGS:-} -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \\n')" -a -v -x %{?**};
%endif
%define  debug_package %{nil}

Name:           onlineconf-yaml
Version:        %{__version}
Release:        %{__release}%{?dist}
Summary:        Utilities for сonvert yml to cdb and import yml to OnlineConf
Group:          OnlineConf
License:        Proprietary
URL:            https://github.com/askobenkov/onlineconf-yaml
Source:         %{name}-%{version}.tar.gz
Prefix:         %{_prefix}
BuildRoot:      %{_tmppath}/%{name}-root

BuildRequires:  golang
BuildRequires:  make
BuildRequires:  which

%description
Utilities for сonvert yml to cdb and import yml to OnlineConf

%prep
%setup -q
mkdir -p go/src/github.com/askobenkov/
ln -sn $(pwd) go/src/github.com/askobenkov/onlineconf-yaml

%build
export GOPATH=$(pwd)/go
cd go/src/github.com/askobenkov/onlineconf-yaml
make build

%install
export GOPATH=$(pwd)/go

[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf $RPM_BUILD_ROOT
install -D $GOPATH/bin/yml2cdb $RPM_BUILD_ROOT%{_bindir}/yml2cdb
install -D $GOPATH/bin/yml2onlineconf $RPM_BUILD_ROOT%{_bindir}/yml2onlineconf
mkdir -m 740 -p $RPM_BUILD_ROOT%{_sysconfdir}/%{name}

%pre

%post

%clean
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root)
%{_bindir}/yml2cdb
%{_bindir}/yml2onlineconf

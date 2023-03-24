Name:		gcc		
Version:	5.4.0
Release:	1%{?dist}
Summary:	install gcc 5.4.0

Group:		avalon
License:	GPL
URL:		http://avalon.com
Source0:	gcc-5.4.0.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
#BuildRequires:	gcc
#Requires:	

%define __debug_install_post   \
   %{_rpmconfigdir}/find-debuginfo.sh %{?_find_debuginfo_opts} "%{_builddir}/%{?buildsubdir}"\
%{nil}
%description
The gcc package contains the GNU Compiler Collection.

%prep
%setup -q


%build
./contrib/download_prerequisites
mkdir gcc-build-5.4.0
cd gcc-build-5.4.0
#$PWD/../configure --prefix=/usr/local/gcc48 --enable-checking=release --enable-languages=c,c++ --disable-multilib
$PWD/../configure --enable-checking=release --enable-languages=c,c++ --disable-multilib
make %{?_smp_mflags}


%install
%{__rm} -rf %{buildroot}
cd gcc-build-5.4.0
make install DESTDIR=%{buildroot}

%clean
%{__rm} -rf %{buildroot}

%pre

%post
echo "/usr/local/lib64" > /etc/ld.so.conf.d/usr_local_lib.conf
/sbin/ldconfig &> /dev/null


%files
%defattr(-, root, root, 0755)
/usr/local/bin
/usr/local/include
/usr/local/lib
/usr/local/lib64
/usr/local/libexec
/usr/local/share

#%doc

%changelog
* Tue Apr 23 2019 lei.guo <lei.guo@funplus.com> - 4.8.5-2
- change

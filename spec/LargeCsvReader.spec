Name:           LargeCsvReader
Version:        %{app_version}
Release:        1%{?dist}
Summary:        App for previewing large csv files

License:        MIT
URL:            https://github.com/RikudouSage/LargeCsvReader
Source0:        %{name}.tar.xz

BuildArch:      x86_64
BuildRequires:  make

%description
App for previewing large csv files

%prep
echo ======================
echo %{buildroot}
echo ======================
%setup -c -T
tar -xJf %{SOURCE0}

%install
rm -rf %{buildroot}
make install DESTDIR=%{buildroot}

%files
%{_bindir}/LargeCsvReader
%{_datadir}/applications/LargeCsvReader.desktop
%{_datadir}/pixmaps/LargeCsvReader.png

%changelog
* Tue Sep 24 2024 RikudouSage <me@rikudousage.com> - %{app_version}-1
- Initial package

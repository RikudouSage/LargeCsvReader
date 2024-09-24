AppId := cz.chrastecky.large_csv_reader

build_windows:
	fyne-cross windows -app-id ${AppId}

build_linux:
	fyne-cross linux -app-id ${AppId} -release

build_macos:
	fyne-cross darwin -app-id ${AppId} -release -macosx-sdk-path bundled

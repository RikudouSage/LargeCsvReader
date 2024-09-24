package main

import (
	"LargeCsvReader/widgets"
	"encoding/csv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io"
	"os"
)

func showPreviewWindow(filePath string, fyneApp fyne.App) {
	window := fyneApp.NewWindow("Preview")
	window.Resize(fyne.NewSize(640, 480))

	typeToCharMap := map[string]rune{
		"Comma":     ',',
		"Semicolon": ';',
		"Tab":       '\t',
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	var csvFile *csv.Reader
	rows := make([][]string, 5)

	previewText := widget.NewLabel("Preview")
	previewText.Hide()

	tableContainer := container.NewStack()
	lineSeparator1 := widget.NewSeparator()
	lineSeparator2 := widget.NewSeparator()

	lineSeparator1.Hide()
	lineSeparator2.Hide()

	createTable := func() *widget.Table {
		table := widget.NewTable(
			func() (int, int) {
				return len(rows), len(rows[0])
			},
			func() fyne.CanvasObject {
				return widgets.NewCellWidget("default (hopefully) large enough text", nil) // placeholder to specify width
			},
			func(id widget.TableCellID, object fyne.CanvasObject) {
				cell := object.(*widgets.CellWidget)
				cell.SetText(rows[id.Row][id.Col])
				cell.OnRightClick = func(event *fyne.PointEvent) {
					items := []*fyne.MenuItem{
						fyne.NewMenuItem("Copy", func() {
							window.Clipboard().SetContent(cell.Text)
						}),
					}
					menu := fyne.NewMenu("", items...)
					canvas := fyne.CurrentApp().Driver().CanvasForObject(cell)
					widget.ShowPopUpMenuAtPosition(menu, canvas, event.AbsolutePosition)
				}
			},
		)

		return table
	}

	var loadMoreButton *widget.Button
	loadMoreButton = widget.NewButton("Load more", func() {
		for range 5 {
			line, err := csvFile.Read()
			if err == io.EOF {
				loadMoreButton.Disable()
				break
			}
			if err != nil {
				panic(err)
			}
			rows = append(rows, line)
		}

		tableContainer.Objects = []fyne.CanvasObject{createTable()}
		tableContainer.Refresh()
	})
	loadMoreButton.Hide()

	separatorSelect := widget.NewSelect([]string{"Comma", "Semicolon", "Tab"}, func(selected string) {
		csvFile = csv.NewReader(file)
		csvFile.LazyQuotes = true
		csvFile.Comma = typeToCharMap[selected]

		_, err = file.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}

		rows = [][]string{}
		for range 5 {
			line, err := csvFile.Read()
			if err == io.EOF {
				loadMoreButton.Disable()
				break
			}
			if err != nil {
				panic(err)
			}
			rows = append(rows, line)
		}

		previewText.Show()
		lineSeparator2.Show()
		lineSeparator1.Show()
		loadMoreButton.Show()

		table := createTable()
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	})

	mainContainer := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Separator"),
			separatorSelect,
			lineSeparator1,
			container.NewHBox(
				previewText,
				loadMoreButton,
			),
			lineSeparator2,
		),
		nil, nil, nil,
		tableContainer,
	)

	window.SetContent(
		container.New(
			layout.NewCustomPaddedLayout(8, 32, 32, 32),
			mainContainer,
		),
	)
	window.SetOnClosed(func() {
		file.Close()
	})

	window.Show()
}

func main() {
	fyneApp := app.NewWithID("1af17320-d9ef-4d2b-aae1-8226b14a177a")
	window := fyneApp.NewWindow("Hello")
	window.Resize(fyne.NewSize(640, 480))
	window.SetMaster()

	openDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if file != nil {
			showPreviewWindow(file.URI().Path(), fyneApp)
		}
	}, window)
	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))

	window.SetContent(container.NewVBox(
		container.NewCenter(
			widget.NewLabel("Using the button below choose a CSV you want to open"),
		),
		container.New(
			layout.NewCustomPaddedLayout(8, 0, 32, 32),
			widget.NewButton("Open CSV", func() {
				openDialog.Show()
			}),
		),
	))

	window.Show()
	fyneApp.Run()
}

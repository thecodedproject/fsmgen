package internal

import (
	flag "flag"
	gopkg "github.com/thecodedproject/gopkg"
	tmpl "github.com/thecodedproject/gopkg/tmpl"
	path "path"
)

var outputPath = flag.String("outdir,o", ".", "output directory for generated files")

type pkgDef struct {
	OutputPath string
	Import gopkg.ImportAndAlias
}

func Generate() error {

	flag.Parse()

	d, err := createPkgDef()
	if err != nil {
		return err
	}

	var files []gopkg.FileContents
	files, err = tmpl.AppendFileContents(
		files,
		// Add file creators here...
		fileTemplate(d),
	)
	if err != nil {
		return err
	}

	return gopkg.LintAndGenerate(files)
}

func createPkgDef() (pkgDef, error) {

	importPath, err := gopkg.PackageImportPath(*outputPath)
	if err != nil {
		return pkgDef{}, err
	}

	pkgName := path.Base(importPath)

	return pkgDef{
		OutputPath: *outputPath,
		Import: gopkg.ImportAndAlias{
			Import: importPath,
			Alias: pkgName,
		},
	}, nil
}


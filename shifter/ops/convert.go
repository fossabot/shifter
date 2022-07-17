/*
copyright 2019 google llc
licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at
    http://www.apache.org/licenses/license-2.0
unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

package ops

import (
	"github.com/google/uuid"
	"log"
	"path/filepath"
	"shifter/generator"
	"shifter/input"
	"shifter/lib"
)

type LogObject struct {
}

type Converter struct {
	UUID       string // Unique ID of the Run
	InputType  string
	SourcePath string
	Generator  string
	OutputPath string
	Flags      map[string]string

	SourceFiles []*FileObject
	OutputFiles []*FileObject

	Logs []*LogObject
}

type DownloadFile struct {
	Link     string `json:"link"`
	Filename string `json:"filename"`
}

// Input Types
const YAML string = "YAML"
const TEMPLATE string = "template"

// Create New Converter
func NewConverter(inputType string, sourcePath string, generator string, outputPath string, flags map[string]string) *Converter {
	// Create New Instance of Converter
	converter := &Converter{}

	// Create UUID for Converter
	converter.UUID = uuid.New().String()

	// Set all the Variables for the Converter
	converter.InputType = inputType
	converter.SourcePath = sourcePath
	converter.Generator = generator
	converter.OutputPath = outputPath
	converter.Flags = flags

	// Process the Path and Create Array of File Objects
	files, err := ProcessPath(converter.SourcePath)
	if err != nil {
		log.Println(err)
	}

	// Set Converter Files
	converter.SourceFiles = files
	if len(converter.SourceFiles) > 0 {
		converter.LoadSourceFiles()
	}

	return converter
}

func (converter *Converter) WriteSourceFiles() {
	// Process Input Objects
	for _, file := range converter.SourceFiles {
		file.WriteFile()
	}
}

func (converter *Converter) LoadSourceFiles() {
	// Process Input Objects
	for _, file := range converter.SourceFiles {
		file.LoadFile()
	}
}

func (converter *Converter) ListSourceFiles() {
	// Process Input Objects
	for _, file := range converter.SourceFiles {
		file.Meta()
	}
}

func (converter *Converter) ListOutputFiles() {
	// Process Input Objects
	for _, file := range converter.OutputFiles {
		file.Meta()
	}
}

func (converter *Converter) ConvertFiles() {
	// Process Input Objects
	for _, file := range converter.SourceFiles {

		var r []lib.Converted
		switch converter.InputType {
		case "yaml":
			sourceFile := input.Yaml(file.Content, converter.Flags)
			r = generator.NewGenerator(converter.Generator, file.Filename, sourceFile)
		case "template":
			sourceFile, values := input.Template(file.Content, converter.Flags)
			r = generator.NewGenerator(converter.Generator, file.Filename, sourceFile, values)
		}

		//outputFileName := fmt.Sprint(idx)
		for k := range r {
			fileObj := &FileObject{
				StorageType: file.StorageType,
				//SourcePath:    (converter.OutputPath + "/" + r[k].Path + r[k].Name + filepath.Ext(file.SourcePath)),
				Path:          (converter.OutputPath + "/" + r[k].Path + r[k].Name),
				Filename:      file.Filename,
				Ext:           filepath.Ext(file.Path),
				Content:       r[k].Payload,
				ContentLength: file.ContentLength,
			}

			// Write Converted File to Storage
			log.Printf("Writing to file %v", fileObj.Filename)
			fileObj.WriteFile()

			// Add Converted File Object to Converter
			converter.OutputFiles = append(converter.OutputFiles, fileObj)
		}
	}
}

func (converter *Converter) BuildDownloadFiles() []*DownloadFile {
	var files []*DownloadFile

	// Process Output Objects
	for _, file := range converter.OutputFiles {
		dlFile := &DownloadFile{}
		dlFile.Link = "https://somefile.com"
		dlFile.Filename = file.Filename
		files = append(files, dlFile)
	}

	return files
}

package devcon

import (
	"strings"
)

type OEMPackage struct {
	Name     string
	Provider string
	Class    string
}

// DPAdd adds a third-party (OEM) driver package to the driver store on the
// local computer. The infFilePath parameter is the fully qualified path and
// name of the INF file for the driver package.
//
// Notes
// DPAdd copies the specified INF file to the `%windir%/Inf` directory and
// renames it OEM*.inf. This file name is unique on the computer, and you cannot
// specify it.
//
// If this INF file already exists in `%windir%/Inf` (as determined by comparing
// the binary files, not by matching the file names) and the catalog (.cat) file
// for the INF is identical to a catalog file in the directory, the INF file is
// not recopied to the `%windir%/Inf` directory.
//
// This command calls the `SetupCopyOEMInf` function with no `CopyStyle` flags.
// `SetupCopyOEMInf` is described in the Microsoft Windows SDK documentation.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-dp-add for more information.
func (dc *DevCon) DPAdd(infFilePath string) error {
	lines, err := dc.run(commandDPAdd, infFilePath)
	if err != nil {
		return err
	}

	// TODO: Parse.
	dc.printResults(lines)

	return nil
}

// DPEnum returns the third-party (OEM) driver packages in the driver store on
// the local computer.
//
// Notes
// DPEnum returns the `OEM*.inf` files in the `%windir%/Inf` on the local computer.
// For each file, this command displays the provider, class, date, and version
// number from the INF file.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-dp-enum for more information.
func (dc *DevCon) DPEnum() ([]OEMPackage, error) {
	lines, err := dc.run(commandDPEnum)
	if err != nil {
		return nil, err
	}

	return parseOEMPackages(lines), nil
}

// DPDelete deletes a third-party (OEM) driver package from the driver store on
// the local computer. This command deletes the INF file, the PNF file, and the
// associated catalog file (.cat).
//
// THe infFileName represents the OEM*.inf file name of the INF file. Windows
// assigns a file name with this format to the INF file when you add the driver
// package to the driver store, such as by using DPAdd().
//
// Specifying true for force deletes the driver package even if a device is
// using it at the time.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-dp-delete for more information.
func (dc *DevCon) DPDelete(infFileName string, force bool) error {
	args := make([]string, 0)
	if force {
		args = append(args, "/f")
	}
	args = append(args, infFileName)

	lines, err := dc.run(commandDPDelete, args...)
	if err != nil {
		return err
	}

	// TODO: Parse.
	dc.printResults(lines)

	return nil
}

func parseOEMPackages(lines []string) []OEMPackage {
	lines = lines[1:]

	oemPackages := make([]OEMPackage, 0)

	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		oemPackage := OEMPackage{}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			if lineIndex == groupStart {
				oemPackage.Name = line
			} else {
				valuePair := parseColonSeparatedLine(line)

				if valuePair[0] == "Provider" {
					oemPackage.Provider = valuePair[1]
				} else if valuePair[0] == "Class" {
					oemPackage.Class = valuePair[1]
				}
			}
		}

		oemPackages = append(oemPackages, oemPackage)
	}

	return oemPackages
}

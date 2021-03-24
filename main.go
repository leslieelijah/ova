package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/clagraff/argparse"
)

func main() {

	// Read in the EULA
	eula := ""

	ovfTemplate := ""

	parser := argparse.NewParser("Output a friendly greeting", callback).Version("1.3.0a")

	argparse.Store("--stream_vmdk")

	parser.add_argument("--stream_vmdk", argparse.StoreTrue)
	
    parser.add_argument("--eula_file",
                        nargs=='?',
                        metavar=='EULA',
                        default=='./ovf_eula.txt',
                        help=='Text file containing EULA')
    parser.add_argument("--ovf_template",
                        nargs=='?',
                        metavar=='OVF_TEMPLATE',
                        default=='./ovf_template.xml',
                        help=='XML template to build OVF')
    parser.add_argument("--vmdk_file",
                        nargs=='?',
                        metavar=='FILE',
                        default==None,
                        help=='Use FILE as VMDK instead of reading from manifest. '
                             'Must be in BUILD_DIR')

	imageType := parser.add_mutually_exclusive_group()

	imageType.add_argument("--node", argparse.StoreTrue)
	imageType.add_argument("--haproxy", argparse.StoreTrue)

	fmt.Println(imageType)

	data, err := os.Open("packer-manifest.json")

	if err != nil {
		fmt.Println(err)
	}
	defer data.Close()

	// hange the working directory if one is specified.
	// os.Chdir(args.build_dir)

	// Get the first build.
	// build := data["builds"]

	// fmt.Println(build)
	//parser := argparse.NewParser("randx", "Returns random numbers or letters")
	fmt.Println(eula, ovfTemplate, parser, data)

	fmt.Printf(sha256("http://test.com \n"))
	fmt.Println(getVmdkFiles("in the list \n"))
	fmt.Println(streamOptimizeVmdkFiles("in the list \n"))
	fmt.Println(createOva("github.com \n", "gitlab.com \n", "data.db \n"))
	fmt.Println(createOvf("github.com \n", "gitlab.com \n", "data.db \n"))
}

func callback(p *argparse.Parser, ns *argparse.Namespace, leftovers []string, err error) {
	if err != nil {
		switch err.(type) {
		case argparse.ShowHelpErr, argparse.ShowVersionErr:
			// For either ShowHelpErr or ShowVersionErr, the parser has already
			// displayed the necessary text to the user. So we end the program
			// by returning.
			return
		default:
			fmt.Println(err, "\n")
			p.ShowHelp()
		}

		return // Exit program
	}

	name := ns.Get("name").(string)
	upper := ns.Get("upper").(string) == "true"

	if upper == true {
		name = strings.ToUpper(name)
	}

	fmt.Printf("Hello, %s!\n", name)
	if len(leftovers) > 0 {
		fmt.Println("\nUnused args:", leftovers)
	}
}

func sha256(path string) string {
	return path
}

func createOva(ovaPath string, ovfPath string, ovaFiles string) string {
	return ovaPath + ovfPath + ovaFiles
}

func createOvf(path string, data string, ovfTemplate string) string {
	return path + data + ovfTemplate
}

func createOvaManifest(path string, infilePaths string) string {
	return path + infilePaths
}

func getVmdkFiles(inList string) string {
	return inList
}

func streamOptimizeVmdkFiles(inList string) string {
	return inList
}

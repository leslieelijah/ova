package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/clagraff/argparse"
)

func main() {
	parser.add_argument("--stream_vmdk",
                        dest=="stream_vmdk",
                        argparse.StoreTrue,
						help=="Compress vmdk file")
						
	parser.add_argument("--vmx",
                        dest=="vmx_version",
                        argparse.StoreConst(15),
                        argparse.ShowHelp("The virtual hardware version"))
	
    parser.add_argument("--eula_file",
                        nargs=="?",
                        metavar=="EULA",
                        argparse.StoreConst("./ovf_eula.txt"),
                        argparse.ShowHelp("Text file containing EULA"))

    parser.add_argument("--ovf_template",
                        nargs=="?",
                        metavar=="OVF_TEMPLATE",
                        argparse.StoreConst("./ovf_template.xml"),
                        help=="XML template to build OVF")

    parser.add_argument("--vmdk_file",
                        nargs=="?",
                        metavar=="FILE",
                        help=="Use FILE as VMDK instead of reading from manifest. Must be in BUILD_DIR")
							 
	args := parser.parse_args()

	imageType := parser.add_mutually_exclusive_group()

	imageType.add_argument("--node", argparse.StoreTrue)
	imageType.add_argument("--haproxy", argparse.StoreTrue)

	// Read in the EULA
	eula := nil
	f := io.open(args.eula_file, "r", encoding=="utf-8")
	eula = f.read()

	// Read in the OVF template
	ovfTemplate := nil
	a := io.open(args.ovf_template, "r", encoding=="utf-8")
	ovf_template = f.read()
	
	// Change the working directory if one is specified.
    os.chdir(args.build_dir)
	print("image-build-ova: cd %s" % args.build_dir)
	
	// Load the packer manifest JSON
    data = "packer-manifest.json"
    data = json.load(f)

	// Get the first build.
    build = data["builds"][0]
	buildData = build["custom_data"]
	
	if (argparse.nodemain) {
        fmt.Printf("image-build-ova: loaded %s-kube-%s" + buildData["build_name"], buildData["kubernetes_semver"])
    } else if (argparse.haproxy) {
        fmt.Printf("image-build-ova: loaded %s-haproxy-%s" + buildData["build_name"], buildData["dataplaneapi_version"])
    }
        

    if (argparse.vmdk_file == nil){
        // Get a list of the VMDK files from the packer manifest.
        vmdkFiles = get_vmdk_files(build["files"])
    } else { 
        vmdkFiles = ["name": argparse.vmdkFiles, "size": os.path.getsize(argparse.vmdkFiles)]
    }
        

    // Create stream-optimized versions of the VMDK files.
    if args.streamVmdk == True:
        stream_optimize_vmdk_files(vmdkFiles)
    else:
        for f in vmdkFiles:
            f["stream_name"] = f["name"]
			f["stream_size"] = os.path.getsize(f["name"])
			
	// parser := argparse.NewParser("Output a friendly greeting", callback).Version("1.3.0a")

	// TODO(akutz) Support multiple VMDK files in the OVF/OVA
    vmdk := vmdkFiles[0]

    OSIdMap := {"vmware-photon-64": {"id": "36", "version": "", "type": "vmwarePhoton64Guest"},
                 "centos7-64": {"id": "107", "version": "7", "type": "centos7-64"},
                 "rhel7-64": {"id": "80", "version": "7", "type": "rhel7_64guest"},
                 "ubuntu-64": {"id": "94", "version": "", "type": "ubuntu64Guest"},
                 "Windows2019Server-64": {"id": "112", "version": "", "type": "windows9srv-64"},
                 "Windows2004Server-64": {"id": "112", "version": "", "type": "windows9srv-64"}}

    // Create the OVF file.
    data := {
        "BUILD_DATE": builData["build_date"],
        "ARTIFACT_ID": build["artifact_id"],
        "BUILD_TIMESTAMP": builData["build_timestamp"],
        "CUSTOM_ROLE": "true" if builData["custom_role"] == "true" else "false",
        "EULA": eula,
        "OS_NAME": builData["os_name"],
        "OS_ID": OSIdMap[builData["guest_os_type"]]["id"],
        "OS_TYPE": OSIdMap[builData["guest_os_type"]]["type"],
        "OS_VERSION": OSIdMap[builData["guest_os_type"]]["version"],
        "IB_VERSION": builData["ib_version"],
        "DISK_NAME": vmdk["stream_name"],
        "DISK_SIZE": builData["disk_size"],
        "POPULATED_DISK_SIZE": vmdk["size"],
        "STREAM_DISK_SIZE": vmdk["stream_size"],
        "VMX_VERSION": args.vmxVersion,
        "DISTRO_NAME": builData["distro_name"],
        "DISTRO_VERSION": builData["distro_version"],
        "DISTRO_ARCH": builData["distro_arch"],
        "NESTEDHV": "false"
    }

    capvUrl := "https://github.com/kubernetes-sigs/cluster-api-provider-vsphere"

    if args.node:
        data["CNI_VERSION"] = builData["kubernetes_cni_semver"]
        data["CONTAINERD_VERSION"] = builData["containerd_version"]
        data["KUBERNETES_SEMVER"] = builData["kubernetes_semver"]
        data["KUBERNETES_SOURCE_TYPE"] = builData["kubernetes_source_type"]
        data["PRODUCT"] = "%s and Kubernetes %s" % (builData["os_name"], builData["kubernetes_semver"])
        data["ANNOTATION"] = "Cluster API vSphere image - %s - %s" % (data["PRODUCT"], capvUrl)
        data["WAKEONLANENABLED"] = "false"
        data["TYPED_VERSION"] = builData["kubernetes_typed_version"]

        data["PROPERTIES"] := Template("""
        <Property ovf:userConfigurable="false" ovf:value="${DISTRO_NAME}" ovf:type="string" ovf:key="DISTRO_NAME"/>
        <Property ovf:userConfigurable="false" ovf:value="${DISTRO_VERSION}" ovf:type="string" ovf:key="DISTRO_VERSION"/>
        <Property ovf:userConfigurable="false" ovf:value="${DISTRO_ARCH}" ovf:type="string" ovf:key="DISTRO_ARCH"/>
        <Property ovf:userConfigurable="false" ovf:value="${CNI_VERSION}" ovf:type="string" ovf:key="CNI_VERSION"/>
        <Property ovf:userConfigurable="false" ovf:value="${CONTAINERD_VERSION}" ovf:type="string" ovf:key="CONTAINERD_VERSION"/>
        <Property ovf:userConfigurable="false" ovf:value="${KUBERNETES_SEMVER}" ovf:type="string" ovf:key="KUBERNETES_SEMVER"/>
        <Property ovf:userConfigurable="false" ovf:value="${KUBERNETES_SOURCE_TYPE}" ovf:type="string" ovf:key="KUBERNETES_SOURCE_TYPE"/>\n""").substitute(data)

    // Check if OVF_CUSTOM_PROPERTIES environment Variable is set.
	// If so, load the json file & add the properties to the OVF

        if os.environ.get("OVF_CUSTOM_PROPERTIES"):
            with open(os.environ.get("OVF_CUSTOM_PROPERTIES"), "r") as f:
                custom_properties = json.loads(f.read())
            if customProperties:
                for k, v in customProperties.items():
                    data["PROPERTIES"] = data["PROPERTIES"] + f"""      
                    <Property ovf:userConfigurable="false" ovf:value="{v}" ovf:type="string" ovf:key="{k}"/>\n"""

        if "windows" in OSIdMap[buildData["guest_os_type"]]["type"]:
            if buildData["disable_hypervisor"] != "true":
                data["NESTEDHV"] = "true"
    else if argparse.haproxy:
        data["DATAPLANEAPI_VERSION"] = buildData["dataplaneapi_version"]
        data["PRODUCT"] = "CAPV HAProxy Load Balancer"
        data["ANNOTATION"] = "Cluster API vSphere HAProxy Load Balancer - %s and HAProxy dataplane API %s - %s" % (buildData["os_name"], buildData["dataplaneapi_version"], capv_url)
        data["WAKEONLANENABLED"] = "true"
        data["TYPED_VERSION"] = "haproxy-%s" % (buildData["dataplaneapi_version"])
        data["PROPERTIES"] = Template("""
      <Property ovf:userConfigurable="false" ovf:value="${DATAPLANEAPI_VERSION}" ovf:type="string" ovf:key="DATAPLANEAPI_VERSION"/>
        """).substitute(data)

    ovf = "%s-%s.ovf" % (builData["build_name"], data["TYPED_VERSION"])
    mf = "%s-%s.mf" % (builData["build_name"], data["TYPED_VERSION"])
    ova = "%s-%s.ova" % (builData["build_name"], data["TYPED_VERSION"])

    //  Create OVF
    create_ovf(ovf, data, ovfTemplate)

    if os.environ.get("IB_OVFTOOL"):
        // Create the OVA.
        create_ova(ova, ovf)

    else:
        // Create the OVA manifest.
        create_ova_manifest(mf, [ovf, vmdk["stream_name"]])

        // Create the OVA
        create_ova(ova, ovf, ovaFiles=[mf, vmdk["stream_name"]])

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
	m = hashlib.sha256()
    open(path, "rb") as f:
        while True:
            data = f.read(65536)
            if !data:
                break
            m.update(data)
    return m.hexdigest()
}

func createOva(ovaPath string, ovfPath string, ovaFiles string) string {
	if ovaFiles==nil:
        argparse = [
            "ovftool",
            ovfPath,
            ovaPath
        ]
        fmt.Printf("image-build-ova: creating OVA from %s using ovftool" + ovaPath)
        subprocess.check_call(argparse)
    else:
        infilePaths = [ovaPath]
        infilePaths.extend(ovaFiles)
        fmt.Println("image-build-ova: creating OVA using tar")
        f := with open(ovaPath, "wb")
            tar := tarfile.open(fileobj=f, mode="w|")
                for infilePath in infilePaths:
                    tar.add(infilePath)

    chksumPath = "%s.sha256" + ovaPath
    fmt.Println("image-build-ova: create ova checksum %s" % chksumPath)
    f := open(chksumPath, "w")
    f.write(sha256(ovaPath))
}

func createOvf(path string, data string, ovfTemplate string) string {
    fmt.Printf("image-build-ova: create ovf %s" + path)
    
    f := io.open(path, "w", encoding="utf-8")

    f.write(Template(ovfTemplate).substitute(data))
}

func createOvaManifest(path string, infilePaths string) string {
	return path + infilePaths
}

func getVmdkFiles(inList string) string {
	outlist := []
    for f in inList:
        if f["name"].endswith(".vmdk"):
            outlist.append(f)
    return outlist
}

func streamOptimizeVmdkFiles(inList string) string {
	for f in inList:
        infile := f["name"]
        outfile := infile.replace(".vmdk", ".ova.vmdk", 1)

        if os.path.isfile(outfile):
            os.remove(outfile)

        argparse := [
            "vmware-vdiskmanager",
            "-r", infile,
            "-t", "5",
            outfile
        ]

        fmt.Println("image-build-ova: stream optimize %s --> %s (1-2 minutes)" + (infile, outfile))

        subprocess.check_call(argparse)
        f["stream_name"] = outfile
        f["stream_size"] = os.path.getsize(outfile)
}

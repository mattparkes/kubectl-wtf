package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"

	//	"crypto/x509"

	//	"k8s.io/cli-runtime"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"k8s.io/client-go/kubernetes"
	//	"k8s.io/client-go/tools/clientcmd"

	hostfile "github.com/lextoumbourou/goodhosts"
	//"k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes"
)

type Result struct {
	//level (0-5) for info, warning, error, etc
	level      int
	problem    string
	suggestion string
}

func Print(message string, severity string, indentation int) {
	indents := ""
	for i := 0; i < indentation; i++ {
		indents += " "
	}

	colouredSeverity := ""
	switch severity {
	case "error":
		colouredSeverity = Sprintf(Bold(Red("[Error]")))
	case "warning":
		colouredSeverity = Sprintf(Brown("[Warn]"))
	case "info":
		colouredSeverity = Sprintf(Cyan("[Info]"))
	case "ok":
		colouredSeverity = Sprintf(Green("[OK]"))
	case "todo":
		colouredSeverity = Sprintf(Blue("[TODO]"))
	default:
		colouredSeverity = ""
	}

	switch severity {
	case "error":
		fmt.Printf("%s%s: %s\n", indents, colouredSeverity, Red(message))
	case "data":
		fmt.Printf("%s%s\n", indents, Magenta(message))
	case "action":
		fmt.Printf("%s%s\n", indents, Bold(message))
	default:
		fmt.Printf("%s%s: %s\n", indents, colouredSeverity, message)
	}

	return
}

//var color aurora.Aurora
func main() {

	//	textPtr := flag.String("text", "", "Text to parse.")
	//	metricPtr := flag.String("metric", "chars", "Metric {chars|words|lines};.")
	//	uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")
	//	fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *textPtr, *metricPtr, *uniquePtr)

	//--version

	//--output / -o quiet|verbose|html|json

	//--no-hooks

	//--hook <path to hook>

	flag.Parse()

	//var colorsFlag = flag.Bool("colors", false, "enable or disable colors")
	//color = aurora.NewAurora(*colorsFlag)

	resourceType := ""
	resourceName := ""

	// Support <type>/<name> syntax (e.g output of kubectl get ingress -o name)
	// Support <type> <name> syntax (typical kubectl usage e.g kubectl get ingress foo)
	if flag.NArg() == 1 && strings.Contains(flag.Arg(0), "/") {
		resourceType = strings.Split(flag.Arg(0), "/")[0]
		resourceName = strings.Split(flag.Arg(0), "/")[1]
	} else if flag.NArg() == 2 {
		resourceType = flag.Arg(0)
		resourceName = flag.Arg(1)
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch resourceType {
	case "ingress":
		checkIngress(resourceName)
	default:
		fmt.Println(fmt.Sprintf("Non-supported resource type '%s'", resourceType))
		//flag.PrintDefaults()
		os.Exit(1)
	}

}

func checkIngress(ingressName string) {
	namespace := "default"

	hostnames, ports, certificates, paths, services, pods := ingressCheckResources(ingressName, namespace)
	certificates = certificates
	paths = paths
	services = services
	pods = pods

	//check ingress class
	//check ingress health

	ingressCheckHostsFile(hostnames)
	ingressCheckDNS(hostnames)
	ingressCheckTCP(hostnames, ports)
	ingressCheckCertificate(hostnames)

	fmt.Println(Sprintf(Blue("\nTODO: INFO if hostnames point to a (known) CDN")))

	fmt.Println(Sprintf(Blue("\nTODO: for each service, checkService(service)")))
}

func checkService() {
	fmt.Println("\nTODO: error if upstreams not set or don't exist")

	// port-forward and netcat

	fmt.Println("\nTODO: for each Pod, checkPod(pod)")
}

func checkPod() {
	// check state
	//  - warn if imagepullbackoff, failing readiness, can't mount/attach volumes

	// grep logs for errors

	//port-forward and netcat
}

func ingressCheckResources(ingressName string, namespace string) (hostnames [2]string, ports [2]int, certificates [1]string, paths [1]string, services [1]string, pods [0]string) {
	Print(fmt.Sprintf("\nChecking Ingress '%s' in namespace '%s' to ensure it exists:", ingressName, namespace), "action", 0)

	//TODO Lookup ingress

	if true {
		Print(fmt.Sprintf("Ingress %s exists in namespace %s:", ingressName, namespace), "ok", 2)

		hostnames = [...]string{"mattparkes.net", "notreal.host"}
		ports = [...]int{80, 443}
		certificates = [...]string{"mattparkes.net"}
		paths = [...]string{"/"}
		services = [...]string{"test"}
		pods = [...]string{}

		//TODO: print ingress details: hosts, paths, upstream services, upstream pods
		//return details (for use by other checks)
		Print(fmt.Sprintf("Hostnames: %s", hostnames), "data", 4)
		Print(fmt.Sprintf("Certificates: %s", certificates), "data", 4)
		Print(fmt.Sprintf("Paths: %s", paths), "data", 4)
		Print(fmt.Sprintf("Backend Services: %s", services), "data", 4)
		Print(fmt.Sprintf("Backend Pods: %s", pods), "data", 4)
	} else {
		Print(fmt.Sprintf("Ingress '%s' does not exist in namespace '%s'\n", Bold(ingressName), Bold(namespace)), "error", 2)
		os.Exit(1)
	}

	return hostnames, ports, certificates, paths, services, pods
}

// Check local OS hosts file and warn if entry found
func ingressCheckHostsFile(hostnames [2]string) {
	Print(fmt.Sprintf("\nChecking local hosts file for %s:", hostnames), "action", 0)

	for _, hostname := range hostnames {
		found := false
		hosts, err := hostfile.NewHosts()
		if err == nil {
			for i, line := range hosts.Lines {
				if line.Hosts != nil {
					for _, host := range line.Hosts {
						if host == hostname {
							found = true
							Print(fmt.Sprintf("Local hosts file entry for '%s' found:", hostname), "warning", 2)
							Print(fmt.Sprintf("line %s:    %s", strconv.Itoa(i+1), line.Raw), "data", 4)
						}
					}
				}
			}
			if !found {
				Print(fmt.Sprintf("No local hosts file entry found for '%s'", hostname), "ok", 2)
			}
		} else {
			Print(fmt.Sprintf("Unable to check local hosts file: %s", err), "warning", 2)
		}
	}
	return
}

func ingressCheckDNS(hostnames [2]string) {
	Print(fmt.Sprintf("\nChecking DNS resolution for %s:", hostnames), "action", 0)

	for _, hostname := range hostnames {
		hosts, err := net.LookupHost(hostname)
		if err == nil {
			Print(fmt.Sprintf("Hostname '%s' resolves to: %s", hostname, hosts), "ok", 2)
		} else {
			Print(fmt.Sprintf("Unable to resolve Hostname '%s'", hostname), "error", 2)
		}
	}
}

func ingressCheckTCP(hostnames [2]string, ports [2]int) {
	Print(fmt.Sprintf("\nChecking TCP Ports %v for %s:", ports, hostnames), "action", 2)

	for _, hostname := range hostnames {
		for _, port := range ports {
			address := fmt.Sprintf("%s:%d", hostname, port)

			_, err := net.DialTimeout("tcp", address, 10*time.Second)
			if err != nil {
				Print(fmt.Sprintf("TCP connection to %s:%d could not be established", hostname, port), "error", 2)
				Print(fmt.Sprintf("%s", err), "data", 4)
			} else {
				Print(fmt.Sprintf("TCP connection to %s:%d successfully established", hostname, port), "ok", 2)
			}
		}
	}
}

func ingressCheckCertificate(hostnames [2]string) {
	Print(fmt.Sprintf("\nTODO: Checking SSL Certificates for %s:", hostnames), "action", 2)
	for _, hostname := range hostnames {
		Print(hostname, "todo", 4)
		//if _, err := cert.Verify(opts); err != nil {
		//	panic("failed to verify certificate: " + err.Error())
		//}

		// get certificate from Kubernetes for each hostname
		// check it's expiry, self-signed, hostnames, etc

		//get certificate from port 443 and check it matches (if not probably a CDN in the way, etc)
	}
}

//func service CheckService(service) {
//}

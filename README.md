### All Crash-course steps

Step1: golang cli application using cobra-cli.  
Step2: zerolog for log levels - info, debug, trace, warn, error.  
Step3: pflag with flags for logs level.  
Step3+: Use Viper to add env vars.  
Step4: fasthttp with cobra command "server" and flags for server port.  
Step4+: Add http requests logging.  
Step5: makefile, distroless dockerfile, github worflow and initial tests.  
Step5+:Trivy vulnerabilities check.  
Step6: k8s.io/client-go to create function to list Kubernetes deployment resources in default namespace. Auth via kubeconfig. add flags for set kubeconfig file. list cli command call function.  
Step6+: Add create/delete command.
Step7: using k8s.io/client-go create list/watch informer for Kubernetes deployment resources. Auth via kubeconfig or in-cluster auth add flags for in-cluster mode. Informer report events in logs. Add envtest unit tests.  
Step7+: add custom logic function for update/delete events using informers cache search.  
Step7++: use config to setup informers start configuration.  
Step8: json api handler to request list deployment resources in informer cache storage.   
Step9: sigs.k8s.io/controller-runtime and controller with logic to report in log each event received with informer.  
Step9+: multi-cluster informers. Dynamically created informers.  
Step10: controller mgr to control informer and controller. Leader election with lease resource. flag to disable leader election. flag for mgr metrics port.  
Step11: custom crd Frontendpage with additional informer, controller with additional reconciliation logic for custom resource.  
Step11++: add multi-project client configuration for management clusters.  
Step12: platform engineering integration based on https://docs.port.io/actions-and-automations/create-self-service-experiences/setup-the-backend. API handler for actions to CRUD custom resource.  
Step12+: Add Update action support for IDP and controller.  
Step12++: Discord notifications integration.  
Step13: use github.com/mark3labs/mcp-go/mcp to create mcp server for api handlers as a mcp tools flag to specify mcp port.  
Step13+: Add delete/update MCP tool.  
Step13++: Add oidc auth to MCP.  
Step14: jwt authentication and authorisation for api and mcp.  
Step14+: add jwt auth for MCP.  
Step15: basic Opentelemetry code instrumentation.  
Step15++: 90% test coverage.  

### Example repos:

https://github.com/dereban25/go-kubernetes-controllers  
https://github.com/Sp1tfire88/k8s-controllers  
https://github.com/oleksandr-san/k8s-controller  
https://github.com/vanelin/k8s-controller  
https://github.com/dmytropakki1995/k8s-controllers  
https://github.com/roman-povoroznyk/kubernetes-controller  
https://github.com/egormak/k8s-controller  
https://github.com/alioss/k8s-controller-crash-course  
https://github.com/Michaelcode2/k8s-controller-sample  
https://github.com/roman-povoroznyk/kubernetes-controller  
https://github.com/fataevalex/k8s-controller  
https://github.com/AdamDubnytskyy/k8s-controller  
https://github.com/creativie/k8s-controller  
https://github.com/dolv/k8s-controller-tutorial  
https://github.com/silhouetteUA/k8s-controller  
https://github.com/ibra86/k8s-controller-patterns  
https://github.com/solaris24251/k8s-controller  
https://github.com/e1jefe/k8s-controller  
https://github.com/dliaudan/k8s-controllers-template  
https://github.com/dolv/k8s-controller-tutorial  
https://github.com/MikeBorovik/k8s-controller-tutorial  
https://github.com/thegostev/go-kubernetes-controllers  
https://github.com/HiulnaraPyvovar/k8s-controller-tutorial  
https://github.com/danteprivet/k8s-controller  
https://github.com/dmytropakki1995/k8s-controllers  
https://github.com/JRaver/k8s-controller-tutorial   

### How to run a controller's Docker container in a Windows environment using git-bash:

```MSYS_NO_PATHCONV=1 docker run -e HOME=/tmp -v C:/Users/bergs/.kube:/tmp/.kube bergshrund/k8s-controller:1ccb937-dirty-linux-amd64 server --kubeconfig=/tmp/.kube/config```

MSYS_NO_PATHCONV=1 - disable path substitution in git-bash for Windows  
-e HOME=/tmp - set variable HOME to new value  
-v C:/Users/localuser/.kube:/tmp/.kube  
package main

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const defaultPort = "8080"

//go:embed templates/*.html
var templateFS embed.FS

type podInfo struct {
	Name      string
	Namespace string
	Phase     string
	PodIP     string
	NodeName  string
	Age       string
}

type podsPageData struct {
	Namespace string
	Pods      []podInfo
	Error     string
	Source    string
	QueriedAt string
}

var (
	homeTemplate = template.Must(template.ParseFS(templateFS, "templates/home.html"))
	podsTemplate = template.Must(template.ParseFS(templateFS, "templates/pods.html"))
)

type server struct {
	client     kubernetes.Interface
	configName string
}

func main() {
	clientset, configName, err := newClientset()
	if err != nil {
		log.Fatalf("create kubernetes client: %v", err)
	}

	srv := server{client: clientset, configName: configName}
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleHome)
	mux.HandleFunc("/pods", srv.handlePods)

	port := envOrDefault("PORT", defaultPort)
	addr := ":" + port
	log.Printf("starting server on %s using %s", addr, configName)
	if err := http.ListenAndServe(addr, requestLogger(mux)); err != nil {
		log.Fatalf("http server stopped: %v", err)
	}
}

func newClientset() (*kubernetes.Clientset, string, error) {
	config, configName, err := loadKubeConfig()
	if err != nil {
		return nil, "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, "", err
	}

	return clientset, configName, nil
}

func loadKubeConfig() (*rest.Config, string, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, "in-cluster service account", nil
	}

  if kubeconfig := strings.TrimSpace(os.Getenv("KUBECONFIG")); kubeconfig != "" {
    config, kubeErr := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if kubeErr == nil {
      return config, "local kubeconfig from KUBECONFIG", nil
    }
    return nil, "", errors.Join(err, kubeErr)
  }

	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return nil, "", errors.Join(err, homeErr)
	}

	kubeconfig := filepath.Join(homeDir, ".kube", "config")
	config, kubeErr := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if kubeErr == nil {
		return config, "local kubeconfig", nil
	}

	return nil, "", errors.Join(err, kubeErr)
}

func (s server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if execErr := homeTemplate.Execute(w, nil); execErr != nil {
		log.Printf("render home template: %v", execErr)
	}
}

func (s server) handlePods(w http.ResponseWriter, r *http.Request) {
	namespace := resolveNamespace(r)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pods, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	data := podsPageData{
		Namespace: namespace,
		Source:    s.configName,
		QueriedAt: time.Now().Format(time.RFC3339),
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.Error = err.Error()
	} else {
		data.Pods = make([]podInfo, 0, len(pods.Items))
		for _, pod := range pods.Items {
			data.Pods = append(data.Pods, podInfo{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Phase:     string(pod.Status.Phase),
				PodIP:     emptyValue(pod.Status.PodIP),
				NodeName:  emptyValue(pod.Spec.NodeName),
				Age:       time.Since(pod.CreationTimestamp.Time).Round(time.Second).String(),
			})
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if execErr := podsTemplate.Execute(w, data); execErr != nil {
		log.Printf("render pods template: %v", execErr)
	}
}

func resolveNamespace(r *http.Request) string {
	if namespace := strings.TrimSpace(r.URL.Query().Get("namespace")); namespace != "" {
		return namespace
	}

	if namespace := strings.TrimSpace(os.Getenv("POD_NAMESPACE")); namespace != "" {
		return namespace
	}

	if namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		trimmed := strings.TrimSpace(string(namespace))
		if trimmed != "" {
			return trimmed
		}
	}

	return metav1.NamespaceDefault
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func emptyValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started).Round(time.Millisecond))
	})
}

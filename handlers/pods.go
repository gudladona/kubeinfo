package handlers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
	"github.com/gudladona87/kubeinfo/models"
	"encoding/json"
)

//PodInfoHandler provides functions to query Pod metadata.
type PodInfoHandler struct {
	CoreClient corev1.CoreV1Interface
}

func (podInfo *PodInfoHandler) ListPods(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pods, err := podInfo.CoreClient.Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Printf("error listing pods: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("There are %d pods in the cluster\n", len(pods.Items))

	resp := models.Response{Message: "OK", PodCount: len(pods.Items)}
	b, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/julienschmidt/httprouter"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/gudladona87/kubeinfo/models"

	k8stesting "k8s.io/client-go/testing"
	metadata "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	core_v1 "k8s.io/api/core/v1"
)

func TestPodInfoHandler_ListPods(t *testing.T) {
	var successReaction k8stesting.ReactionFunc = func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &core_v1.PodList{Items: []core_v1.Pod{{
			ObjectMeta: metadata.ObjectMeta{Name: "testpod"},
		}}}, nil
	}

	var errorReaction k8stesting.ReactionFunc = func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("cannot reach kube API Server")
	}

	type fields struct {
		Clientset corev1.CoreV1Interface
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		p httprouter.Params
	}

	type want struct {
		status   int
		response models.Response
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "TestListpods Success",
			fields: fields{
				Clientset: func() corev1.CoreV1Interface {
					fakeClientSet := &fake.Clientset{
						Fake: k8stesting.Fake{
							ReactionChain: []k8stesting.Reactor{
								&k8stesting.SimpleReactor{Verb: "list", Resource: "pods", Reaction: successReaction},
							},
						},
					}
					return fakeClientSet.CoreV1()
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo", nil),
				p: httprouter.Params{},
			},
			want: want{
				status: http.StatusOK,
				response: models.Response{
					PodCount: 1,
					Message:  "OK",
				},
			},
		},
		{
			name: "TestListpods Error",
			fields: fields{
				Clientset: func() corev1.CoreV1Interface {
					fakeClientSet := &fake.Clientset{
						Fake: k8stesting.Fake{
							ReactionChain: []k8stesting.Reactor{
								&k8stesting.SimpleReactor{Verb: "list", Resource: "pods", Reaction: errorReaction},
							},
						},
					}
					return fakeClientSet.CoreV1()
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo", nil),
				p: httprouter.Params{},
			},
			want: want{
				status: http.StatusInternalServerError,
				response: models.Response{
					PodCount: 1,
					Message:  "OK",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			podInfo := &PodInfoHandler{
				CoreClient: tt.fields.Clientset,
			}
			podInfo.ListPods(tt.args.w, tt.args.r, tt.args.p)

			resp := tt.args.w.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			apiResp := models.Response{}

			if tt.want.status != resp.StatusCode {
				t.Errorf("ListPods() response = %v, want %v", resp.StatusCode, tt.want.status)
			}

			if len(body) > 0 {
				err := json.Unmarshal(body, &apiResp)
				if err != nil {
					t.Fatal("error unmarshalling JSON response")
				}

				if apiResp.PodCount != tt.want.response.PodCount {
					t.Errorf("ListPods() pod count = %v, want %v", apiResp.PodCount, tt.want.response.PodCount)
				}

				if apiResp.Message != tt.want.response.Message {
					t.Errorf("ListPods() response message = %v, want %v", apiResp.Message, tt.want.response.Message)
				}
			}
		})
	}
}

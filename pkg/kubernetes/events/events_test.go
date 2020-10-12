package events

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	defer GinkgoRecover()
	kubeClient := fake.NewSimpleClientset()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "foo",
			UID:       "bar",
		},
	}

	eventRecorder := GenericEventRecorder()
	if err := eventRecorder.Initialize(pod, kubeClient, "test"); err != nil {
		GinkgoT().Fatal("Error initializing event recorder")
	}
}

func TestGenericEventRecording(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(GenericEventRecorder().object)
	assert.NotNil(GenericEventRecorder().recorder)
	assert.NotNil(GenericEventRecorder().watcher)

	events := GenericEventRecorder().watcher.ResultChan()

	GenericEventRecorder().NormalEvent("TestReason", "Test message")
	<-events

	GenericEventRecorder().WarnEvent("TestReason", "Test message")
	<-events

	GenericEventRecorder().ErrorEvent(errors.New("test"), "TestReason", "Test message")
	<-events
}

func TestSpecificEventRecording(t *testing.T) {
	assert := assert.New(t)

	kubeClient := fake.NewSimpleClientset()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "foo",
			UID:       "bar",
		},
	}

	eventRecorder, err := NewEventRecorder(pod, kubeClient, "test")

	assert.Nil(err)
	assert.NotNil(eventRecorder.object)
	assert.NotNil(eventRecorder.recorder)
	assert.NotNil(eventRecorder.watcher)

	events := eventRecorder.watcher.ResultChan()

	eventRecorder.NormalEvent("TestReason", "Test message")
	<-events

	eventRecorder.WarnEvent("TestReason", "Test message")
	<-events

	eventRecorder.ErrorEvent(errors.New("test"), "TestReason", "Test message")
	<-events
}

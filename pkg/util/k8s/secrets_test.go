package k8sutil_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/samsung-cnct/cma-aws/pkg/util/k8s"

	_ "github.com/samsung-cnct/cma-aws/pkg/util/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkingSSHSecretCRUDer struct{}

func (dummy *WorkingSSHSecretCRUDer) Create(*corev1.Secret) (*corev1.Secret, error) {
	return &corev1.Secret{}, nil
}

func (dummy *WorkingSSHSecretCRUDer) Delete(name string, options *v1.DeleteOptions) error {
	return nil
}

func (dummy *WorkingSSHSecretCRUDer) Get(name string, options v1.GetOptions) (*corev1.Secret, error) {
	output := corev1.Secret{}
	output.Data = make(map[string][]byte)
	output.Type = corev1.SecretTypeSSHAuth
	output.Data[corev1.SSHAuthPrivateKey] = []byte("Hi Mom")
	return &output, nil
}

type WrongTypeSSHSecretCRUDer struct{}

func (dummy *WrongTypeSSHSecretCRUDer) Create(*corev1.Secret) (*corev1.Secret, error) {
	return &corev1.Secret{}, nil
}

func (dummy *WrongTypeSSHSecretCRUDer) Delete(name string, options *v1.DeleteOptions) error {
	return nil
}

func (dummy *WrongTypeSSHSecretCRUDer) Get(name string, options v1.GetOptions) (*corev1.Secret, error) {
	output := corev1.Secret{}
	output.Data = make(map[string][]byte)
	output.Data["something"] = []byte("Hi Not Mom")
	return &output, nil
}

type EmptySSHSecretCRUDer struct{}

func (dummy *EmptySSHSecretCRUDer) Create(*corev1.Secret) (*corev1.Secret, error) {
	return &corev1.Secret{}, nil
}

func (dummy *EmptySSHSecretCRUDer) Delete(name string, options *v1.DeleteOptions) error {
	return nil
}

func (dummy *EmptySSHSecretCRUDer) Get(name string, options v1.GetOptions) (*corev1.Secret, error) {
	output := corev1.Secret{}
	output.Data = make(map[string][]byte)
	output.Type = corev1.SecretTypeSSHAuth
	output.Data[corev1.SSHAuthPrivateKey] = []byte("")
	return &output, nil
}

type BrokenSecretCRUDer struct{}

func (dummy *BrokenSecretCRUDer) Create(*corev1.Secret) (*corev1.Secret, error) {
	return nil, fmt.Errorf("A create error happened")
}

func (dummy *BrokenSecretCRUDer) Delete(name string, options *v1.DeleteOptions) error {
	return fmt.Errorf("A delete error happened")
}

func (dummy *BrokenSecretCRUDer) Get(name string, options v1.GetOptions) (*corev1.Secret, error) {
	return nil, fmt.Errorf("A Get error happened")
}

var _ = Describe("K8S Util Secret Functions", func() {
	Context("when dealing with SSH secrets", func() {
		var (
			sshSecret k8sutil.SSHSecret
		)
		Context("getting a secret", func() {
			Context("when a get works", func() {
				var (
					err    error
					secret []byte
				)
				BeforeEach(func() {
					sshSecret = k8sutil.NewSSHSecret(&WorkingSSHSecretCRUDer{})
					secret, err = sshSecret.Get("something")
				})
				It("should return the secret", func() {
					Expect(secret).To(Equal([]byte("Hi Mom")))
				})
				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when you get a secret that's the wrong type", func() {
				var (
					err    error
					secret []byte
				)
				BeforeEach(func() {
					sshSecret = k8sutil.NewSSHSecret(&WrongTypeSSHSecretCRUDer{})
					secret, err = sshSecret.Get("something")
				})
				It("should have a nil", func() {
					Expect(secret).To(HaveLen(0))
				})
				It("should have an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when you get a secret that's empty", func() {
				var (
					err    error
					secret []byte
				)
				BeforeEach(func() {
					sshSecret = k8sutil.NewSSHSecret(&EmptySSHSecretCRUDer{})
					secret, err = sshSecret.Get("something")
				})
				It("should have a nil", func() {
					Expect(secret).To(HaveLen(0))
				})
				It("should have an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when you get an error", func() {
				var (
					err    error
					secret []byte
				)
				BeforeEach(func() {
					sshSecret = k8sutil.NewSSHSecret(&BrokenSecretCRUDer{})
					secret, err = sshSecret.Get("something")
				})
				It("should have a nil", func() {
					Expect(secret).To(HaveLen(0))
				})
				It("should have an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})

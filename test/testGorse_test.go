package main_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zhenghaoz/gorse/client"
)

var _ = Describe("TestGorse", func() {
	It("test", func() {
		gorse := client.NewGorseClient("http://127.0.0.1:8087", "")

		scores, err := gorse.DeleteUser(context.Background(), "40310ed2-70d7-44b9-8eb0-e8b17c48b5ab")
		Expect(err).ToNot(HaveOccurred())
		_ = scores
	})
})

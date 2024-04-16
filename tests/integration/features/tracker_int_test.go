package features_test

import (
	"fmt"

	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"

	dsciv1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/dscinitialization/v1"
	featurev1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/features/v1"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/feature"
	"github.com/opendatahub-io/opendatahub-operator/v2/tests/integration/features/fixtures"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Feature tracking capability - for keeping relationship between features enabled by operator and resources created in the cluster", func() {
	Context("ensuring feature trackers indicate status and phase", func() {

		const appNamespace = "default"

		var (
			dsci *dsciv1.DSCInitialization
		)

		BeforeEach(func() {
			dsci = fixtures.NewDSCInitialization("default")
		})

		It("should indicate successful installation in FeatureTracker", func() {
			featuresHandler := feature.ClusterFeaturesHandler(dsci, func(handler *feature.FeaturesHandler) error {
				verificationFeatureErr := feature.CreateFeature("always-working-feature").
					For(handler).
					UsingConfig(envTest.Config).
					Load()

				Expect(verificationFeatureErr).ToNot(HaveOccurred())

				return nil
			})

			// when
			Expect(featuresHandler.Apply()).To(Succeed())

			// then
			featureTracker, err := fixtures.GetFeatureTracker(envTestClient, appNamespace, "always-working-feature")
			Expect(err).ToNot(HaveOccurred())
			Expect(*featureTracker.Status.Conditions).To(ContainElement(
				MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(conditionsv1.ConditionAvailable),
					"Status": Equal(corev1.ConditionTrue),
					"Reason": Equal(string(featurev1.FeatureCreated)),
				}),
			))
		})

		It("should indicate failure in preconditions", func() {
			// given
			featuresHandler := feature.ClusterFeaturesHandler(dsci, func(handler *feature.FeaturesHandler) error {
				verificationFeatureErr := feature.CreateFeature("precondition-fail").
					For(handler).
					UsingConfig(envTest.Config).
					PreConditions(func(f *feature.Feature) error {
						return fmt.Errorf("during test always fail")
					}).
					Load()

				Expect(verificationFeatureErr).ToNot(HaveOccurred())

				return nil
			})

			// when
			Expect(featuresHandler.Apply()).ToNot(Succeed())

			// then
			featureTracker, err := fixtures.GetFeatureTracker(envTestClient, appNamespace, "precondition-fail")
			Expect(err).ToNot(HaveOccurred())
			Expect(*featureTracker.Status.Conditions).To(ContainElement(
				MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(conditionsv1.ConditionDegraded),
					"Status": Equal(corev1.ConditionTrue),
					"Reason": Equal(string(featurev1.PreConditions)),
				}),
			))
		})

		It("should indicate failure in post-conditions", func() {
			// given
			featuresHandler := feature.ClusterFeaturesHandler(dsci, func(handler *feature.FeaturesHandler) error {
				verificationFeatureErr := feature.CreateFeature("post-condition-failure").
					For(handler).
					UsingConfig(envTest.Config).
					PostConditions(func(f *feature.Feature) error {
						return fmt.Errorf("during test always fail")
					}).
					Load()

				Expect(verificationFeatureErr).ToNot(HaveOccurred())

				return nil
			})

			// when
			Expect(featuresHandler.Apply()).ToNot(Succeed())

			// then
			featureTracker, err := fixtures.GetFeatureTracker(envTestClient, appNamespace, "post-condition-failure")
			Expect(err).ToNot(HaveOccurred())
			Expect(*featureTracker.Status.Conditions).To(ContainElement(
				MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(conditionsv1.ConditionDegraded),
					"Status": Equal(corev1.ConditionTrue),
					"Reason": Equal(string(featurev1.PostConditions)),
				}),
			))
		})

		It("should correctly indicate source in the feature tracker", func() {
			// given
			featuresHandler := feature.ClusterFeaturesHandler(dsci, func(handler *feature.FeaturesHandler) error {
				emptyFeatureErr := feature.CreateFeature("always-working-feature").
					For(handler).
					UsingConfig(envTest.Config).
					Load()

				Expect(emptyFeatureErr).ToNot(HaveOccurred())

				return nil
			})

			// when
			Expect(featuresHandler.Apply()).To(Succeed())

			// then
			featureTracker, err := fixtures.GetFeatureTracker(envTestClient, appNamespace, "always-working-feature")
			Expect(err).ToNot(HaveOccurred())
			Expect(featureTracker.Spec.Source).To(
				MatchFields(IgnoreExtras, Fields{
					"Name": Equal("default-dsci"),
					"Type": Equal(featurev1.DSCIType),
				}),
			)
		})

		It("should correctly indicate app namespace in the feature tracker", func() {
			// given
			featuresHandler := feature.ClusterFeaturesHandler(dsci, func(handler *feature.FeaturesHandler) error {
				emptyFeatureErr := feature.CreateFeature("empty-feature").
					For(handler).
					UsingConfig(envTest.Config).
					Load()

				Expect(emptyFeatureErr).ToNot(HaveOccurred())

				return nil
			})

			// when
			Expect(featuresHandler.Apply()).To(Succeed())

			// then
			featureTracker, err := fixtures.GetFeatureTracker(envTestClient, appNamespace, "empty-feature")
			Expect(err).ToNot(HaveOccurred())
			Expect(featureTracker.Spec.AppNamespace).To(Equal("default"))
		})

	})

})

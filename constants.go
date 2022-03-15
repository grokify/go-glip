package glip

const (
	APIPathGlipFiles               = "/restapi/v1.0/glip/files"
	APIPathGlipGroups              = "/restapi/v1.0/glip/groups"
	APIPathGlipPosts               = "/restapi/v1.0/glip/posts"
	GlipWebhookV1BaseURLProduction = "https://hooks.ringcentral.com/webhook/"
	GlipWebhookV2BaseURLProduction = "https://hooks.ringcentral.com/webhook/v2/"
	GlipWebhookV1BaseURLSandbox    = "https://hooks-glip.devtest.ringcentral.com/webhook/"
	GlipWebhookV2BaseURLSandbox    = "https://hooks-glip.devtest.ringcentral.com/webhook/v2/"
	AttachmentTypeCard             = "Card"
	HeaderValidationToken          = "Validation-Token"
)

var (
	WebhookBaseURL string = "https://hooks.glip.com/webhook/"
)

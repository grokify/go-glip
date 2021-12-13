package glip

const (
	ApiPathGlipFiles               = "/restapi/v1.0/glip/files"
	ApiPathGlipGroups              = "/restapi/v1.0/glip/groups"
	ApiPathGlipPosts               = "/restapi/v1.0/glip/posts"
	GlipWebhookV1BaseURLProduction = "https://hooks.ringcentral.com/webhook/"
	GlipWebhookV2BaseURLProduction = "https://hooks.ringcentral.com/webhook/v2/"
	GlipWebhookV1BaseURLSandbox    = "https://hooks-glip.devtest.ringcentral.com/webhook/"
	GlipWebhookV2BaseURLSandbox    = "https://hooks-glip.devtest.ringcentral.com/webhook/v2/"
	AttachmentTypeCard             = "Card"
)

var (
	WebhookBaseURL string = "https://hooks.glip.com/webhook/"
)

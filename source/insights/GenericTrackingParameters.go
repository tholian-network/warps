package insights

var GenericTrackingParameters []TrackingParameter = []TrackingParameter{
	{
		Domains: []string{"mailchimp.com"},
		Parameters: []string{
			"mc_cid",
			"mc_eid",
		},
	},
	{
		Domains: []string{"adobe.com"},
		Parameters: []string{
			"mkt_tok",
		},
	},
	{
		Domains: []string{"matomo.org"},
		Parameters: []string{
			"mtm_campaign",
			"mtm_cid",
			"mtm_content",
			"mtm_group",
			"mtm_keyword",
			"mtm_medium",
			"mtm_placement",
			"mtm_source",
			"matomo_campaign",
			"matomo_cid",
			"matomo_content",
			"matomo_group",
			"matomo_keyword",
			"matomo_medium",
			"matomo_placement",
			"matomo_source",
		},
	},
	{
		Domains: []string{"microsoft.com"},
		Parameters: []string{
			"cvid",
			"msclkid",
			"oicd",
		},
	},
	{
		Domains: []string{"piwik.org"},
		Parameters: []string{
			"pk_campaign",
			"pk_cid",
			"pk_content",
			"pk_kwd",
			"pk_keyword",
			"pk_medium",
			"piwik_campaign",
			"piwik_cid",
			"piwik_content",
			"piwik_kwd",
			"piwik_keyword",
			"piwik_medium",
		},
	},
	{
		Domains: []string{"google.com"},
		Parameters: []string{
			"utm_campaign",
			"utm_cid",
			"utm_content",
			"utm_medium",
			"utm_name",
			"utm_reader",
			"utm_referrer",
			"utm_social",
			"utm_social-type",
			"utm_source",
			"utm_term",
		},
	},
}

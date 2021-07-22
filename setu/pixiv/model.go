package pixiv

import "time"

type ErrResult struct {
	HasError bool `json:"has_error"`
	Errors   struct {
		System struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"system"`
	} `json:"errors"`
	Error string `json:"error"`
}

type LoginSuccessResult struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ProfileImageUrls struct {
			Px16X16   string `json:"px_16x16"`
			Px50X50   string `json:"px_50x50"`
			Px170X170 string `json:"px_170x170"`
		} `json:"profile_image_urls"`
		ID                     string `json:"id"`
		Name                   string `json:"name"`
		Account                string `json:"account"`
		MailAddress            string `json:"mail_address"`
		IsPremium              bool   `json:"is_premium"`
		XRestrict              int    `json:"x_restrict"`
		IsMailAuthorized       bool   `json:"is_mail_authorized"`
		RequirePolicyAgreement bool   `json:"require_policy_agreement"`
	} `json:"user"`
	Response struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ProfileImageUrls struct {
				Px16X16   string `json:"px_16x16"`
				Px50X50   string `json:"px_50x50"`
				Px170X170 string `json:"px_170x170"`
			} `json:"profile_image_urls"`
			ID                     string `json:"id"`
			Name                   string `json:"name"`
			Account                string `json:"account"`
			MailAddress            string `json:"mail_address"`
			IsPremium              bool   `json:"is_premium"`
			XRestrict              int    `json:"x_restrict"`
			IsMailAuthorized       bool   `json:"is_mail_authorized"`
			RequirePolicyAgreement bool   `json:"require_policy_agreement"`
		} `json:"user"`
	} `json:"response"`
}
type IllustResult struct {
	Illusts []struct {
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Type      string `json:"type"`
		ImageUrls struct {
			SquareMedium string `json:"square_medium"`
			Medium       string `json:"medium"`
			Large        string `json:"large"`
		} `json:"image_urls"`
		Caption  string `json:"caption"`
		Restrict int    `json:"restrict"`
		User     struct {
			ID               int    `json:"id"`
			Name             string `json:"name"`
			Account          string `json:"account"`
			ProfileImageUrls struct {
				Medium string `json:"medium"`
			} `json:"profile_image_urls"`
			IsFollowed bool `json:"is_followed"`
		} `json:"user"`
		Tags []struct {
			Name           string `json:"name"`
			TranslatedName string `json:"translated_name"`
		} `json:"tags"`
		Tools          []string    `json:"tools"`
		CreateDate     time.Time   `json:"create_date"`
		PageCount      int         `json:"page_count"`
		Width          int         `json:"width"`
		Height         int         `json:"height"`
		SanityLevel    int         `json:"sanity_level"`
		XRestrict      int         `json:"x_restrict"`
		Series         interface{} `json:"series"`
		MetaSinglePage struct {
			OriginalImageURL string `json:"original_image_url"`
		} `json:"meta_single_page,omitempty"`
		MetaPages []struct {
			ImageUrls struct {
				Original string `json:"original"`
			} `json:"image_urls"`
		} `json:"meta_pages"`
		TotalView      int  `json:"total_view"`
		TotalBookmarks int  `json:"total_bookmarks"`
		IsBookmarked   bool `json:"is_bookmarked"`
		Visible        bool `json:"visible"`
		IsMuted        bool `json:"is_muted"`
	} `json:"illusts"`
	SearchSpanLimit int `json:"search_span_limit"`
}

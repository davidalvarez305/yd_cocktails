const clickIdKeys = ["gclid", "gbraid", "wbraid", "msclkid", "fbclid", "li_fat_id"];

export class MarketingHelper {
    constructor() {
        this.landingPage = new URL(window.location.href);
        this.referrer = document.referrer;
        this.user = JSON.parse(localStorage.getItem("user")) || null;

        this.language = navigator.language || navigator.userLanguage;

        this.longitude = null;
        this.latitude = null;

        this.clickId = null;
        this.facebookClickId = null;
        this.facebookClientId = null;

        this.userAgent = navigator.userAgent;

        if (this.user) {
            if (this.user.landingPage) this.landingPage = new URL(this.user.landingPage);
            if (this.user.referrer) this.referrer = this.user.referrer;
        }

        // Get Click ID
        if (this.isPaid(this.landingPage.searchParams)) {
            this.clickId = this.getClickId(this.landingPage.searchParams);
        };

        // Get FB Click ID
        const fbClickId = this.landingPage.searchParams.get("fbclid");
        if (fbClickId) this.facebookClickId = fbClickId;

        // Get FB Client ID
        const fbp = this.getCookie("_fbp");
        if (fbp) this.facebookClientId = fbp;

        this.data = new FormData();
    }

    populate() {
        if (this.clickId) this.data.set("click_id", this.clickId);
        if (this.facebookClickId) this.data.set("facebook_click_id", this.facebookClickId);
        if (this.facebookClientId) this.data.set("facebook_client_id", this.facebookClientId);

        if (this.landingPage) this.data.set("landing_page", this.landingPage);
        if (this.referrer) this.data.set("referrer", this.referrer);
        if (this.language) this.data.set("language", this.language);

        // Append source, medium, and channel based on URL or referrer
        const source = this.landingPage.searchParams.get("source") || this.getSource();
        const medium = this.landingPage.searchParams.get("medium") || this.getMedium();
        const channel = this.landingPage.searchParams.get("channel") || this.getChannel();
        const utm_source = this.landingPage.searchParams.get("utm_source");

        if (medium) this.data.set("medium", medium);
        if (channel) this.data.set("channel", channel);
        if (source && !utm_source) this.data.set("source", source);
        if (utm_source) this.data.set("source", utm_source);

        // Handle geolocation (conditionally append if available)
        if (this.longitude) this.data.set("longitude", this.longitude);
        if (this.latitude) this.data.set("latitude", this.latitude);

        // Append all marketing parameters
        this.landingPage.searchParams.forEach((value, key) => value && this.data.set(key, value));
    }

    getCookie(name) {
        const cookies = document.cookie.split(';');

        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(name + '=')) {
                return cookie.substring(name.length + 1);
            }
        }

        return null;
    }

    getSource() {
        let url;
        try {
            url = new URL(this.referrer);
        } catch (error) {
            return "";
        }

        let host = url.hostname.toLowerCase();

        if (host.startsWith("www.")) {
            host = host.slice(4);
        }

        const parts = host.split(".");

        // Handle cases like ftp.google.com or ads.google.com
        if (
            parts.length > 2 &&
            !["com", "net", "org", "edu", "gov", "mil", "int"].includes(
                parts[parts.length - 2]
            )
        ) {
            host = parts.slice(-3).join(".");
        } else if (parts.length > 1) {
            // Check if the last part is a two-letter country code
            const lastPart = parts[parts.length - 1];
            if (lastPart.length === 2 && lastPart !== "co") {
                // 'co' is a special case like .co.uk, .co.in, etc.
                host = parts.slice(-3).join(".");
            } else {
                host = parts.slice(-2).join(".");
            }
        }

        return host;
    }

    getClickId() {
        for (const key of clickIdKeys) {
            if (this.landingPage.searchParams.has(key)) return this.landingPage.searchParams.get(key);
        }
    }

    isPaid() {
        for (const key of clickIdKeys) {
            if (this.landingPage.searchParams.has(key)) {
                return true;
            }
        }

        return false;
    }

    getMedium() {
        // No referrer means the user accessed the website directly
        if (this.referrer.length === 0) return "direct";

        // Non-empty referrer and no querystring === organic
        if (this.landingPage.searchParams.size === 0) return "organic";

        // Paid ads
        if (this.isPaid(this.landingPage.searchParams)) return "paid";

        // Querystring + non-empty referrer and no click id === referral
        return "referral";
    }

    getUserLocation() {
        const options = {
            enableHighAccuracy: true,
            timeout: 3000,
            maximumAge: 0,
        };

        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                function (position) {
                    this.latitude = position.coords.latitude;
                    this.longitude = position.coords.longitude;
                },
                function (error) {
                    console.error("Error getting user location:", error.message);
                },
                options
            );
        } else {
            console.error("Geolocation is not supported by your browser.");
        }
    }

    getChannel() {
        const displayNetworks = ["googleads.g.doubleclick.net"];

        const searchEngines = [
            { domain: "google" },
            { domain: "bing" },
            { domain: "yahoo" },
            { domain: "ecosia" },
            { domain: "duckduckgo" },
            { domain: "yandex" },
            { domain: "baidu" },
            { domain: "naver" },
            { domain: "ask.com" },
            { domain: "adsensecustomsearchads" },
            { domain: "aol" },
            { domain: "brave" },
        ];

        const majorSocialNetworks = [
            "facebook",
            "instagram",
            "twitter",
            "linkedin",
            "pinterest",
            "snapchat",
            "reddit",
            "whatsapp",
            "wechat",
            "telegram",
            "discord",
            "vkontakte",
            "weibo",
            "line",
            "kakaotalk",
            "qq",
            "viber",
            "telegram",
            "tumblr",
            "flickr",
            "meetup",
            "tagged",
            "badoo",
            "myspace",
        ];

        const majorVideoPlatforms = [
            "youtube",
            "tiktok",
            "vimeo",
            "dailymotion",
            "twitch",
            "bilibili",
            "youku",
            "rutube",
            "vine",
            "peertube",
            "ig tv",
            "veoh",
            "metacafe",
            "vudu",
            "vidyard",
            "rumble",
            "bit chute",
            "brightcove",
            "viddler",
            "vzaar",
        ];

        // Check display platforms
        for (let platform of displayNetworks) {
            if (this.referrer.includes(platform)) {
                return "display";
            }
        }

        // Check search engines
        for (let engine of searchEngines) {
            if (this.referrer.includes(engine.domain)) {
                return "search";
            }
        }

        // Check social networks
        for (let network of majorSocialNetworks) {
            if (this.referrer.includes(network)) {
                return "social";
            }
        }

        // Check video platforms
        for (let platform of majorVideoPlatforms) {
            if (this.referrer.includes(platform)) {
                return "video";
            }
        }

        return "other";
    }
}
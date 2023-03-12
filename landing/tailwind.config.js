const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
    darkMode: "class",
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                "bgSet": {
                    "light": "#FFFFFF",
                    "dark": "#0A1018"
                },
                "blackSet": {
                    "light": "#000000",
                    "dark": "#FFFFFF"
                },
                "fontGrey": {
                    "light": "#3D3D3D",
                    "dark": "#E2E2E2"
                },
                "orangeSet": {
                    "light": "#DF7902",
                    "dark": "#F59300"
                },
                "pizazz": "#F59300",
                "electricViolet": "#C41CFF",
            },
        },
        fontFamily: {
            "sans": ["Poppins", ...defaultTheme.fontFamily.sans]
        },
        fontSize: {
            'largeTitle': ["24px", {
                "lineHeight": "28px",
                "fontWeight": "700"
            }],
            'largeTitle2': ["58px", {
                "lineHeight": "68px",
                "fontWeight": "700"
            }],
            'largeTitle3': ["62px", {
                "lineHeight": "72px",
                "fontWeight": "700"
            }],
            'headline': ["14px", {
                "lineHeight": "20px",
                "fontWeight": "400"
            }],
            'headline2': ["16px", {
                "lineHeight": "20px",
                "fontWeight": "400"
            }],
            'headline3': ["20px", {
                "lineHeight": "24px",
                "fontWeight": "400"
            }],
            "button": ["16px", {
                "lineHeight": "20px",
                "fontWeight": "500"
            }],
            "body3": ["12px", {
                "lineHeight": "16px",
                "fontWeight": "400"
            }],
        },
    },
    plugins: [],
}

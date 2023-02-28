const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
    darkMode: "class",
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                "bg": {
                    "light": "#FFFFFF",
                    "dark": "#0A1018"
                }
            },
        },
        fontFamily: {
            "sans": ["Poppins", ...defaultTheme.fontFamily.sans]
        },
    },
    plugins: [],
}

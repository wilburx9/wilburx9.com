import sun from './images/sun.svg'
import moon from './images/moon.svg'
import logo from './images/logo.svg'
import React, {useEffect, useState} from "react";

export default function Header() {
    const [theme, setTheme] = useDarkMode()
    let isDarkTheme = theme === "dark"
    return (
        <div className="flex flex-row w-full">
            <img src={logo} alt="Logo"/>
            <img src={isDarkTheme ? sun : moon} alt={`Switch to ${isDarkTheme ? "light" : "dark"} mode`}
                 onClick={() => setTheme(isDarkTheme ? "light" : "dark")}/>
        </div>
    );
}

function useDarkMode() {
    const [theme, setTheme] = useState(typeof window !== "undefined" ? localStorage.theme : "dark")

    useEffect(() => {
        const colorTheme = theme === "dark" ? "light" : "dark"
        const root = window.document.documentElement

        root.classList.remove(colorTheme)
        root.classList.add(theme)

        if (typeof window != "undefined") {
            localStorage.setItem("theme", theme)
        }
    }, [theme])

    return [theme, setTheme] as const
}
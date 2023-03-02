import React from 'react';
import './App.css';
import Header from "./Header";
import Content from "./Content";
import Footer from "./Footer";

export default function App() {
    return (
        <div className="h-screen w-screen bg-bgSet-light dark:bg-bgSet-dark">
            <div className="max-w-[1192px] h-screen  flex flex-col mx-auto">
                <Header/>
                <Content/>
                <Footer/>
            </div>
        </div>
    )
}


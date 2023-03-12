import React from 'react';
import './App.css';
import Header from "./Header";
import Content from "./Content";
import Footer from "./Footer";

export default function App() {
    return (
        <div className="h-screen w-screen bg-bgSet-light dark:bg-bgSet-dark">
            <div className="max-w-screen-xl h-screen flex flex-col mx-auto px-4 md:px-20 xl:px-4 2xl:px-0">
                <Header/>
                <Content/>
                <Footer/>
            </div>
        </div>
    )
}


import React from 'react';
import './App.css';
import Header from "./Header";
import Content from "./Content";
import Footer from "./Footer";

export default function App() {
    return (
        <div className="h-screen w-screen bg-bg-light dark:bg-bg-dark flex flex-col">
            <Header/>
            <Content/>
            <Footer/>
        </div>
    )
}


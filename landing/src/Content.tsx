export default function Content() {
    return (
        <div className="grow flex flex-col-reverse xl:flex-row gap-x-6 justify-center gap-y-4 pt-3">
            <Left/>
            <Right/>
        </div>
    
    );
}

function Left() {
    return (
        <div className="grow basis-0 flex flex-col justify-center px-0 lg:px-20 xl:px-0">
            <h1 className="text-center md:text-left  text-blackSet-light dark:text-blackSet-dark text-largeTitle md:text-largeTitle2 xl:text-largeTitle3">
                I am
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-pizazz to-electricViolet"> Wilberforce
                    </span>,<br/>a Software Engineer and Visual Storyteller.</h1>
            <h2 className="text-center md:text-left pt-2 text-headline md:text-headline2 xl:text-headline3 text-fontGrey-light dark:text-fontGrey-dark">
                Seasoned software engineer and avid photographer. Combining technical expertise with a creative eye
                to deliver impactful digital experiences.
            </h2>
            <div className="flex flex-row gap-x-4 pt-8">
                <a href="/blog" className="grow basis-0 bg-bgSet-dark dark:bg-bgSet-light bg-cover
                text-bgSet-light dark:text-bgSet-dark rounded-[12px] text-button py-4
                hover:bg-orangeSet-light dark:hover:bg-orangeSet-dark hover:text-bgSet-light dark:hover:text-bgSet-light
                focus:bg-orangeSet-light dark:focus:bg-orangeSet-dark focus:text-bgSet-light dark:focus:text-bgSet-light
                active:bg-orangeSet-light dark:active:bg-orangeSet-dark active:text-bgSet-light dark:active:text-bgSet-light
                text-center no-underline">
                    Blog
                </a>
                <a href="https://www.linkedin.com/in/wilburx9"
                   className="grow basis-0 border border-blackSet-light bg-cover dark:border-blackSet-dark rounded-[12px]
                text-bgSet-dark dark:text-bgSet-light text-button py-4 no-underline
                hover:bg-orangeSet-light dark:hover:bg-orangeSet-dark hover:text-bgSet-light dark:hover:text-bgSet-light
                focus:bg-orangeSet-light dark:focus:bg-orangeSet-dark focus:text-bgSet-light dark:focus:text-bgSet-light
                active:bg-orangeSet-light dark:active:bg-orangeSet-dark active:text-bgSet-light dark:active:text-bgSet-light
                hover:border-transparent focus:border-transparent active:border-transparent
                hover:dark:border-transparent focus:dark:border-transparent active:dark:border-transparent text-center">
                    LinkedIn
                </a>
            </div>
        </div>
    )
}

function Right() {
    return (
        <div className="grow basis-0">
            <div className="group flex justify-center items-center w-full h-full bg-contain bg-no-repeat bg-center
            bg-[url('./images/pattern_light.svg')]
            dark:bg-[url('./images/pattern_dark.svg')]
            hover:bg-[url('./images/pattern.svg')]
            hover:dark:bg-[url('./images/pattern.svg')]">
                <div className="w-[43.5%] h-[43.5%] bg-auto bg-no-repeat bg-center bg-contain mx-auto my-auto
            bg-[url('./images/me_bw.png')] group-hover:bg-[url('./images/me.png')]"/>
            </div>
        </div>
    )
}
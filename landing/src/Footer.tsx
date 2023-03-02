import github from './images/github.svg'
import twitter from './images/twitter.svg'
import instagram from './images/instagram.svg'

export default function Footer() {
    let year = new Date().getFullYear()
    let socials = [
        {name: "GitHub", icon: github, url: "github.com/wilburt"},
        {name: "Twitter", icon: twitter, url: "twitter.com/wilburx09"},
        {name: "Instagram", icon: instagram, url: "instagram.com/wilburx9"},
    ]
    
    return (
        <div className="flex flex-row w-full justify-between pb-[70px] pt-5 items-center">
            <div className="flex flex-row gap-x-6 items-center">
                {socials.map(s => <a href={`https://${s.url}`}><img src={s.icon} alt={s.icon}/></a>)}
            </div>
            <p className="text-body3 font-normal leading-16 text-bgSet-dark dark:text-bgSet-light">
                {`Copyright ${String.fromCharCode(169)} ${year} All rights reserved.`}
            </p>
        </div>
    )
}

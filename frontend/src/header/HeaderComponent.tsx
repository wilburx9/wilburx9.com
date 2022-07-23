import {
  Flex,
  Icon,
  IconButton,
  Image, Link,
  useColorMode,
  useColorModeValue
} from "@chakra-ui/react";
import React from "react";
import avatar from './avatar.png'
import {IconType} from "react-icons";
import {VscGithub} from "react-icons/vsc";
import {ImLinkedin, ImTwitter} from "react-icons/im";
import {RiMoonFill, RiSunFill} from "react-icons/ri";

class Social {
  name: string
  url: string
  icon: IconType

  constructor(name: string, icon: IconType, url: string) {
    this.name = name
    this.url = url;
    this.icon = icon;
  }
}

let socials: Social[] = [
  new Social("LinkedIn", ImLinkedin, "https://www.linkedin.com/in/wilburx9"),
  new Social("Twitter", ImTwitter, "https://twitter.com/wilburx09"),
  new Social("Github", VscGithub, "https://github.com/wilburt"),
]

export const HeaderComponent = () => (
  <Flex w='full' alignItems='center' pt={16}>
    <Image src={avatar} boxSize='64px' alt="Wilbur's Avatar"/>
    <Flex flex='1' justifyContent='end' alignItems='center'>
      {socials.map(e => <Link href={e.url} isExternal>
        <Icon as={e.icon!} marginStart={9} boxSize={5} display='block'/>
      </Link>)}
      <ThemeIcon/>
    </Flex>
  </Flex>
)


const ThemeIcon: React.FC = () => {
  const {toggleColorMode} = useColorMode()
  const text = useColorModeValue("dark", "light")
  const SwitchIcon = useColorModeValue(RiMoonFill, RiSunFill)

  return (
    <IconButton
      size="lg"
      variant="ghost"
      color="current"
      marginStart={6}
      onClick={toggleColorMode}
      icon={<SwitchIcon/>}
      aria-label={`Switch to ${text} mode`}
    />
  )
}

import {
  Box,
  Flex, Heading,
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
import {IoArrowDown} from "react-icons/io5";

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

export const TopComponent = () => (
  <Box w='full' h='100vh'>
    <Flex direction='column' h='100%'>
      <TopSection/>
      <Flex direction='column' flex='1' justify='space-between'>
        This is a box
        <BottomSection/>
      </Flex>
    </Flex>
  </Box>
)

const TopSection = function () {
  const {toggleColorMode} = useColorMode()
  const text = useColorModeValue("dark", "light")
  const SwitchIcon = useColorModeValue(RiMoonFill, RiSunFill)
  return <Flex alignItems='center' pt='4vh'>
    <Image src={avatar} boxSize='64px' alt="Wilbur's Avatar"/>
    <Flex flex='1' justifyContent='end' alignItems='center'>
      {socials.map(e => <Link href={e.url} isExternal>
        <Icon as={e.icon!} marginStart={9} boxSize={5} display='block'/>
      </Link>)}
      <IconButton
        size="lg"
        variant="ghost"
        color="current"
        marginStart={6}
        onClick={toggleColorMode}
        isRound={true}
        icon={<SwitchIcon/>}
        aria-label={`Switch to ${text} mode`}
      />
    </Flex>
  </Flex>
}

const MiddleSection = function () {
  
}

const BottomSection = function () {
  const color = useColorModeValue('#1A1B22', 'white')
  return <Flex direction='row' align='end' w='full' mb={4}>
    <IconButton
      isRound={true}
      mb={10}
      variant="ghost"
      color="current"
      fontSize='32px'
      minH='70px'
      minW='70px'
      icon={<IoArrowDown/>}
      aria-label={'Scroll down'}/>
    <Heading as='h1' textAlign='right' flex='1' fontSize='10vh' color={color}>
      I am Wilberforce{<br/>}
      Software Engineer{<br/>}
      at <span style={{color: "#8D949D"}}>Mindvalley</span>
    </Heading>
  </Flex>;
}

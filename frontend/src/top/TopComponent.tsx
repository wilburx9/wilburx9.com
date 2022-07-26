import {
  Box, Fade,
  Flex, Heading,
  Icon,
  IconButton,
  Image, Link, Text,
  useColorMode,
  useColorModeValue
} from "@chakra-ui/react";
import React, {useContext, useEffect} from "react";
import avatar from './avatar.png'
import {IconType} from "react-icons";
import {VscGithub} from "react-icons/vsc";
import {ImLinkedin, ImTwitter} from "react-icons/im";
import {RiArrowRightSLine, RiMoonFill, RiSunFill} from "react-icons/ri";
import {IoArrowDown} from "react-icons/io5";
import {DataContext, DataValue} from "../DataProvider";
import {ArticleModel} from "../articles/ArticleModel";

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
      <MiddleSection/>
      <BottomSection/>
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
      {socials.map(e => <Link href={e.url} isExternal key={e.url}>
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
  const {articles} = useContext<DataValue>(DataContext)

  useEffect(() => {
  }, [articles])

  let firstArticle: ArticleModel | null = articles[0]
  const postColor = useColorModeValue('black', 'white')
  return <Fade in={firstArticle != null} unmountOnExit>
    <Box width={{base: '20vw', md: '70vw', 'lg': '40vw'}}>
      <Box mt='5vh'>
        <Icon as={RiArrowRightSLine} verticalAlign='middle' color={'#596065'}/>
        <Text as='span' verticalAlign='middle' fontWeight='medium' fontSize='md' color={'#8D949D'}>Read my latest
          article</Text>
      </Box>
      <Box _hover={{opacity: 0.7}}>
        <Link href={firstArticle?.url} _hover={{textDecoration: 'none'}}>
          <Text textAlign='start' fontWeight='normal' fontSize='xl' color={postColor} mt={4}
                noOfLines={5}>{firstArticle?.title}</Text>
        </Link>
      </Box>
    </Box>
  </Fade>
}

const BottomSection = function () {
  const color = useColorModeValue('#1A1B22', 'white')
  return <Flex flex='1' direction={{base: 'column-reverse', md: 'row'}} align='end' w='full' mb={{base: 0, md: 4}}>
    <IconButton
      isRound={true}
      mb={{base: 6, md: 10}}
      variant="ghost"
      color="current"
      fontSize='32px'
      minH='70px'
      minW='70px'
      icon={<IoArrowDown/>}
      aria-label={'Scroll down'}/>
    <Flex flex='1' justifyContent='flex-end'>
      <Heading display='inline-block' alignSelf='flex-end' as='h1' textAlign='end' fontSize='min(10vw, 10vh)'
               color={color} mb={{base: '5vh', md: 0}}>
        I am Wilberforce{<br/>}
        Software Engineer{<br/>}
        at <span style={{color: "#8D949D"}}>Mindvalley</span>
      </Heading>
    </Flex>
  </Flex>;
}

import React, {useContext} from "react";
import {DataContext} from "../DataProvider";
import {
  Box,
  Heading,
  HStack,
  Icon,
  LinkBox,
  LinkOverlay,
  SimpleGrid,
  Text,
  useColorModeValue,
  VStack
} from "@chakra-ui/react";
import {Language, RepoModel} from "./RepoModel";
import {IconType} from "react-icons";
import {SiDart, SiGo, SiJava, SiKotlin, SiTypescript} from "react-icons/si";
import {DiCss3, DiSwift} from "react-icons/di";
import {IoLogoJavascript} from "react-icons/io";
import {AiFillStar} from "react-icons/ai";
import {RiGitBranchFill} from "react-icons/ri";


export const ReposComponent = () => {
  const {repos} = useContext(DataContext)

  if (!repos || repos.length === 0) return <Box/>

  return (
    <VStack align='start'
            w='full'>
      <Heading pt={20}
               size='xl'
               as='h4'>And here are my open-source <Box
        as='span'
        color='cyan.500'>contributions</Box>.</Heading>
      <SimpleGrid columns={{base: 1, md: 2}}
                  spacing={8}
                  w='full'
                  py={10}>
        {repos.map((repo) =>
          <RepoComponent {...repo} key={repo.id}/>
        )}
      </SimpleGrid>
    </VStack>
  )
}

function RepoComponent(props: RepoModel) {
  // Assign icons to the languages and filter out those with no icon afterwards
  let languages = props.languages.map(function (l) {
    let icon = lngIconMap.get(l.name.toLowerCase())
    return new _Language(l, icon?.icon, icon?.size)
  }).filter((l) => l.icon != null)

  let hoverBorderColor = useColorModeValue('#BFC6CF', '#596065');
  let bgColor = useColorModeValue('#F7F8F8', '#26272D');

  return <LinkBox>
    <LinkOverlay href={props.url}
                 isExternal
                 role="group">
      <Box
        w='full'
        h='full'
        bg={bgColor}
        borderWidth='1.5px'
        borderRadius='lg'
        borderColor={bgColor}
        py={5}
        px={6}
        _groupHover={{borderColor: hoverBorderColor}}>
        <VStack w='full'
                h='full'
                alignItems='flex-start'
                justify='center'>
          <Heading mb={1}
                   size='md'>{props.name}</Heading>
          <Text mb={3}
                fontSize='sm'
                align="start">{props.description}</Text>
          <HStack w='full'
                  justify='space-between'
                  pt={6}>
            <HStack spacing={0}>
              <Icon as={RiGitBranchFill}/>
              <Text fontSize='sm'
                    fontWeight='bold'
                    pl={2}
                    pr={6}>{props.forks}</Text>
              <Icon as={AiFillStar}
                    pr={0.5}/>
              <Text fontSize='sm'
                    fontWeight='bold'
                    pl={2}
                    pr={6}>{props.stars}</Text>
            </HStack>
            <HStack spacing={4}>
              {languages.map((l) => <LanguageComponent {...l} key={l.name}/>)}
            </HStack>
          </HStack>
        </VStack>
      </Box>
    </LinkOverlay>
  </LinkBox>
}

function LanguageComponent(props: _Language) {
  return <Icon as={props.icon!}
               color={props.color}
               w={props.iconSize}
               h={props.iconSize}/>
}

let lngIconMap: Map<string, { icon: IconType, size: number }> = new Map([
  ["java", {icon: SiJava, size: 5}],
  ["kotlin", {icon: SiKotlin, size: 3}],
  ["dart", {icon: SiDart, size: 4}],
  ["go", {icon: SiGo, size: 7}],
  ["swift", {icon: DiSwift, size: 6}],
  ["typescript", {icon: SiTypescript, size: 4}],
  ["javascript", {icon: IoLogoJavascript, size: 4}],
  ["css", {icon: DiCss3, size: 4}],
])

class _Language {
  name: string;
  color: string;
  icon?: IconType;
  iconSize?: number

  constructor(l: Language, icon?: IconType, iconSize?: number) {
    this.name = l.name
    this.color = l.color
    this.icon = icon
    this.iconSize = iconSize
  }
}

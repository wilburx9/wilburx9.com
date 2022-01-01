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
import {
  AiFillStar,
  CgGitFork, DiCss3,
  DiSwift,
  SiGo,
  IoLogoJavascript,
  SiDart,
  SiKotlin,
  SiTypescript, SiJava
} from "react-icons/all";
import {IconType} from "react-icons";


export const ReposComponent = () => {
  const {repos} = useContext(DataContext)

  if (!repos || repos.length === 0) return <Box/>

  return (
    <VStack align='start'>
      <Heading pt={16} size='xl' align="start" fontWeight="black">&#47;&#47; Open-source Projects</Heading>
      <SimpleGrid columns={{base: 1, md: 2, lg: 3}} spacing={10} py={4}>
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

  let hoverBorderColor = useColorModeValue('black', 'white');
  let bgColor = useColorModeValue('gray.100', 'gray.900');

  return <LinkBox>
    <LinkOverlay href={props.url} isExternal role="group">
      <Box
        w='full'
        h='full'
        bg={bgColor}
        borderWidth='1px'
        borderRadius='lg'
        p={4}
        _groupHover={{borderColor: hoverBorderColor}}>
        <VStack w='full' h='full' alignItems='flex-start' justify='center'>
          <Heading mb={1} size='sm' align="start">{props.name}</Heading>
          <Text mb={3} fontSize='sm' align="start">{props.description}</Text>
          <HStack w='full' justify='space-between'>
            <HStack spacing={0}>
              <Icon as={CgGitFork}/>
              <Text fontSize='xs' fontWeight='bold' pr={4}>{props.forks}</Text>
              <Icon as={AiFillStar} pr={0.5}/>
              <Text fontSize='xs' fontWeight='bold' pr={4}>{props.stars}</Text>
            </HStack>
            <HStack>
              {languages.map((l) => <LanguageComponent {...l} key={l.name}/>)}
            </HStack>
          </HStack>
        </VStack>
      </Box>
    </LinkOverlay>
  </LinkBox>
}

function LanguageComponent(props: _Language) {
  return <Icon as={props.icon!} color={props.color} w={props.iconSize} h={props.iconSize}/>
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

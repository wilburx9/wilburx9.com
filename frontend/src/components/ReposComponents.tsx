import React, {useContext, useEffect} from "react";
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
import {Language, RepoModel} from "../models/RepoModel";
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


export const ReposComponents = () => {
  const {fetchRepos, repos} = useContext(DataContext)

  useEffect(() => {
    fetchRepos?.()
  }, [fetchRepos])

  if (!repos) return <Box/>

  return (
    <SimpleGrid columns={{base: 1, md: 2, lg: 3}} spacing={10} p={4}>
      {repos.map((repo) =>
        <Item {...repo} key={repo.name}/>
      )}
    </SimpleGrid>
  )
}

function Item(props: RepoModel) {
  // Assign icons to the languages and filter out those with no icon afterwards
  let languages = props.languages.map(function (l) {
    let icon = lngIconMap.get(l.name.toLowerCase())
    return new Lng(l, icon?.icon, icon?.size)
  }).filter((l) => l.icon != null)

  let borderColor = useColorModeValue('black', 'white');

  return <LinkBox>
    <LinkOverlay href={props.url} isExternal role="group">
      <Box w='full' h='full' borderWidth='1px' borderRadius='lg' p={4} _groupHover={{borderColor: borderColor}}>
        <VStack w='full' h='full' alignItems='flex-start' justify='center'>
          <Heading mb={1} fontSize="16px" align="start" fontWeight="bold">{props.name}</Heading>
          <Text mb={3} fontSize="14px" align="start" fontWeight="regular">{props.description}</Text>
          <HStack w='full' justify='space-between'>
            <HStack spacing={0}>
              <Icon as={CgGitFork}/>
              <Text fontSize='12px' fontWeight='medium' pr={4}>{props.forks}</Text>
              <Icon as={AiFillStar} pr={0.5}/>
              <Text fontSize='12px' fontWeight='medium' pr={4}>{props.stars}</Text>
            </HStack>
            <HStack>
              {languages.map((l) => <Icon as={l.icon!} color={l.color} w={l.iconSize} h={l.iconSize}/>)}
            </HStack>
          </HStack>
        </VStack>
      </Box>
    </LinkOverlay>
  </LinkBox>
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

class Lng {
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
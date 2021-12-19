import {
  Box,
  Heading,
  Image,
  LinkBox, LinkOverlay,
  Text,
  SimpleGrid, Icon, Circle, useColorModeValue, AspectRatio, Flex, VStack,
} from "@chakra-ui/react";
import React, {useContext} from "react";
import {DataContext} from "../DataProvider";
import {ArticleModel} from "../models/ArticleModel";
import {HiArrowRight} from "react-icons/hi";

export const ArticlesComponent = () => {
  const {articles} = useContext(DataContext)

  if (!articles || articles.length === 0) return <Box/>

  return (
    <VStack align='start'>
      <Heading pt={6} size='xl' align="start" fontWeight="black">&#47;&#47;Articles</Heading>
      <SimpleGrid columns={{base: 1, md: 2, lg: 3}} spacing={10} py={4}>
        {articles.map((article) => <ArticleComponent {...article} key={article.title}/>)}
      </SimpleGrid>
    </VStack>
  )
}

function ArticleComponent(props: ArticleModel) {
  return <LinkBox>
    <LinkOverlay href={props.url} isExternal role="group">
      <Box w='full' h='full' borderRadius='lg' borderWidth='1px' pb={5}
           _groupHover={{bg: useColorModeValue('gray.100', 'gray.900')}}>
        <Flex direction='column' align='start' justify='space-between' h='full' spacing={0}>
          <AspectRatio w='full' ratio={3 / 2}>
            <Image
              src={props.thumbnail}
              alt={props.title}
              borderTopRadius='lg'
              objectFit="cover"
              w='full'/>
          </AspectRatio>
          <Heading px={6} pt={6} size='md' align="start">{props.title}</Heading>
          <Text px={6} pt={3} noOfLines={3} fontSize='md' align="start">{props.excerpt}</Text>
          <Box ml='auto'>
            <Circle bg={useColorModeValue('gray.200', 'gray.700')} boxSize={10} mt={4} mx={9}>
              <Icon as={HiArrowRight} color={useColorModeValue('gray.800', 'gray.200')}/>
            </Circle>
          </Box>
        </Flex>
      </Box>
    </LinkOverlay>
  </LinkBox>
}

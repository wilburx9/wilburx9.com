import {
  Box,
  Heading,
  Image,
  LinkBox, LinkOverlay,
  Text,
  SimpleGrid, Icon, Circle, useColorModeValue, AspectRatio, Flex,
} from "@chakra-ui/react";
import React, {useContext, useEffect} from "react";
import {DataContext} from "../DataProvider";
import {ArticleModel} from "../models/ArticleModel";
import {HiArrowRight} from "react-icons/hi";

export const ArticlesComponent = () => {
  const {fetchArticles, articles} = useContext(DataContext)

  useEffect(() => {
    fetchArticles?.()
  }, [fetchArticles])

  if (!articles) return <Box/>

  return (
    <SimpleGrid columns={{base: 1, md: 2, lg: 3}} spacing={10} py={4}>
      {articles.map((article) => <Column {...article} key={article.title}/>)}
    </SimpleGrid>
  )
}

function Column(props: ArticleModel) {
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
          <Heading px={6} pt={6} fontSize="20px" align="start" fontWeight="bold">{props.title}</Heading>
          <Text px={6} pt={3} noOfLines={3} fontSize="16px" align="start" fontWeight="regular">{props.excerpt}</Text>
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

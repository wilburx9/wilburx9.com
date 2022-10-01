import {
  Box,
  Heading,
  Image,
  LinkBox, LinkOverlay,
  SimpleGrid, Icon, AspectRatio, Flex, VStack, Text, useColorModeValue,
} from "@chakra-ui/react";
import React, {useContext} from "react";
import {DataContext} from "../DataProvider";
import {ArticleModel} from "./ArticleModel";
import {FiArrowRight} from "react-icons/fi";

export const ArticlesComponent = () => {
  const {articles} = useContext(DataContext)

  if (!articles || articles.length === 0) return <Box/>

  return (
    <VStack align='start'
            id="articles">
      <Heading pt={20}
               size='xl'
               as='h4'>Here are some of my from <Box
        as='span'
        color='secondary.300'>articles</Box>.</Heading>
      <SimpleGrid columns={{base: 1, md: 2, lg: 3}}
                  spacing={10}
                  py={4}>
        {articles.map((article) => <ArticleComponent {...article} key={article.id}/>)}
      </SimpleGrid>
    </VStack>
  )
}

function ArticleComponent(props: ArticleModel) {
  return <LinkBox mt={10}>
    <LinkOverlay href={props.url}
                 isExternal
                 role="group">
      <Box w='full'
           h='full'
           pb={5}
           _groupHover={{boxShadow: useColorModeValue('2xl', 'dark-lg')}}>
        <Flex direction='column'
              align='start'
              justify='space-between'
              h='full'>
          <AspectRatio w='full'
                       ratio={2.2385}>
            <Image
              src={props.thumbnail}
              alt={props.title}
              objectFit="cover"
              w='full'
              p={2}/>
          </AspectRatio>
          <Heading
            px={3}
            pt={4}
            size='md'
            as='h2'
            fontWeight='semibold'>
            {props.title}
          </Heading>
          <Flex
            mx={3}
            align='center'
            my={4}
            bg='secondary.400'
            borderRadius='md'
            px={4}
            py={2}>
            <Text fontWeight='semibold'
                  color='white'
                  fontSize='xs'>Read Here</Text>
            <Icon ml={6}
                  as={FiArrowRight}
                  boxSize={3}
                  color='white'/>
          </Flex>
        </Flex>
      </Box>
    </LinkOverlay>
  </LinkBox>
}

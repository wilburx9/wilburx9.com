import {Box, Flex, Heading, HStack, Image, StackDivider, Text, VStack} from "@chakra-ui/react";
import React, {useContext, useEffect} from "react";
import {DataContext} from "./DataProvider";
import {Utils} from "./Utils";
import {Article} from "./models/Article";

export const ArticleComponent = () => {
  const {fetchArticles, articles} = useContext(DataContext)

  useEffect(() => {
    fetchArticles?.()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fetchArticles])

  if (!articles) return <Box/>

  let groups = Utils.groupArray(articles, 2);
  return (
    <HStack divider={<StackDivider borderColor="gray.200"/>}>
      {groups.map(group => (
        <VStack key={group[0].title}>
          {group.map(article => (
            <Column {...article} key={article.title}/>
          ))}
        </VStack>
      ))}
    </HStack>
  )
}


function Column(props: Article) {
  return <Box w="100%">
    <Flex>
      <Image
        src={props.thumbnail}
        alt={props.title}
        objectFit="cover"
        boxSize="100px"
        borderRadius="md"/>
      <VStack alignItems="flex-start" w="100%" px={4} justifyContent="left">
        <Heading mb={1} fontSize="20px" align="start" fontWeight="medium">{props.title}</Heading>
        <Text noOfLines={2} fontSize="16px" align="start" fontWeight="regular">{props.excerpt}</Text>
      </VStack>
    </Flex>
  </Box>
}

import {Box, Flex, Heading, Image, Stack, StackDivider, Text, useColorModeValue, VStack} from "@chakra-ui/react";
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

  let borderColor = useColorModeValue('gray.400', 'gray.600');

  if (!articles) return <Box/>

  let groups = Utils.groupArray(articles, 2);
  return (
    <Stack
      direction={{base: 'column', md: 'row'}}
      divider={<StackDivider
        display={{base: 'none', md: 'inherit'}}
        borderColor={borderColor}/>}>

      {groups.map((group) => <Box
        w={{base: 'full', md: '50%'}}>
        <VStack key={group[0].title}>
          {group.map(function (article, index) {
            let isLast = index === (group.length - 1)
            return <Box w="full" pb={isLast ? 0 : 4} px={4}>
              <Column {...article} key={article.title}/>
            </Box>;
          })}
        </VStack>
      </Box>)}
    </Stack>
  )
}


function Column(props: Article) {
  return <Flex>
    <Image
      src={props.thumbnail}
      alt={props.title}
      objectFit="cover"
      boxSize="100px"
      borderRadius="md"/>
    <VStack alignItems="flex-start" w="full" px={4} justifyContent="left">
      <Heading mb={1} fontSize="20px" align="start" fontWeight="bold">{props.title}</Heading>
      <Text noOfLines={3} fontSize="16px" align="start" fontWeight="regular">{props.excerpt}</Text>
    </VStack>
  </Flex>
}

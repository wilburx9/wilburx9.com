import {Box, Flex, Heading, Image, Text, useColorModeValue, VStack} from "@chakra-ui/react";
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

  let borderColor = useColorModeValue('gray.400', 'red.600');

  if (!articles) return <Box/>

  let groups = Utils.groupArray(articles, 2);
  return (
    <Flex direction={{base: 'column-reverse', md: 'row'}}>
      {groups.map((group, index) =>
        <Box
          borderLeftColor={borderColor} // TODO: Not working. Fix!
          borderLeft={{base: '0px', md: (index === 0) ? '0px' : '1px'}}
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
    </Flex>
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

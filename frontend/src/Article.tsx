import {Box, Flex, Heading, HStack, Image, StackDivider, Text, VStack} from "@chakra-ui/react";
import React from "react";

export const Article = () => {
  return (
    <HStack divider={<StackDivider borderColor="gray.200"/>}>
      <Column/>
      <Column/>
    </HStack>
  )
}

function Column() {
  return <Box w="100%">
    <Flex>
      <Image
        src="https://source.unsplash.com/300x300/?cute,cat"
        alt="Post title"
        objectFit="cover"
        boxSize="100px"
        borderRadius="md"/>
      <VStack alignItems="center" w="100%" px={4} justifyContent="center">
        <Heading mb={1} fontSize="20px"  align="start"  fontWeight="medium">Modern online and offline payments for Africa</Heading>
        <Text noOfLines={2} fontSize="16px" align="start" fontWeight="regular">
          "The quick brown fox jumps over the lazy dog" is an English-language pangram—a
          sentence that contains all of the letters of the English alphabet. Owing to
          its existence, Chakra was created.
        </Text>
      </VStack>
    </Flex>
  </Box>
}
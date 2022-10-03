import {Box, HStack, Icon, LinkBox, LinkOverlay, Text} from "@chakra-ui/react";
import React from "react";
import {SiGo, SiReact} from "react-icons/si";

export const AttributionComponent = function () {
  let year = new Date().getFullYear();
  return <Box w='full'
              pt={20}
              pb={4}
              alignItems='center'>
    <LinkBox mt={{base: 4, md: 0}}>
      <LinkOverlay href='https://github.com/wilburt/wilburx9.dev'
                   isExternal>
        <HStack justify='center'>
          <Text fontSize='xs'>Wiburx9 &copy; {year} &#8226; Built with</Text>
          <Icon as={SiGo}
                color='#8D949D'
                w={6}
                h={6}/>
          <Icon as={SiReact}
                color='#00D8FF'/>
        </HStack>
      </LinkOverlay>
    </LinkBox>
  </Box>;
}


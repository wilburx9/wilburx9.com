import {Flex, HStack, Icon, LinkBox, LinkOverlay, Text} from "@chakra-ui/react";
import React from "react";
import {IconType} from "react-icons";
import {AiFillGithub, AiFillLinkedin, AiOutlineTwitter, SiGo, SiReact} from "react-icons/all";

class Attr {
  name: string
  url: string
  icon: IconType

  constructor(name: string, url: string, icon: IconType) {
    this.name = name;
    this.url = url;
    this.icon = icon;
  }

}

let attributions: Attr[] = [
  new Attr("LinkedIn", "https://www.linkedin.com/in/wilburx9", AiFillLinkedin),
  new Attr("Github", "https://github.com/wilburt", AiFillGithub),
  new Attr("Twitter", "https://twitter.com/wilburx09", AiOutlineTwitter),
];

export const AttributionComponent = () => (
  <Flex flexDir={{base: 'column', md: 'row'}} w='full' mt={8} mb={5} align={{base: 'center', md: 'flex-end'}}>
    <HStack spacing={4} flexGrow={1} justify='center'>
      {attributions.map(a => <LinkBox key={a.url}>
        <LinkOverlay href={a.url} isExternal>
          <Icon as={a.icon!} boxSize={8}/>
        </LinkOverlay>
      </LinkBox>)}
    </HStack>
    <LinkBox mt={{base: 4, md: 0}}>
      <LinkOverlay href='https://github.com/wilburt/wilburx9.dev' isExternal>
        <HStack>
          <Text fontSize='xs'>Built with</Text>
          <Icon as={SiGo} color='#00ADD8' w={6} h={6}/>
          <Icon as={SiReact} color='#00D8FF'/>
        </HStack>
      </LinkOverlay>
    </LinkBox>
  </Flex>
)


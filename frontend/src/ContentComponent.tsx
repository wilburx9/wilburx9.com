import React, {useContext, useEffect} from "react";
import {DataContext} from "./DataProvider";
import {Box, Center, Fade, Flex, Spinner, useColorModeValue} from "@chakra-ui/react";
import {ColorModeSwitcher} from "./theme/ColorModeSwitcher";
import {ArticlesComponent} from "./articles/ArticlesComponent";
import {ReposComponent} from "./repos/ReposComponent";
import {AttributionComponent} from "./footer/AttributionComponent";
import {ContactComponent} from "./contact/ContactComponent";

export const ContentComponent = () => {
  const {hasData} = useContext(DataContext)

  useEffect(() => {
    hasData?.()
  }, [hasData])

  let footerColorBg = useColorModeValue('gray.100', 'gray.900');
  let hasAnyData = hasData?.();

  if (hasAnyData == null || !hasAnyData) {
    return <Fade in={!hasAnyData} unmountOnExit transition={{exit: {duration: 2}}}>
      <Center h='100vh' w='100vw'> <Spinner size='xl'/></Center>
    </Fade>
  }

  return (
    <Fade in={hasAnyData} unmountOnExit transition={{enter: {duration: 2}}}>
      <Flex flexDir="column">
        <ColorModeSwitcher ml='auto'/>
        <Box
          flex='1'
          maxW="container.xl"
          px={[5, null, 10]}
          alignSelf="center">
          <ArticlesComponent/>
          <ReposComponent/>
          <ContactComponent/>
        </Box>
        <Box bg={footerColorBg}>
          <Box maxW="container.xl" mx='auto'>
            <AttributionComponent />
          </Box>
        </Box>
      </Flex>
    </Fade>
  );
}
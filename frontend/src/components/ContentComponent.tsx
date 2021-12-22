import React, {useContext, useEffect} from "react";
import {DataContext} from "../DataProvider";
import {Box, Center, Fade, Flex, Spinner} from "@chakra-ui/react";
import {ColorModeSwitcher} from "./ColorModeSwitcher";
import {ArticlesComponent} from "./ArticlesComponent";
import {ReposComponent} from "./ReposComponent";
import {AttributionComponent} from "./AttributionComponent";

export const ContentComponent = () => {
  const {hasData} = useContext(DataContext)

  useEffect(() => {
    hasData?.()
  }, [hasData])
  console.log("Has data? " + hasData?.())

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
          <AttributionComponent/>
        </Box>
      </Flex>
    </Fade>
  );
}
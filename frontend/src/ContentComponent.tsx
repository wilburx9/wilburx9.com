import React, {useContext, useEffect} from "react";
import {DataContext} from "./DataProvider";
import {Box, Center, Fade, Spinner, VStack} from "@chakra-ui/react";
import {ArticlesComponent} from "./articles/ArticlesComponent";
import {ReposComponent} from "./repos/ReposComponent";
import {AttributionComponent} from "./footer/AttributionComponent";
import {ContactComponent} from "./contact/ContactComponent";
import {HeaderComponent} from "./header/HeaderComponent";

export const ContentComponent = () => {
  const {hasData} = useContext(DataContext)

  useEffect(() => {
    hasData?.()
  }, [hasData])

  let hasAnyData = hasData?.();

  if (hasAnyData == null || !hasAnyData) {
    return <Fade in={!hasAnyData} unmountOnExit transition={{exit: {duration: 2}}}>
      <Center h='100vh' w='100vw'> <Spinner size='xl'/></Center>
    </Fade>
  }

  return (
    <Fade in={hasAnyData} unmountOnExit transition={{enter: {duration: 2}}}>
      <Box maxW="container.xl" mx='auto'>
        <VStack>
          <HeaderComponent/>
          <ArticlesComponent/>
          <ReposComponent/>
          <ContactComponent/>
          <AttributionComponent/>
        </VStack>
      </Box>
    </Fade>
  );
}

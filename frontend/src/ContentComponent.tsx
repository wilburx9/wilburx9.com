import React from "react";
import {Box, VStack} from "@chakra-ui/react";
import {ArticlesComponent} from "./articles/ArticlesComponent";
import {ReposComponent} from "./repos/ReposComponent";
import {AttributionComponent} from "./footer/AttributionComponent";
import {ContactComponent} from "./contact/ContactComponent";
import {TopComponent} from "./top/TopComponent";

export const ContentComponent = () => (
  <Box maxW="container.xl" mx='auto' px={{base: '2.5', sm: '4', xl: '0'}}>
    <VStack>
      <TopComponent/>
      <ArticlesComponent/>
      <ReposComponent/>
      <ContactComponent/>
      <AttributionComponent/>
    </VStack>
  </Box>
)

import React from "react";
import {Box, Button, FormControl, Heading, HStack, Input, Textarea, useColorModeValue, VStack} from "@chakra-ui/react";
import {HiOutlineArrowRight} from "react-icons/all";
import HCaptcha from '@hcaptcha/react-hcaptcha';

export const ContactComponent = () => {

  function handleCaptchaSuccess(token: string) {

  }

  let bColor = useColorModeValue('black', 'white');
  let focusStyle = {borderColor: bColor, boxShadow: 'none', borderWidth: "1px"}
  return <Box mt={4} my={6} py={6} px={20} borderRadius='xl' bg={useColorModeValue('gray.100', 'gray.900')}>
    <FormControl>
      <Heading mb={4} size='lg' align="start">&#47;&#47;Let's work together</Heading>
      <HStack alignItems='flex-end' spacing={10}>

        <VStack flexGrow={3} spacing={6}>
          <HStack w='full' spacing={6}>
            <Input _focus={focusStyle} id='name' autoComplete='name' type='text'
                   placeholder='Name' variant='filled' isInvalid={false} isRequired={true}/>
            <Input _focus={focusStyle} id='email' autoComplete='email' type='email'
                   placeholder='Email' variant='filled' isInvalid={false} isRequired={true}/>
          </HStack>

          <Input _focus={focusStyle} id='subject' type='text' placeholder='Subject'
                 variant='filled' isInvalid={false} isRequired={true}/>
          <Textarea _focus={focusStyle} id='message' resize='vertical' placeholder='Message'
                    variant='filled' isInvalid={false} isRequired={true}/>
        </VStack>

        <VStack spacing={6}>
          <HCaptcha sitekey={process.env.REACT_APP_H_CAPTCHA_SITE_KEY!}
                    theme={useColorModeValue('light', 'dark')} size='normal'
                    onVerify={(token) => handleCaptchaSuccess(token)}/>
          <Button color='white' isLoading={false} type='submit' borderRadius='lg' loadingText='Sending'
                  spinnerPlacement='end'
                  rightIcon={<HiOutlineArrowRight/>} w='full' disabled={false}>
            Send
          </Button>
        </VStack>
      </HStack>
    </FormControl>
  </Box>
}
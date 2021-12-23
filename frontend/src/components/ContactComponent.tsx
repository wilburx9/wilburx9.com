import React from "react";
import {Box, Button, FormControl, HStack, Input, Textarea, useColorModeValue, VStack} from "@chakra-ui/react";
import {HiOutlineArrowRight} from "react-icons/all";
import HCaptcha from '@hcaptcha/react-hcaptcha';

export const ContactComponent = () => {

  function handleCaptchaSuccess(token: string) {

  }

  return <Box mt={4} my={6} py={6} px={20} borderRadius='xl' bg={useColorModeValue('gray.100', 'gray.900')}>
    <FormControl>
      <HStack alignItems='flex-end'>

        <VStack flexGrow={1}>
          <HStack w='full'>
            <Input id='name' autoComplete='name' type='text' placeholder='Name' isInvalid={false} isRequired={true}/>
            <Input id='email' autoComplete='email' type='email' placeholder='Email'
                   isInvalid={false} isRequired={true}/>
          </HStack>
          <Input id='subject' type='text' placeholder='Subject' isInvalid={false} isRequired={true}/>
          <Textarea id='message' resize='vertical' placeholder='Message' isInvalid={false} isRequired={true}/>
        </VStack>

        <VStack>
          <HCaptcha sitekey={process.env.REACT_APP_H_CAPTCHA_SITE_KEY!}
                    theme={useColorModeValue('light', 'dark')}
                    onVerify={(token) => handleCaptchaSuccess(token)}/>
          <Button isLoading={false} type='submit' loadingText='Sending' spinnerPlacement='end'
                  rightIcon={<HiOutlineArrowRight/>} w='full' disabled={false}>
            Send
          </Button>
        </VStack>
      </HStack>
    </FormControl>
  </Box>
}
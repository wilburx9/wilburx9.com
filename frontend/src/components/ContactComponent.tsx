import React from "react";
import {
  Box,
  Button, FormControl,
  FormErrorMessage,
  Heading,
  HStack,
  Input,
  Textarea, useColorMode,
  VStack
} from "@chakra-ui/react";
import {HiOutlineArrowRight} from "react-icons/all";
import HCaptcha from '@hcaptcha/react-hcaptcha';
import {Form, Formik, Field} from 'formik';
import * as Yup from 'yup';

export const ContactComponent = () => {
  const {colorMode} = useColorMode()

  function handleCaptchaSuccess(token: string) {
    console.log("HCaptcha token  = " + token)
  }

  let isLightTheme = colorMode === 'light'

  return <Box mt={4} my={6} py={6} px={20} borderRadius='xl' bg={isLightTheme ? 'gray.100' : 'gray.900'}>
    <Heading mb={4} size='lg' align="start">&#47;&#47;Let's work together</Heading>
    <Formik
      initialValues={{name: '', email: '', subject: '', message: ''}}
      onSubmit={(values, actions) => {
        setTimeout(() => {
          alert(JSON.stringify(values, null, 2))
          actions.setSubmitting(false)
        }, 1000)
      }}
      validationSchema={Yup.object({
        name: Yup.string().required('Required'),
        email: Yup.string().email('Invalid email address').required('Required'),
        subject: Yup.string().required('Required'),
        message: Yup.string().required('Required'),
      })}
    >
      {(formik) => {
        return (
          <Form>
            <HStack align='flex-end' spacing={10}>

              <VStack flexGrow={3} spacing={6}>
                <HStack w='full' spacing={6} align='flex-start'>
                  <Field name='name'>
                    {({field, form}: { field: any; form: any }) => (
                      <FormControl isInvalid={form.errors.name && form.touched.name}>
                        <Input {...field}  id='name' autoComplete='name' type='text' placeholder='Name' />
                        <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name='email'>
                    {({field, form}: { field: any; form: any }) => (
                      <FormControl isInvalid={form.errors.email && form.touched.email}>
                        <Input {...field}  id='email' autoComplete='email' type='email' placeholder='Email' />
                        <FormErrorMessage>{form.errors.email}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </HStack>
                <Field name='subject'>
                  {({field, form}: { field: any; form: any }) => (
                    <FormControl isInvalid={form.errors.subject && form.touched.subject}>
                      <Input {...field}  id='subject' type='text' placeholder='Subject' />
                      <FormErrorMessage>{form.errors.subject}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name='message'>
                  {({field, form}: { field: any; form: any }) => (
                    <FormControl isInvalid={form.errors.message && form.touched.message}>
                      <Textarea {...field}  id='message' resize='vertical' placeholder='Message'  />
                      <FormErrorMessage>{form.errors.message}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
              </VStack>

              <VStack spacing={6}>
                <HCaptcha sitekey={process.env.REACT_APP_H_CAPTCHA_SITE_KEY!}
                          theme={isLightTheme ? 'light' : 'dark'}
                          onVerify={(token) => handleCaptchaSuccess(token)}/>
                <Button color='white' isLoading={formik.isSubmitting} type='submit' borderRadius='lg'
                        loadingText='Sending'
                        spinnerPlacement='end'
                        rightIcon={<HiOutlineArrowRight size='20px'/>} w='full'>
                  Send
                </Button>
              </VStack>
            </HStack>
          </Form>
        );
      }}
    </Formik>
  </Box>
}
import React, {useContext, useRef, useState} from "react";
import {
  Box,
  Button, FormControl,
  FormErrorMessage,
  Heading,
  HStack,
  Input, Textarea, useColorMode, useToast,
  VStack
} from "@chakra-ui/react";
import {HiOutlineArrowRight} from "react-icons/all";
import HCaptcha from '@hcaptcha/react-hcaptcha';
import {Form, Formik, Field, FormikHelpers} from 'formik';
import * as Yup from 'yup';
import {DataContext} from "../DataProvider";
import { ContactData } from "../models/ContactModel";

export const ContactComponent = () => {
  const {postEmail} = useContext(DataContext)
  let isLightTheme = useColorMode().colorMode === 'light'
  let captchaRef = useRef<HCaptcha>(null);
  const [token, setToken] = useState<string | null>('');
  const toast = useToast()

  async function handleValidForm(values: FormData, token: string, actions: FormikHelpers<FormData>) {
    let data: ContactData = {
      sender_name: values.name.trim(),
      sender_email: values.email.trim(),
      subject: values.subject.trim(),
      message: values.message.trim(),
      captcha_response: token
    }
    let response = await postEmail?.(data)

    toast({
      status: response?.success ? 'success' : 'error',
      title: response?.success ? 'Success' : 'Error',
      description: response?.message,
      duration: 5000,
      isClosable: true
    })

    setToken(null)
    captchaRef.current?.resetCaptcha()
    actions.resetForm()
  }

  function onFormSubmit(values: FormData, actions: FormikHelpers<FormData>) {
    if (token && token.length > 0) {
      handleValidForm(values, token, actions)
    } else {
      captchaRef.current?.execute({async: true}).then(({response}) => {
        return handleValidForm(values, response, actions);
      })
    }
  }

  return <Box mt={4} my={6} py={6} px={20} borderRadius='xl' bg={isLightTheme ? 'gray.100' : 'gray.900'}>
    <Heading mb={6} size='lg' align="start">&#47;&#47;Let's work together</Heading>
    <Formik<FormData>
      initialValues={{name: '', email: '', subject: '', message: ''}}
      onSubmit={(values, actions) => onFormSubmit(values, actions)}
      validationSchema={Yup.object({
        name: Yup.string().trim().required('Required'),
        email: Yup.string().trim().email('Invalid email address').required('Required'),
        subject: Yup.string().trim().required('Required'),
        message: Yup.string().trim().required('Required'),
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
                        <Input {...field} id='name' autoComplete='name' type='text' placeholder='Name'/>
                        <FormErrorMessage>{form.errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name='email'>
                    {({field, form}: { field: any; form: any }) => (
                      <FormControl isInvalid={form.errors.email && form.touched.email}>
                        <Input {...field} id='email' autoComplete='email' type='email' placeholder='Email'/>
                        <FormErrorMessage>{form.errors.email}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </HStack>
                <Field name='subject'>
                  {({field, form}: { field: any; form: any }) => (
                    <FormControl isInvalid={form.errors.subject && form.touched.subject}>
                      <Input {...field} id='subject' type='text' placeholder='Subject'/>
                      <FormErrorMessage>{form.errors.subject}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name='message'>
                  {({field, form}: { field: any; form: any }) => (
                    <FormControl isInvalid={form.errors.message && form.touched.message}>
                      <Textarea {...field} id='message' resize='vertical' placeholder='Message'/>
                      <FormErrorMessage>{form.errors.message}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
              </VStack>

              <VStack spacing={6}>
                <HCaptcha sitekey={process.env.REACT_APP_H_CAPTCHA_SITE_KEY!}
                          theme={isLightTheme ? 'light' : 'dark'}
                          ref={captchaRef}
                          onExpire={() => setToken(null)}
                          onVerify={setToken}/>
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

type FormData = {
  name: string,
  email: string,
  subject: string,
  message: string
}

export type FormResponse = {
  success: boolean;
  message: string
}
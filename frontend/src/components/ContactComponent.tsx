import React, {useContext, useRef, useState} from "react";
import {
  Box,
  Button, FormControl,
  FormErrorMessage,
  Heading,
  Input, Stack, Textarea, useBreakpointValue, useColorMode, useToast,
  VStack
} from "@chakra-ui/react";
import {HiOutlineArrowRight} from "react-icons/all";
import HCaptcha from '@hcaptcha/react-hcaptcha';
import {Form, Formik, Field, FormikHelpers} from 'formik';
import * as Yup from 'yup';
import {DataContext} from "../DataProvider";
import {ContactData} from "../models/ContactModel";
import {getAnalyticsParams, logAnalyticsEvent} from "../analytics/firebase";
import {AnalyticsEvent} from "../analytics/events";
import {AnalyticsKey} from "../analytics/keys";
import {getInputFilledStyle} from "../theme";

export const ContactComponent = () => {
  const {postEmail} = useContext(DataContext)
  const [token, setToken] = useState<string | null>('');
  const toast = useToast()
  let isLightMode = useColorMode().colorMode === 'light'
  let captchaRef = useRef<HCaptcha>(null);
  let isNormalCaptchaSize = useBreakpointValue({base: false, md: false, lg: true})
  let isSmallButton = useBreakpointValue({base: false, lg: true})


  async function handleValidForm(values: FormData, token: string, actions: FormikHelpers<FormData>) {
    let data = new ContactData(
      values.name.trim(),
      values.email.trim(),
      values.subject.trim(),
      values.message.trim(),
      token
    )
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
    if (response?.success === true) {
      actions.resetForm()
    } else {
      actions.setSubmitting(false)
    }
  }

  function onFormSubmit(values: FormData, actions: FormikHelpers<FormData>) {
    if (token && token.length > 0) {
      handleValidForm(values, token, actions)
    } else {
      captchaRef.current?.execute({async: true}).then(({response}) => {
        return handleValidForm(values, response, actions);
      }).catch((reason) => {
        actions.setSubmitting(false)
        logCaptchaError(reason)
      })
    }
  }

  return <Box mt={4} my={6} py={6} px={{base: 5, md: 6, lg: 20}} borderRadius='xl'
              bg={isLightMode ? 'gray.50' : 'gray.900'}>
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
            <Stack direction={{base: 'column', md: 'row'}} align={{base: 'stretch', md: 'flex-end'}} spacing={10}>
              <VStack flexGrow={3} spacing={6}>
                <Stack direction={{base: 'column', md: 'row'}} w='full' spacing={6} align='flex-start'>
                  <Field name='name'>
                    {({field, form}: { field: any; form: any }) => (
                      <FormControl isInvalid={form.errors.name && form.touched.name}>
                        <Input {...field} id='name' autoComplete='name' type='text' placeholder='Name'
                               colorScheme='blackAlpha'/>
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
                </Stack>
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
                      {/*For some reason, the input style applied in the custom theme for Textarea doesn't work. So I'm applying ut here manually*/}
                      <Textarea {...field} {...getInputFilledStyle(isLightMode)} id='message' resize='vertical'
                                placeholder='Message'/>
                      <FormErrorMessage>{form.errors.message}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
              </VStack>

              <Stack direction={{base: 'row-reverse', md: 'column'}} spacing={6}
                     align={{base: 'flex-start', md: 'center'}}>
                <HCaptcha sitekey={process.env.REACT_APP_H_CAPTCHA_SITE_KEY!}
                          theme={isLightMode ? 'light' : 'dark'}
                          size={isNormalCaptchaSize ? 'normal' : 'compact'}
                          onError={logCaptchaError}
                          ref={captchaRef}
                          onExpire={() => setToken(null)}
                          onVerify={setToken}/>
                <Button color='white' isLoading={formik.isSubmitting} type='submit' borderRadius='lg'
                        loadingText='Sending'
                        spinnerPlacement='end'
                        px={{base: 0, md: 20, lg: 0}}
                        size={isSmallButton ? 'md' : 'lg'}
                        bg={isLightMode ? 'blackAlpha.600' : 'whiteAlpha.100'}
                        rightIcon={<HiOutlineArrowRight size='20px'/>}
                        w='full'
                        _hover={{background: isLightMode ? 'blackAlpha.800' : 'whiteAlpha.200'}}
                        _active={{background: isLightMode ? 'blackAlpha.900' : 'whiteAlpha.300'}}
                >
                  Send
                </Button>
              </Stack>
            </Stack>
          </Form>
        );
      }}
    </Formik>
  </Box>
}

function logCaptchaError(reason: any) {
  let params = getAnalyticsParams()
  params.set(AnalyticsKey.reason, reason)
  logAnalyticsEvent(AnalyticsEvent.captchaFailure, params)
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
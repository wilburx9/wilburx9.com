import {redact} from "../common/Utils";

export interface ContactResponse {
  data?: string;
  success?: boolean;
}

export class ContactData {
  sender_name: string;
  sender_email: string;
  subject: string
  message: string
  captcha_response: string


  constructor(name: string, email: string, subject: string, message: string, captchaResponse: string) {
    this.sender_name = name;
    this.sender_email = email;
    this.subject = subject;
    this.message = message;
    this.captcha_response = captchaResponse;
  }

  public redact(): ContactData {
    this.sender_name = redact(this.sender_name)
    this.sender_email = redact(this.sender_email)
    this.subject = redact(this.subject)
    this.message = redact(this.message)
    this.captcha_response = redact(this.captcha_response, 4)
    return this
  }
}



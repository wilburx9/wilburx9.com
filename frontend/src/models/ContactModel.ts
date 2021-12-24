export interface ContactResponse {
  data?:    string;
  success?: boolean;
}

export type ContactData = {
  sender_name: string;
  sender_email: string;
  subject: string
  message: string
  captcha_response: string
}

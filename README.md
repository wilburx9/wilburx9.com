# wilburx9.com
The source code for wilburx9.com, which is built on [Ghost](https://ghost.org). 
The code lives in two directories _frontend_ and _backend_; the former is a Ghost theme, and the latter contains Go services 
that manage newsletters.

## Frontend
The custom template which was originally cloned from the [Ghost Starter theme ](https://github.com/TryGhost/Starter).
It's written in handlebars, css and a mixture of JQuery and vanilla JS.

### Routing
Deviating from the usual Ghost themes, I repurposed "/" into a self-aggrandizing landing page. 
So the blog lives on "/blog" and the listing pages of individual tags are on "/blog/tag_slug". 
This was implemented using a [custom route](frontend/routes.yaml).
Consequently, an empty page was created on Ghost dashboard with the slug "blog" and a custom title and description. 
See the [Ghost doc](https://ghost.org/docs/themes/routing/#the-default-collection) on this.

### External Articles
Some of my articles live on Kodeco and Medium, and can't be imported to Ghost. 
However, I added "external articles" that contain nothing but bookmark cards pointing to the original article. 
These articles have the "#external" private tag. Special provision has been made such so that such articles are not indexed, and clicking the post-cards from any such article on any article listing section opens the original article.

### Running Locally
1. Clone this repo
2. [Install Ghost](https://ghost.org/tutorials/local-ghost/).
3. Create a symlink in Ghost's installations theme directory that points to this repo's frontend directory by running `ln -s  PROJECT_DIR/frontend GHOST_DIR/content/themes/wilburx9`.
4. `cd PROJECT_DIR/frontend` and run `yarn dev` to build and watch for new changes.
5. `cd GHOST_DIR` and run `ghost restart`.
6. Go to installed themes in the design settings of Ghost dashboard and select "wilburx9".

### Deploying
Deployment can be done manually or using the [CD workflow](.github/workflows/frontend.yaml).
* **Manually**
  1. Clone the repo
  2. `cd PROJECT_DIR/frontend`.
  3. `yarn pretest` to build.
  4. `yarn zip` to package the theme into a zip file.
  5. Go to Theme settings and upload the generated zip file.
* **CD Workflow**: The steps for building, deploying and applying the theme are packaged into a CD workflow which runs when there's a push event on the _live_ branch that changed the frontend directory. To take advantage of this:
  1. Fork the repo
  2. Create a [Ghost custom integration](https://ghost.org/integrations/custom-integrations/), add Admin API Key and Url to your [projects secrets](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository) and add these as `GHOST_ADMIN_API_KEY` and `GHOST_ADMIN_API_URL` respectively.
  3. Add `GEN_SOURCEMAPS` to [GitHub variables](https://docs.github.com/en/actions/learn-github-actions/variables) with `true` or `false` depending on if you want to deploy the theme along with the source maps.
  4. After deployment, the workflow clears the cdn cache to ensure the new theme reflects immediately. And this works on the assumption that your website is deployed on a Lightsail instance behind a Lightsail Distribution. So create a [GitHub OIDC provider AWS on aws IAM](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services) and ensure [the policy has a _Resource_](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_job-functions_create-policies.html) that points to the ARN of the Lightsail instance's Distribution and has a `lightsail:ResetDistributionCache` in the `Actions` array. Then, add the ARN of the role and AWS region to GitHub secrets as `AWS_ROLE` and `AWS_REGION`.

## Backend
Ghost's newsletter implementation is not very customizable as I wanted readers to choose the kind of newsletter to want to receive. 
Hence, the backend directory contains two Go Services deployed to AWS Lambda:
 - **subscribe** service receives the email and newsletter preferences from the frontend, validates the captcha and forwards it to MailerLite.
 - **broadcast** when a new article is published, Ghost hits this webhook to broadcast to appropriate subscribers.

### Deploying
Deployment is done by a [CD workflow](.github/workflows/backend.yaml) which is triggered when there's a new push event on the _live_ branch that changed the backend directory.
* **Secrets**: AWS Systems Manager Parameter Store is used to store the secrets used by both services. These secrets are:
  1. `WILBURX9_ALLOWED_ORIGINS`: The origins allowed by the Lambdas; should be the site's homepage. Without this, the subscription form will fail because of good old CORS error.
  2. `WILBURX9_EMAIL_SENDER`: The email for the sender of the newsletter.
  3. `WILBURX9_MAILER_LITE_TOKEN`: [API token for MailerLite](https://www.mailerlite.com/help/where-to-find-the-mailerlite-api-key-groupid-and-documentation).
  4. `WILBURX9_TURNSTILE_HOSTNAME`: Your website domain configured on Cloudflare [Tunrnstile](https://developers.cloudflare.com/turnstile/) dashboard.
  5. `WILBURX9_TURNSTILE_SECRET`: Turnstile site's Secret key.
* Create Lambdas
  1. Create two Lambda Functions on AWS using the Go 1.x runtime and x86_64 architecture. Ensure their roles has a statement that allows reading from `ssm:GetParameter`; this is to ensure the services can read the secrets created above.
  2. Add the function names to [GitHub variables](https://docs.github.com/en/actions/learn-github-actions/variables) using  `LAMBDA_FUNCTION_SUBSCRIBE` and `LAMBDA_FUNCTION_BROADCAST`.
* IAM Role
  1. Add `lambda:UpdateFunctionCode` to the `Actions` array of the IAM role created when deploying the frontend. This is so the CLI pipeline can update Lambda function.
  2. Add the ARNs of the Lambda functions to the  `Resource` array of same IAM role.
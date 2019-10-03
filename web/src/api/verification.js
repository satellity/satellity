class Verification {
  constructor(api) {
    this.api = api;
  }

  create(params) {
    const data = {purpose: params.purpose, email: params.email, recaptcha: params.recaptcha};
    return this.api.axios.post('/email_verifications', data);
  }
}

export default Verification;

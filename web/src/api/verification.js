class Verification {
  constructor(api) {
    this.api = api;
  }

  create(params) {
    params = {'email': params.email, 'recaptcha': params.recaptcha};
    return this.api.axios.post('/email_verifications', params).then((resp) => {
      return resp.data;
    });
  }
}

export default Verification;

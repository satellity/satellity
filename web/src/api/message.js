class Message {
  constructor(api) {
    this.api = api;
  }

  index(id, offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B')
    }
    return this.api.axios.get(`/groups/${id}/messages?offset=${offset}`).then((resp) => {
      return resp.data;
    });
  }

  create(id, params) {
    const data = {body: params.body};
    return this.api.axios.post(`/groups/${id}/messages`, data).then((resp) => {
      return resp.data;
    });
  }
}

export default Message;

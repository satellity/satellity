class Comment {
  constructor(api) {
    this.api = api;
  }

  index(id) {
    return this.api.axios.get(`/topics/${id}/comments`).then((resp) => {
      return resp.data;
    });
  }

  create(params) {
    const data = {topic_id: params.topic_id, body: params.body};
    return this.api.axios.post('/comments', data).then((resp) => {
      return resp.data;
    });
  }

  update(id, params) {
    const data = {topic_id: params.topic_id, body: params.body};
    return this.api.axios.post(`/comments/${id}`, data).then((resp) => {
      return resp.data;
    });
  }

  show(id) {
    return this.api.axios.get(`/comments/${id}`).then((resp) => {
      return resp.data;
    });
  }

  delete(id) {
    return this.api.axios.delete(`/comments/${id}`).then((resp) => {
      return resp;
    });
  }
}

export default Comment;

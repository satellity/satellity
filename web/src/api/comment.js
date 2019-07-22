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
    return this.api.axios.post('/comments', params).then((resp) => {
      return resp.data;
    });
  }

  update(id, params) {
    return this.api.axios.post(`/comments/${id}`, params).then((resp) => {
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

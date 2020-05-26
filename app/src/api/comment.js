class Comment {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
  }

  index(id) {
    return this.api.axios.get(`/topics/${id}/comments`);
  }

  create(params) {
    const data = {topic_id: params.topic_id, body: params.body};
    return this.api.axios.post('/comments', data);
  }

  update(id, params) {
    const data = {topic_id: params.topic_id, body: params.body};
    return this.api.axios.post(`/comments/${id}`, data);
  }

  show(id) {
    return this.api.axios.get(`/comments/${id}`);
  }

  delete(id) {
    return this.api.axios.delete(`/comments/${id}`);
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/comments');
  }

  delete(id) {
    return this.api.axios.delete(`/admin/comments/${id}`)
  }
}

export default Comment;

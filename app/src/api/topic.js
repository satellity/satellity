const i18n = window.i18n;

class Topic {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
  }

  index(offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B').replace(' ', '%2B');
    }
    return this.api.axios.get(`/topics?offset=${offset}`);
  }

  create(params) {
    const data = {title: params.title, body: params.body.trim(), topic_type: params.topic_type, category_id: params.category_id, draft: params.draft};
    return this.api.axios.post('/topics', data);
  }

  update(id, params) {
    const data = {title: params.title, body: params.body.trim(), topic_type: params.topic_type, category_id: params.category_id, draft: params.draft};
    return this.api.axios.post(`/topics/${id}`, data);
  }

  show(id) {
    return this.api.axios.get(`/topics/${id}`);
  }

  action(action, id) {
    if (action === 'like' || action === 'unlike' || action === 'bookmark' || action === 'unsave') {
      return this.api.axios.post(`/topics/${id}/${action}`, {});
    }
    return new Promise((resolve) => {
      resolve({error: i18n.t('topic.action.invalid')});
    });
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/topics');
  }

  delete(id) {
    return this.api.axios.delete(`/admin/topics/${id}`)
  }
}

export default Topic;

require('codemirror/lib/codemirror.css');
require('codemirror/theme/xq-light.css');
require('codemirror/mode/markdown/markdown.js');
import style from './style.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {Controlled as CodeMirror} from 'react-codemirror2'
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';
const validate = require('uuid-validate');

class TopicNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (d !== null && d !== undefined && d !== '') {
      categories = JSON.parse(atob(d));
    }
    let id = this.props.match.params.id;
    if (id === null || id == undefined) {
      id = ''
    }
    this.state = {topic_id: id, title: '', category_id: '', body: '', categories: categories, loading: true, submitting: false};
    this.handleChange = this.handleChange.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
    // TODO handle authentication
    if (!this.api.user.loggedIn()) {
      props.history.push('/');
    }
  }

  componentDidMount() {
    if (validate(this.state.topic_id)) {
      this.api.topic.show(this.props.match.params.id, (resp) => {
        resp.data.loading = false;
        this.setState(resp.data);
      });
    } else {
      this.setState({loading: false});
    }
    this.api.category.index((resp) => {
      let category_id = '';
      if (resp.data.length > 0) {
        category_id = resp.data[0].category_id;
      }
      this.setState({categories: resp.data, category_id: category_id});
    });
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleBodyChange(editor, data, value) {
    this.setState({body: value});
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    const history = this.props.history;
    const data = {title: this.state.title, body: this.state.body, category_id: this.state.category_id};
    if (validate(this.state.topic_id)) {
      this.api.topic.update(this.state.topic_id, data, (resp) => {
        history.push('/');
      });
      return
    }
    this.api.topic.create(data, (resp) => {
      history.push('/');
    });
  }

  render() {
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} onBodyChange={this.handleBodyChange} state={this.state} />
    )
  }
}

// TODO jsx editor format
const View = ({onSubmit, onChange, onBodyChange, state}) => {
  const categories = state.categories.map((c) => {
    return (
      <option value={c.category_id} key={c.category_id}>{c.name}</option>
    )
  });

  let title = <h2>Create a new topic</h2>;
  if (validate(state.topic_id)) {
    title = <h2>Edit: {state.title}</h2>
  }

  const loadingView = (
    <div className={style.form_loading}>
      <LoadingView style='md-ring'/>
    </div>
  )

  return (
    <div className='container'>
      <main className='section main'>
        {state.loading && loadingView}
        <div className={style.form}>
          {title}
          <form onSubmit={(e) => onSubmit(e)}>
            <div>
              <label name='title'>Title *</label>
              <input type='text' name='title' pattern='.{3,}' required value={state.title} autoComplete='off' onChange={(e) => onChange(e)} />
            </div>
            <div>
              <label name='category'>Category</label>
              <div className='select'>
                <select name='category_id' value={state.category_id} onChange={(e) => onChange(e)}>
                  {categories}
                </select>
              </div>
            </div>
            <div className={style.topic_body}>
              <CodeMirror
                value={state.body}
                options={{
                  mode: 'markdown',
                  theme: 'xq-light',
                  lineNumbers: true,
                  lineWrapping: true
                }}
                onBeforeChange={(editor, data, value) => onBodyChange(editor, data, value)}
              />
            </div>
            <div className='action'>
              <input type='submit' value='SUBMIT' />
            </div>
          </form>
        </div>
      </main>
      <aside className='section aside'>
        <ol className={style.rules}>
          <li className={style.rule}>
            1. To be a kind human, keep goodwill towards others
          </li>
          <li className={style.rule}>
            2. It's a good habits to preview before posting.
          </li>
          <li className={style.rule}>
            3. Welcome to share. Enjoy!
          </li>
        </ol>
      </aside>
    </div>
  )
};

export default TopicNew;

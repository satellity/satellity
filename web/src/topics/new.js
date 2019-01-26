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
    this.handleCategoryClick = this.handleCategoryClick.bind(this);
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
      let category_id = this.state.category_id;
      if (category_id === '' && resp.data.length > 0) {
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

  handleCategoryClick(e, value) {
    this.setState({category_id: value});
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
        this.setState({submitting: false});
        history.push('/');
      });
      return
    }
    this.api.topic.create(data, (resp) => {
      this.setState({submitting: false});
      history.push('/');
    });
  }

  render() {
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} onBodyChange={this.handleBodyChange} onCategoryClick={this.handleCategoryClick} state={this.state} />
    )
  }
}

// TODO jsx editor format
const View = (props) => {
  const categories = props.state.categories.map((c) => {
    return (
      <span key={c.category_id} className={`${style.category} ${c.category_id === props.state.category_id ? style.categoryCurrent : ''}`} onClick={(e) => props.onCategoryClick(e, c.category_id)}>{c.alias}</span>
    )
  });

  let title = <h2>Create a new topic</h2>;
  if (validate(props.state.topic_id)) {
    title = <h2>Edit: {props.state.title}</h2>
  }

  const loadingView = (
    <div className={style.form_loading}>
      <LoadingView style='md-ring'/>
    </div>
  )

  return (
    <div className='container'>
      <main className='section main'>
        {props.state.loading && loadingView}
        <div className={style.form}>
          {title}
          <form onSubmit={(e) => props.onSubmit(e)}>
            <div className={style.categories}>
              {categories}
            </div>
            <div>
              <input type='text' name='title' pattern='.{3,}' required value={props.state.title} autoComplete='off' placeholder='Title *' onChange={(e) => props.onChange(e)} />
            </div>
            <div className={style.topic_body}>
              <CodeMirror
                value={props.state.body}
                options={{
                  mode: 'markdown',
                  theme: 'xq-light',
                  lineNumbers: true,
                  lineWrapping: true,
                  placeholder: 'Text (optional)'
                }}
                onBeforeChange={(editor, data, value) => props.onBodyChange(editor, data, value)}
              />
            </div>
            <div className='action'>
              <button className='btn submit' disabled={props.state.submitting}>
                { props.state.submitting && <LoadingView style='sm-ring blank'/> }
                &nbsp;SUBMIT
              </button>
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

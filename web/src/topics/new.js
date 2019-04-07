require('codemirror/lib/codemirror.css');
require('codemirror/theme/xq-light.css');
require('codemirror/mode/markdown/markdown.js');
import style from './style.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import {Controlled as CodeMirror} from 'react-codemirror2'
import showdown from 'showdown';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';
const validate = require('uuid-validate');

class TopicNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.converter = new showdown.Converter();
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (!!d) {
      categories = JSON.parse(atob(d));
    }
    let id = this.props.match.params.id;
    // false , 0 , "" , null , undefined , and NaN
    if (!id) {
      id = ''
    }
    this.state = {
      topic_id: id,
      title: '',
      body: '',
      categories: categories,
      preview: false,
      loading: true,
      submitting: false
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleCategoryClick = this.handleCategoryClick.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handlePreview = this.handlePreview.bind(this);
  }

  componentDidMount() {
    if (validate(this.state.topic_id)) {
      this.api.topic.show(this.props.match.params.id).then((data) => {
        data.loading = false;
        this.setState(data);
      });
    } else {
      this.setState({loading: false});
    }
    this.api.category.index().then((data) => {
      let category_id = this.state.category_id;
      if (!category_id && data.length > 0) {
        category_id = data[0].category_id;
      }
      this.setState({categories: data, category_id: category_id});
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.location.pathname != prevProps.location.pathname) {
      if (this.props.location.pathname === '/topics/new') {
        this.setState({
          topic_id: '',
          title: '',
          body: '',
          preview: false,
        });
      }
    }
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleCategoryClick(e, value) {
    e.preventDefault();
    this.setState({category_id: value});
  }

  handleBodyChange(editor, data, value) {
    this.setState({body: value});
  }

  handlePreview(e) {
    e.preventDefault();
    this.setState({body_html: this.converter.makeHtml(this.state.body), preview: !this.state.preview});
  }

  handleSubmit(e) {
    e.preventDefault();
    if (this.state.submitting) {
      return
    }
    this.setState({submitting: true});
    const history = this.props.history;
    const data = {title: this.state.title, body: this.state.body, category_id: this.state.category_id};
    // TODO should update submitting always
    if (validate(this.state.topic_id)) {
      this.api.topic.update(this.state.topic_id, data).then((data) => {
        this.setState({submitting: false});
        history.push('/');
      });
      return
    }
    this.api.topic.create(data).then((data) => {
      this.setState({submitting: false});
      history.push('/');
    });
  }

  render() {
    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }

    let state = this.state;
    const categories = state.categories.map((c) => {
      return (
        <span key={c.category_id} className={`${style.category} ${c.category_id === state.category_id ? style.categoryCurrent : ''}`} onClick={(e) => this.handleCategoryClick(e, c.category_id)}>{c.alias}</span>
      )
    });

    let title = <h1>{i18n.t('topic.title.new')}</h1>;
    if (validate(state.topic_id)) {
      title = <h1>{i18n.t('topic.title.edit', {name: state.title})}</h1>
    }

    const loadingView = (
      <div className={style.form_loading}>
        <LoadingView style='md-ring'/>
      </div>
    )

    let form = (
      <form onSubmit={this.handleSubmit}>
        <div className={style.categories}>
          {categories}
        </div>
        <div>
          <input type='text' name='title' pattern='.{3,}' required value={state.title} autoComplete='off' placeholder='Title *' onChange={this.handleChange} />
        </div>
        <div className={style.preview}> <FontAwesomeIcon className={style.eye} icon={['far', 'eye']} onClick={this.handlePreview} /> </div>
        <div className={style.topic_body}>
          {
            !state.preview &&
            <CodeMirror
              className='editor'
              value={state.body}
              options={{
                mode: 'markdown',
                theme: 'xq-light',
                lineNumbers: true,
                lineWrapping: true,
                placeholder: 'Text (optional)'
              }}
              onBeforeChange={(editor, data, value) => this.handleBodyChange(editor, data, value)}
            />
          }
          {
            state.preview &&
            <article className={`md ${style.preview_body}`} dangerouslySetInnerHTML={{__html: state.body_html}}>
            </article>
          }
        </div>
        <div className='action'>
          <button className='btn submit' disabled={state.submitting}>
            { state.submitting && <LoadingView style='sm-ring blank'/> }
            &nbsp;SUBMIT
          </button>
        </div>
      </form>
    )

    return (
      <div className='container'>
        <main className='section main'>
          {state.loading && loadingView}
          <div className={style.form}>
            {!state.loading && title}
            {!state.loading && form}
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
  }
}

export default TopicNew;

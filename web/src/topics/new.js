require('codemirror/lib/codemirror.css');
require('codemirror/theme/xq-light.css');
require('codemirror/mode/markdown/markdown.js');
import style from './new.scss';
import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import {Controlled as CodeMirror} from 'react-codemirror2'
import showdown from 'showdown';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Base64 from '../components/base64.js';
import API from '../api/index.js';
import LoadingView from '../loading/loading.js';
const validate = require('uuid-validate');

class New extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.base64 = new Base64();
    this.converter = new showdown.Converter();
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (!!d) {
      categories = JSON.parse(this.base64.decode(d));
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
      draft: false,
      categories: categories,
      preview: false,
      loading: true,
      submitting: false
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleCategoryClick = this.handleCategoryClick.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleDraft = this.handleDraft.bind(this);
    this.handlePreview = this.handlePreview.bind(this);
  }

  componentDidMount() {
    if (validate(this.state.topic_id)) {
      this.api.topic.show(this.state.topic_id).then((data) => {
        data.loading = false;
        this.setState(data);
      });
    } else {
      this.api.topic.show('draft').then((data) => {
        if (!data) {
          data = {};
        }
        data.loading = false;
        this.setState(data);
      });
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
    this.setState({submitting: true, draft: false}, () => {
      this.submitForm();
    });
  }

  handleDraft(e) {
    e.preventDefault();
    this.setState({submitting: true, draft: true}, () => {
      this.submitForm();
    });
  }

  submitForm() {
    const history = this.props.history;
    const data = {title: this.state.title, body: this.state.body, category_id: this.state.category_id, draft: this.state.draft};
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
        <span key={c.category_id} className={`${style.category} ${c.category_id === state.category_id ? style.active : ''}`} onClick={(e) => this.handleCategoryClick(e, c.category_id)}>{c.alias}</span>
      )
    });

    let title = <h1>{i18n.t('topic.title.new')}</h1>;
    if (validate(state.topic_id)) {
      title = <h1>{i18n.t('topic.title.edit', {name: state.title})}</h1>
    }

    const loadingView = (
      <div className={style.loading}>
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
        <div className={style.actions}>
          <FontAwesomeIcon className={style.eye} icon={['far', 'eye']} onClick={this.handlePreview} />
        </div>
        <div className={style.body}>
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
        <article className={`md ${style.preview}`} dangerouslySetInnerHTML={{__html: state.body_html}}>
        </article>
      }
    </div>
    <div>
      {
        state.submitting &&
        <div className={style.submitting}>
          <LoadingView style='sm-ring'/>
          <span> Submitting </span>
        </div>
      }
      {
        !state.submitting &&
        <div>
          <button type="submit" className='btn topic' disabled={state.submitting}>

            &nbsp;{i18n.t('general.submit')}
          </button>
          <a className={style.draft} onClick={this.handleDraft}>{i18n.t('general.draft')}</a>
        </div>
      }
    </div>
  </form>
    )

    return (
      <div className='container'>
        <main className='column main'>
          {state.loading && loadingView}
          <div className={style.form}>
            {!state.loading && title}
            {!state.loading && form}
          </div>
        </main>
        <aside className='column aside'>
          <ol className={style.rules} dangerouslySetInnerHTML={{__html: i18n.t('topic.rules')}}></ol>
        </aside>
      </div>
    )
  }
}

export default New;
